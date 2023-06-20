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

var ProjectRequirementController = new(projectRequirementController)

type projectRequirementController struct{}

func (c *projectRequirementController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Name == "" || instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name & projectId",
		})
	}

	if currentUserId, err := helper.GetCurrentUserID(ctx); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user info",
		})
	} else {
		instance.OwnerID = currentUserId
	}

	duplicated, success, err := service.ProjectRequirementService.Add(instance)
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

	_ = ants.Submit(func(d *model.ProjectRequirement) func() {
		//需求类型 1-功能 2-接口 3-性能 4-安全 5-体验 6-改进 7-其他
		t := "未知"
		switch d.Type {
		case 1:
			t = "功能"
		case 2:
			t = "接口"
		case 3:
			t = "性能"
		case 4:
			t = "安全"
		case 5:
			t = "体验"
		case 6:
			t = "改进"
		case 7:
			t = "其他"
		}
		//优先级 1-低 2-中 3-高
		p := "低"
		switch d.Priority {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		s := "内部"
		if d.Source == 1 {
			s = "客户"
		}
		f := "待评估"
		switch d.Feasibility {
		case -1:
			f = "不可行"
		case 1:
			f = "低"
		case 2:
			f = "中"
		case 3:
			f = "高"
		}
		status := "待设计"
		switch d.Status {
		case 1:
			status = "待设计"
		case 2:
			status = "待评审"
		case 3:
			status = "评审通过"
		case 9:
			status = "已完成"
		case -1:
			status = "评审未通过"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.Name,
			Content: fmt.Sprintf("[%s-%s]: %s, 可行性: %s, 状态: %s, 来源:%s",
				p,
				t,
				d.Detail,
				f,
				status,
				s,
			),
			Route:      fmt.Sprintf("/project/detail/%d/requirement?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, d.OwnerID)
	}(instance))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectRequirementController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
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

	success, err := service.ProjectRequirementService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		c.addSearchDocument(instance.ID)
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectRequirementController) addSearchDocument(id uint64) {
	_ = ants.Submit(func(id uint64) func() {
		d, e := service.ProjectRequirementService.Instance(id)
		if e != nil {
			logger.Errorln(e)
			return func() {
			}
		}
		//需求类型 1-功能 2-接口 3-性能 4-安全 5-体验 6-改进 7-其他
		t := "未知"
		switch d.Type {
		case 1:
			t = "功能"
		case 2:
			t = "接口"
		case 3:
			t = "性能"
		case 4:
			t = "安全"
		case 5:
			t = "体验"
		case 6:
			t = "改进"
		case 7:
			t = "其他"
		}
		//优先级 1-低 2-中 3-高
		p := "低"
		switch d.Priority {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		s := "内部"
		if d.Source == 1 {
			s = "客户"
		}
		f := "待评估"
		switch d.Feasibility {
		case -1:
			f = "不可行"
		case 1:
			f = "低"
		case 2:
			f = "中"
		case 3:
			f = "高"
		}
		status := "待设计"
		switch d.Status {
		case 1:
			status = "待设计"
		case 2:
			status = "待评审"
		case 3:
			status = "评审通过"
		case 9:
			status = "已完成"
		case -1:
			status = "评审未通过"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.Name,
			Content: fmt.Sprintf("[%s-%s]: %s, 可行性: %s, 状态: %s, 来源:%s",
				p,
				t,
				d.Detail,
				f,
				status,
				s,
			),
			Route:      fmt.Sprintf("/project/detail/%d/requirement?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, d.OwnerID)
	}(id))
}

func (c *projectRequirementController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
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

	success, err := service.ProjectRequirementService.Delete(instance.ID)
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

func (c *projectRequirementController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectRequirementListRequest)
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

	total, list, err := service.ProjectRequirementService.PaginateBetweenTimes(&instance.ProjectRequirement, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectRequirementController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectRequirement)
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
	dept, err := service.ProjectRequirementService.Instance(condition.ID)
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

// StatusDesigned 状态改为设计完成，即待评审
func (c *projectRequirementController) StatusDesigned(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	success, err := service.ProjectRequirementService.UpdateStatus(instance.ID, model.ProjectRequirementStatusPendingReview)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		c.addSearchDocument(instance.ID)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

// StatusReview 评审，状态可能改为评审通过或不通过
func (c *projectRequirementController) StatusReview(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 || instance.Status == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id & status",
		})
	}
	success, err := service.ProjectRequirementService.UpdateStatus(instance.ID, instance.Status)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if success {
		c.addSearchDocument(instance.ID)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

// StatusCompleted 状态改为已完成
func (c *projectRequirementController) StatusCompleted(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRequirement)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	success, err := service.ProjectRequirementService.UpdateStatus(instance.ID, model.ProjectRequirementStatusCompleted)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if success {
		c.addSearchDocument(instance.ID)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}
