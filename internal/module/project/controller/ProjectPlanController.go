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
)

var ProjectPlanController = new(projectPlanController)

type projectPlanController struct{}

func (c *projectPlanController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectPlan)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.PlanName == "" || instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name & projectId",
		})
	}

	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectPlanAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CreateUserID = uid
	}

	duplicated, success, err := service.ProjectPlanService.Add(instance)
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
		Title: instance.PlanName,
		Content: fmt.Sprintf("%s\n目标：%s\n范围：%s\n进度：%s\n资源%s",
			instance.PlanDesc,
			instance.Target,
			instance.Scope,
			instance.Schedule,
			instance.Resource,
		),
		Route:      fmt.Sprintf("/project/detail/%d/plan?id=%d", instance.ProjectID, instance.ID),
		CreateTime: instance.CreateTime,
		UpdateTime: instance.UpdateTime,
	}, instance.CreateUserID))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectPlanController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectPlan)
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

	// 先取出原来的数据
	old, err := service.ProjectPlanService.Instance(&model.ProjectPlan{ID: instance.ID})
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectPlanUpdate); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectPlanService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			d, e := service.ProjectPlanService.Instance(&model.ProjectPlan{ID: id})
			if e != nil {
				logger.Errorln(e)
				return func() {}
			}
			return data.AddDocumentAntsWrapper(&search.Document{
				ID:    d.ID,
				Title: d.PlanName,
				Content: fmt.Sprintf("%s\n目标：%s\n范围：%s\n进度：%s\n资源%s",
					instance.PlanDesc,
					instance.Target,
					instance.Scope,
					instance.Schedule,
					instance.Resource,
				),
				Route:      fmt.Sprintf("/project/detail/%d/plan?id=%d", d.ProjectID, d.ID),
				CreateTime: d.CreateTime,
				UpdateTime: d.UpdateTime,
			}, d.CreateUserID)
		}(instance.ID))
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectPlanController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectPlan)
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

	// 先取出原来的数据
	old, err := service.ProjectPlanService.Instance(&model.ProjectPlan{ID: instance.ID})
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectPlanDelete); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectPlanService.Delete(instance.ID)
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

func (c *projectPlanController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectPlanListRequest)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " project_id",
		})
	}

	if _, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectPlan); err != nil {
		return err
	} else if !success {
		return nil
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

	total, list, err := service.ProjectPlanService.PaginateBetweenTimes(&instance.ProjectPlan, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectPlanController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectPlan)
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
	dept, err := service.ProjectPlanService.Instance(condition)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if _, success, err := helper.CheckResourceCodeInProject(ctx, dept.ProjectID, constant.ResourceProjectPlanInstance); err != nil {
		return err
	} else if !success {
		return nil
	}
	return ctx.JSON(&server.CommonResponse{
		Data: dept,
	})
}

// ExecutingPlanByProject 通过项目ID获取正在执行的计划
func (c *projectPlanController) ExecutingPlanByProject(ctx *fiber.Ctx) error {
	condition := new(model.ProjectPlan)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if condition.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId",
		})
	}

	if _, success, err := helper.CheckResourceCodeInProject(ctx, condition.ProjectID, constant.ResourceProjectPlanInstance); err != nil {
		return err
	} else if !success {
		return nil
	}
	condition.Status = model.ProjectPlanStatusStarted
	plan, err := service.ProjectPlanService.Instance(condition)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: plan,
	})
}
