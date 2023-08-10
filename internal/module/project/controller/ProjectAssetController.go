package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/project/service"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	ucService "github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/server"
	"sync"
)

var ProjectAssetController = new(projectAssetController)

type projectAssetController struct{}

func (c *projectAssetController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectAsset)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Name == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name",
		})
	}

	// 判断权限
	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectAssetAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectAssetService.Add(instance)
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

func (c *projectAssetController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectAsset)
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
	oldInstance, err := service.ProjectAssetService.Instance(instance.ID)
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
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectAssetUpdate); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectAssetService.Update(instance)
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

func (c *projectAssetController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectAsset)
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
	oldInstance, err := service.ProjectAssetService.Instance(instance.ID)
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
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectAssetDelete); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectAssetService.Delete(instance.ID)
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

func (c *projectAssetController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectAssetListRequest)
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

	if _, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectAsset); err != nil {
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

	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " uid",
		})
	}

	total, list, err := service.ProjectAssetService.PaginateBetweenTimes(&instance.ProjectAsset, uid, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	var wg sync.WaitGroup
	for _, asset := range list {
		pa := asset
		wg.Add(1)
		go func(f *domain.ProjectAsset) {
			defer wg.Done()
			// 获取creator信息
			user, err := ucService.UserService.Instance(&ucModel.User{ID: f.CreatorID})
			if err != nil {
				logger.Errorln(err)
				return
			}
			f.Creator = &ucModel.User{
				ID:       user.ID,
				Username: user.Username,
				RealName: user.RealName,
			}
		}(pa)
	}
	wg.Wait()
	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  list,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

func (c *projectAssetController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectAsset)
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
	dept, err := service.ProjectAssetService.Instance(condition.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if _, success, err := helper.CheckResourceCodeInProject(ctx, dept.ProjectID, constant.ResourceProjectAssetInstance); err != nil {
		return err
	} else if !success {
		return nil
	}

	return ctx.JSON(&server.CommonResponse{
		Data: dept,
	})
}
