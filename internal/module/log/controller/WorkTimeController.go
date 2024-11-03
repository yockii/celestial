package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/log/domain"
	"github.com/yockii/celestial/internal/module/log/model"
	"github.com/yockii/celestial/internal/module/log/service"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"time"
)

var WorkTimeController = new(workTimeController)

type workTimeController struct{}

func (c *workTimeController) Add(ctx *fiber.Ctx) error {
	instance := new(model.WorkTime)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ProjectID == 0 || instance.WorkTime <= 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId & work time",
		})
	}

	instance.UserID, _ = helper.GetCurrentUserID(ctx)
	if instance.UserID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " userId",
		})
	}

	// 判断填报的数据周期范围，周一到周日的数据，最晚应在下周一9点前填报完成，超时不允许填报
	now := time.Now()

	var earliestAllowedTime time.Time
	// 判断now是否周一9点前
	if now.Weekday() == time.Monday && now.Hour() < 9 {
		// 允许填报上周工作
		earliestAllowedTime = time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, time.Local)
		// 最晚允许的时间是周天0点
	} else {
		// 最早允许当周周一0点
		offset := int(time.Monday - now.Weekday())
		if offset > 0 {
			offset = -6
		}

		weekStart := now.AddDate(0, 0, offset)
		earliestAllowedTime = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.Local)
	}
	latestAllowedTime := earliestAllowedTime.AddDate(0, 0, 7).Add(time.Hour * 9)

	if !(instance.StartDate > earliestAllowedTime.UnixMilli() && instance.EndDate < latestAllowedTime.UnixMilli()) {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  "填报时间必须在本周 周一至下周一9点前",
		})
	}

	duplicated, success, err := service.WorkTimeService.Add(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if duplicated {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		})
	}
	if !success {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *workTimeController) Update(ctx *fiber.Ctx) error {
	instance := new(model.WorkTime)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}

	success, err := service.WorkTimeService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *workTimeController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.WorkTime)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}

	success, err := service.WorkTimeService.Delete(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *workTimeController) List(ctx *fiber.Ctx) error {
	instance := new(domain.WorkTimeListRequest)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	paginate := new(server.Paginate)
	if err := ctx.QueryParser(paginate); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if paginate.Limit == 0 {
		paginate.Limit = 10
	}

	tcList := make(map[string]*server.TimeCondition)
	if instance.StartDateCondition != nil {
		tcList["start_date"] = instance.StartDateCondition
	} else {
		now := time.Now()
		tcList["start_date"] = &server.TimeCondition{
			Start: getThisMondayTimestamp(now),
			End:   getThisSundayTimestamp(now),
		}
	}
	if instance.EndDateCondition != nil {
		tcList["end_date"] = instance.EndDateCondition
	}
	if instance.CreateTimeCondition != nil {
		tcList["create_time"] = instance.CreateTimeCondition
	}
	if instance.ReviewTimeCondition != nil {
		tcList["review_time"] = instance.ReviewTimeCondition
	}

	instance.UserID, _ = helper.GetCurrentUserID(ctx)

	total, list, err := service.WorkTimeService.PaginateBetweenTimes(&instance.WorkTime, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  list,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

func getThisSundayTimestamp(now time.Time) database.DateTime {
	offset := int(time.Sunday - now.Weekday())
	if offset < 0 {
		offset += 7
	}
	weekSunday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	fmt.Println(weekSunday.Format(time.DateTime))
	return database.DateTime(weekSunday)
}

func getThisMondayTimestamp(now time.Time) database.DateTime {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekMonday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	fmt.Println(weekMonday.Format(time.DateTime))
	return database.DateTime(weekMonday)
}

func (c *workTimeController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.WorkTime)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if condition.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	dept, err := service.WorkTimeService.Instance(condition.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: dept,
	})
}

func (c *workTimeController) Statistics(ctx *fiber.Ctx) error {
	instance := new(domain.WorkTimeStatisticsRequest)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	list, err := service.WorkTimeService.Statistics(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: list,
	})
}

func (c *workTimeController) Mine(ctx *fiber.Ctx) error {
	instance := new(domain.WorkTimeListRequest)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	userID, _ := helper.GetCurrentUserID(ctx)
	if userID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " userId",
		})
	}

	data, err := service.WorkTimeService.StatisticsMine(userID, instance.StartDateCondition)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: data,
	})
}
