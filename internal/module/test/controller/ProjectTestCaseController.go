package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/test/domain"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/celestial/internal/module/test/service"
	"github.com/yockii/ruomu-core/server"
	"sync"
)

var ProjectTestCaseController = new(projectTestCaseController)

type projectTestCaseController struct{}

func (c *projectTestCaseController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTestCase)
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

	// 判断权限
	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectTestCaseAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectTestCaseService.Add(instance)
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

func (c *projectTestCaseController) BatchAdd(ctx *fiber.Ctx) error {
	var list []*domain.ProjectTestCaseWithItems
	if err := ctx.BodyParser(&list); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if len(list) == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " test case list",
		})
	}

	// 判断权限
	var uid uint64
	var err error
	var success bool
	if uid, success, err = helper.CheckResourceCodeInProject(ctx, list[0].ProjectID, constant.ResourceProjectTestCaseAdd); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectTestCaseService.BatchAdd(list, uid)
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
		Data: success,
	})
}

func (c *projectTestCaseController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTestCase)
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
	oldInstance, err := service.ProjectTestCaseService.Instance(instance.ID)
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
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTestCaseUpdate); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectTestCaseService.Update(instance)
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

func (c *projectTestCaseController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTestCase)
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
	oldInstance, err := service.ProjectTestCaseService.Instance(instance.ID)
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
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTestCaseDelete); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectTestCaseService.Delete(instance.ID)
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

func (c *projectTestCaseController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectTestCaseListRequest)
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

	if _, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectTestCaseList); err != nil {
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

	total, list, err := service.ProjectTestCaseService.PaginateBetweenTimes(&instance.ProjectTestCase, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *projectTestCaseController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectTestCase)
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
	dept, err := service.ProjectTestCaseService.Instance(condition.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if _, success, err := helper.CheckResourceCodeInProject(ctx, dept.ProjectID, constant.ResourceProjectTestCaseInstance); err != nil {
		return err
	} else if !success {
		return nil
	}
	return ctx.JSON(&server.CommonResponse{
		Data: dept,
	})
}

func (c *projectTestCaseController) ListWithItems(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectTestCaseListRequest)
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

	if _, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectTestCaseList); err != nil {
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
	total, list, err := service.ProjectTestCaseService.PaginateBetweenTimes(&instance.ProjectTestCase, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.ProjectTestCaseWithItems
	var wg sync.WaitGroup
	for _, item := range list {
		wg.Add(1)
		result := &domain.ProjectTestCaseWithItems{
			ProjectTestCase: *item,
		}
		resultList = append(resultList, result)
		go func(ctcwi *domain.ProjectTestCaseWithItems) {
			defer wg.Done()
			_, itemList, err := service.ProjectTestCaseItemService.PaginateBetweenTimes(&model.ProjectTestCaseItem{
				TestCaseID: ctcwi.ID,
			}, -1, -1, "", nil)
			if err != nil {
				logger.Errorln(err)
				return
			}
			ctcwi.Items = itemList
		}(result)
	}
	wg.Wait()

	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  resultList,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}
