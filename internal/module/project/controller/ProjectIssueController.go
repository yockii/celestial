package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/project/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
)

var ProjectIssueController = new(projectIssueController)

type projectIssueController struct{}

func (c *projectIssueController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectIssue)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ProjectID == 0 || instance.Title == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId & title",
		})
	}

	if uid, err := helper.GetCurrentUserID(ctx); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError + err.Error(),
		})
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectIssueService.Add(instance)
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

	_ = ants.Submit(data.AddDocumentAntsWrapper(&search.Document{
		ID:    instance.ID,
		Title: instance.Title,
		Content: fmt.Sprintf("%s\n原因:%s\n解决方式:%s",
			instance.Content,
			instance.IssueCause,
			instance.SolveMethod,
		),
		Route:      fmt.Sprintf("/project/detail/%d/issue?id=%d", instance.ProjectID, instance.ID),
		CreateTime: instance.CreateTime,
		UpdateTime: instance.UpdateTime,
	}, instance.CreatorID))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectIssueController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectIssue)
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

	success, err := service.ProjectIssueService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			d, e := service.ProjectIssueService.Instance(id)
			if e != nil {
				logger.Errorln(e)
				return func() {}
			}
			return data.UpdateDocumentAntsWrapper(&search.Document{
				ID:    instance.ID,
				Title: instance.Title,
				Content: fmt.Sprintf("%s\n原因:%s\n解决方式:%s",
					instance.Content,
					instance.IssueCause,
					instance.SolveMethod,
				),
				UpdateTime: instance.UpdateTime,
			}, d.CreatorID, d.AssigneeID)
		}(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectIssueController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectIssue)
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

	success, err := service.ProjectIssueService.Delete(instance.ID)
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

func (c *projectIssueController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectIssueListRequest)
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

	total, list, err := service.ProjectIssueService.PaginateBetweenTimes(&instance.ProjectIssue, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectIssueController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectIssue)
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
	dept, err := service.ProjectIssueService.Instance(condition.ID)
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

func (c *projectIssueController) Assign(ctx *fiber.Ctx) error {
	instance := new(model.ProjectIssue)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ID == 0 || instance.AssigneeID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id, assignee_id",
		})
	}

	success, err := service.ProjectIssueService.Assign(instance.ID, instance.AssigneeID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(data.UpdateDocumentAntsWrapper(&search.Document{ID: instance.ID}, instance.CreatorID, instance.AssigneeID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectIssueController) UpdateStatus(statusList ...uint8) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		instance := new(model.ProjectIssue)
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

		if len(statusList) == 0 {
			if instance.Status == model.ProjectIssueStatusAssigned {
				statusList = []uint8{model.ProjectIssueStatusAssigned}
			} else if instance.Status == model.ProjectIssueStatusResolved {
				statusList = []uint8{model.ProjectIssueStatusResolved}
			} else {
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeParamNotEnough,
					Msg:  server.ResponseMsgParamNotEnough + " status",
				})
			}
		}

		success, err := service.ProjectIssueService.UpdateStatus(instance.ID, statusList[0])
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}

		if success {
			_ = ants.Submit(data.UpdateDocumentAntsWrapper(&search.Document{ID: instance.ID}, instance.CreatorID, instance.AssigneeID))
		}

		return ctx.JSON(&server.CommonResponse{
			Data: success,
		})
	}
}
