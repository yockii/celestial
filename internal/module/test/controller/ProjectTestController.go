package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/celestial/internal/module/test/service"
	"github.com/yockii/ruomu-core/server"
)

var ProjectTestController = new(projectTestController)

type projectTestController struct{}

func (c *projectTestController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId",
		})
	}

	// 判断权限
	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectTestAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CreatorID = uid
	}

	success, err := service.ProjectTestService.Add(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
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

func (c *projectTestController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
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

	// 先取出原始数据
	oldInstance, err := service.ProjectTestService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if oldInstance == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTestDelete); err != nil {
		return err
	}

	success, err = service.ProjectTestService.Delete(instance.ID)
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

func (c *projectTestController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
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

	// 先取出原始数据
	oldInstance, err := service.ProjectTestService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if oldInstance == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTestUpdate); err != nil {
		return err
	}

	success, err = service.ProjectTestService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
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

func (c *projectTestController) Close(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
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

	// 先取出原始数据
	oldInstance, err := service.ProjectTestService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if oldInstance == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var uid uint64
	var success bool
	if uid, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTestClose); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CloserID = uid
	}

	success, err = service.ProjectTestService.Close(instance.ID, instance.CloserID)
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

func (c *projectTestController) List(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
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
			Msg:  server.ResponseMsgParamNotEnough + " projectId",
		})
	}

	list, err := service.ProjectTestService.List(instance)
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

func (c *projectTestController) Instance(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTest)
	if err := ctx.QueryParser(instance); err != nil {
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

	instance, err := service.ProjectTestService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}
