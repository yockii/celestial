package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/project/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
	"strconv"
	"strings"
)

var ProjectChangeController = new(projectChangeController)

type projectChangeController struct{}

func (c *projectChangeController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectChange)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ProjectID == 0 || instance.Title == "" || instance.Type == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId & title & type",
		})
	}

	// 判断权限
	if uid, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectChangeAdd); err != nil {
		return err
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectChangeService.Add(instance)
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

	status := "待评审"
	_ = ants.Submit(data.AddDocumentAntsWrapper(&search.Document{
		ID:    instance.ID,
		Title: instance.Title,
		Content: fmt.Sprintf("原因:%s\n计划:%s\n评审:%s\n风险:%s\n状态:%s",
			instance.Reason,
			instance.Plan,
			instance.Review,
			instance.Risk,
			status,
		),
		Route:      fmt.Sprintf("/project/detail/%d/change?id=%d", instance.ProjectID, instance.ID),
		CreateTime: instance.CreateTime,
		UpdateTime: instance.UpdateTime,
	}, instance.ApplyUserID))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectChangeController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectChange)
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

	// 先取出旧数据
	oldInstance, err := service.ProjectChangeService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	// 判断权限
	if _, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectChangeUpdate); err != nil {
		return err
	}

	success, err := service.ProjectChangeService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			d, e := service.ProjectChangeService.Instance(id)
			if e != nil {
				logger.Errorln(e)
				return func() {}
			}
			status := "待评审"
			switch d.Status {
			case 1:
				status = "待评审"
			case 2:
				status = "已批准"
			case -1:
				status = "已拒绝"
			case 9:
				status = "已关闭"
			}
			relatedIdList := []uint64{d.ApplyUserID}
			reviewerIDList := strings.Split(d.ReviewerIDList, ",")
			for _, reviewerIDStr := range reviewerIDList {
				reviewerID, _ := strconv.ParseUint(strings.TrimSpace(reviewerIDStr), 10, 64)
				if reviewerID > 0 {
					relatedIdList = append(relatedIdList, reviewerID)
				}
			}
			return data.AddDocumentAntsWrapper(&search.Document{
				ID:    d.ID,
				Title: d.Title,
				Content: fmt.Sprintf("原因:%s\n计划:%s\n评审:%s\n风险:%s\n状态:%s",
					d.Reason,
					d.Plan,
					d.Review,
					d.Risk,
					status,
				),
				Route:      fmt.Sprintf("/project/detail/%d/change?id=%d", d.ProjectID, d.ID),
				CreateTime: d.CreateTime,
				UpdateTime: d.UpdateTime,
			}, relatedIdList...)
		}(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectChangeController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectChange)
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

	// 先取出旧数据
	oldInstance, err := service.ProjectChangeService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	// 判断权限
	if _, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectChangeDelete); err != nil {
		return err
	}

	success, err := service.ProjectChangeService.Delete(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(data.DeleteDocumentsAntsWrapper(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectChangeController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectChangeListRequest)
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
	if instance.CreateTimeCondition != nil {
		tcList["create_time"] = instance.CreateTimeCondition
	}
	if instance.UpdateTimeCondition != nil {
		tcList["update_time"] = instance.UpdateTimeCondition
	}

	total, list, err := service.ProjectChangeService.PaginateBetweenTimes(&instance.ProjectChange, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectChangeController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectChange)
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
	dept, err := service.ProjectChangeService.Instance(condition.ID)
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
