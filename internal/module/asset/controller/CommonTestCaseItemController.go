package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/module/asset/domain"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/celestial/internal/module/asset/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
)

var CommonTestCaseItemController = new(commonTestCaseItemController)

type commonTestCaseItemController struct{}

func (c *commonTestCaseItemController) Add(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCaseItem)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Content == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " content",
		})
	}

	duplicated, success, err := service.CommonTestCaseItemService.Add(instance)
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

	c.addSearchDocument(instance.TestCaseID)

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *commonTestCaseItemController) addSearchDocument(testCaseID uint64) {
	_ = ants.Submit(func(parentId uint64) func() {
		// 获取父级
		testCase, err := service.CommonTestCaseService.Instance(parentId)
		if err != nil {
			logger.Errorln(err)
			return func() {}
		}
		_, items, err := service.CommonTestCaseItemService.PaginateBetweenTimes(&model.CommonTestCaseItem{TestCaseID: parentId}, -1, -1, "", nil)
		if err != nil {
			logger.Errorln(err)
			return func() {}
		}
		content := testCase.Remark
		relatedUidMap := map[uint64]struct{}{}
		relatedUidMap[testCase.CreatorID] = struct{}{}
		updateTime := testCase.UpdateTime
		if len(items) > 0 {
			for _, item := range items {
				content += "\n" + item.Content
				relatedUidMap[item.CreatorID] = struct{}{}
				if item.UpdateTime > (updateTime) {
					updateTime = item.UpdateTime
				}
			}
		}
		var relatedUidList []uint64
		for uid, _ := range relatedUidMap {
			relatedUidList = append(relatedUidList, uid)
		}

		return func() {
			_ = data.AddDocument(&search.Document{
				ID:         testCase.ID,
				Title:      testCase.Name,
				Content:    content,
				Route:      fmt.Sprintf("/asset/testcase?id=%d", testCase.ID),
				CreateTime: testCase.CreateTime,
				UpdateTime: updateTime,
			}, relatedUidList...)
		}
	}(testCaseID))
}

func (c *commonTestCaseItemController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCaseItem)
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

	success, err := service.CommonTestCaseItemService.Delete(instance.ID)
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

func (c *commonTestCaseItemController) Update(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCaseItem)
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

	success, err := service.CommonTestCaseItemService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		instance, err = service.CommonTestCaseItemService.Instance(instance.ID)
		if err != nil {
			logger.Errorln(err)
		} else {
			c.addSearchDocument(instance.TestCaseID)
		}
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *commonTestCaseItemController) List(ctx *fiber.Ctx) error {
	instance := new(domain.CommonTestCaseItemListRequest)
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

	total, list, err := service.CommonTestCaseItemService.PaginateBetweenTimes(&instance.CommonTestCaseItem, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *commonTestCaseItemController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.CommonTestCaseItem)
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
	dept, err := service.CommonTestCaseItemService.Instance(condition.ID)
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
