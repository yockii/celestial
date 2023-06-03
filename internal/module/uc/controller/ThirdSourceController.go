package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/domain"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/server"
)

var ThirdSourceController = new(thirdSourceController)

type thirdSourceController struct{}

// Add 添加第三方登录源
func (c *thirdSourceController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ThirdSource)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	// 处理必填
	if instance.Name == "" || instance.Code == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " source name or source code",
		})
	}
	duplicated, success, err := service.ThirdSourceService.Add(instance)
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

// Update 更新第三方登录源
func (c *thirdSourceController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ThirdSource)
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
			Msg:  server.ResponseMsgParamNotEnough + " source id",
		})
	}
	success, err := service.ThirdSourceService.Update(instance)
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

// Delete 删除第三方登陆源
func (c *thirdSourceController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ThirdSource)
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
			Msg:  server.ResponseMsgParamNotEnough + " source id",
		})
	}
	success, err := service.ThirdSourceService.Delete(instance.ID)
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

// List 列出第三方登陆源
func (c *thirdSourceController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ThirdSourceListRequest)
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

	total, list, err := service.ThirdSourceService.PaginateBetweenTimes(&instance.ThirdSource, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

// Instance 获取第三方登陆源实例
func (c *thirdSourceController) Instance(ctx *fiber.Ctx) error {
	instance := new(model.ThirdSource)
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
			Msg:  server.ResponseMsgParamNotEnough + " source id",
		})
	}
	success, err := service.ThirdSourceService.Instance(instance)
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

func (c *thirdSourceController) PublicList(ctx *fiber.Ctx) error {
	list, err := service.ThirdSourceService.ListAll(&model.ThirdSource{
		Status: 1,
	})
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	var result []*domain.ThirdSourcePublic
	for _, item := range list {
		configJson := item.Configuration
		cj := gjson.Parse(configJson)
		appKey := ""
		if item.Code == "dingtalk" {
			appKey = cj.Get("appKey").String()
		} else {
			continue
		}

		tsp := &domain.ThirdSourcePublic{
			ID:     item.ID,
			Name:   item.Name,
			CorpID: item.CorpId,
			AppKey: appKey,
		}
		result = append(result, tsp)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}
