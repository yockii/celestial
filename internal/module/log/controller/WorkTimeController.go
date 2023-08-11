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
