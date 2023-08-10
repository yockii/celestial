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

var ProjectModuleController = new(projectModuleController)

type projectModuleController struct{}

func (c *projectModuleController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectModule)
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
			Msg:  server.ResponseMsgParamNotEnough + " projectID/name",
		})
	}

	// 判断权限
	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectModuleAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectModuleService.Add(instance)
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
		ID:         instance.ID,
		Title:      instance.Name,
		Content:    instance.Remark + ", \n状态：待评审",
		Route:      fmt.Sprintf("/project/detail/%d/module?id=%d", instance.ProjectID, instance.ID),
		CreateTime: instance.CreateTime,
		UpdateTime: instance.UpdateTime,
	}, instance.CreatorID))
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectModuleController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectModule)
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

	// 先取出
	old, err := service.ProjectModuleService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectModuleUpdate); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectModuleService.Update(instance, old)
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

func (c *projectModuleController) addSearchDocument(id uint64) {
	_ = ants.Submit(func(id uint64) func() {
		d, e := service.ProjectModuleService.Instance(id)
		if e != nil {
			logger.Errorln(e)
			return func() {}
		}
		//状态 1-待评审 2-评审通过待开发 9-已完成 -1-评审不通过
		status := "待评审"
		switch d.Status {
		case 1:
			status = "待评审"
		case 2:
			status = "评审通过待开发"
		case 9:
			status = "已完成"
		case -1:
			status = "评审不通过"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:         d.ID,
			Title:      d.Name,
			Content:    d.Remark + ", \n状态: " + status,
			Route:      fmt.Sprintf("/project/detail/%d/module?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, d.CreatorID)
	}(id))
}

func (c *projectModuleController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectModule)
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

	// 先取出
	old, err := service.ProjectModuleService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectModuleDelete); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectModuleService.Delete(old)
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

func (c *projectModuleController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectModuleListRequest)
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

	if _, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectModule); err != nil {
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

	total, list, err := service.ProjectModuleService.PaginateBetweenTimes(&instance.ProjectModule, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectModuleController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectModule)
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
	dept, err := service.ProjectModuleService.Instance(condition.ID)
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

// Review 项目模块评审
func (c *projectModuleController) Review(ctx *fiber.Ctx) error {
	instance := new(model.ProjectModule)
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

	// 先取出
	old, err := service.ProjectModuleService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectModuleReview); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectModuleService.UpdateStatus(old, instance.Status)
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
