package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/server"
)

var SysConfigController = new(sysConfigController)

type sysConfigController struct{}

// updatableSysConfigs 允许通过接口修改的配置键白名单
var updatableSysConfigs = map[string]string{
	constant.SysConfigKeyUsernamePasswordEnabled: "是否启用用户名密码登录/注册",
}
func (c *sysConfigController) List(ctx *fiber.Ctx) error {
	list, err := service.SysConfigService.List()
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

func (c *sysConfigController) Update(ctx *fiber.Ctx) error {
	instance := new(model.SysConfig)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	comment, has := updatableSysConfigs[instance.Key]
	if !has {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgDataNotMatch + ": 不允许修改的配置项",
		})
	}
	if instance.Value != "true" && instance.Value != "false" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + ": 值只能为 true/false",
		})
	}
	if err := service.SysConfigService.Set(instance.Key, instance.Value, comment); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: true,
	})
}
