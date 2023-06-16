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
	"sync"
)

var CommonTestCaseController = new(commonTestCaseController)

type commonTestCaseController struct{}

func (c *commonTestCaseController) Add(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCase)
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

	duplicated, success, err := service.CommonTestCaseService.Add(instance)
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
		Content:    instance.Remark,
		Route:      fmt.Sprintf("/asset/testcaselib?id=%d", instance.ID),
		CreateTime: instance.CreateTime,
		UpdateTime: instance.UpdateTime,
	}, instance.CreatorID))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *commonTestCaseController) Update(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCase)
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

	success, err := service.CommonTestCaseService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			return func() {
				d, e := service.CommonTestCaseService.Instance(id)
				if e != nil {
					logger.Errorln(err)
					return
				}
				_ = data.AddDocument(&search.Document{
					ID:         d.ID,
					Title:      d.Name,
					Content:    d.Remark,
					Route:      fmt.Sprintf("/asset/testcaselib?id=%d", d.ID),
					CreateTime: d.CreateTime,
					UpdateTime: d.UpdateTime,
				}, d.CreatorID)
			}
		}(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *commonTestCaseController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.CommonTestCase)
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

	success, err := service.CommonTestCaseService.Delete(instance.ID)
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

func (c *commonTestCaseController) List(ctx *fiber.Ctx) error {
	instance := new(domain.CommonTestCaseListRequest)
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

	total, list, err := service.CommonTestCaseService.PaginateBetweenTimes(&instance.CommonTestCase, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
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

func (c *commonTestCaseController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.CommonTestCase)
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
	dept, err := service.CommonTestCaseService.Instance(condition.ID)
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

func (c *commonTestCaseController) ListWithItem(ctx *fiber.Ctx) error {
	instance := new(domain.CommonTestCaseListRequest)
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

	total, list, err := service.CommonTestCaseService.PaginateBetweenTimes(&instance.CommonTestCase, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.CommonTestCaseWithItem
	var wg sync.WaitGroup
	for _, item := range list {
		wg.Add(1)
		result := &domain.CommonTestCaseWithItem{
			CommonTestCase: *item,
		}
		resultList = append(resultList, result)
		go func(ctcwi *domain.CommonTestCaseWithItem) {
			defer wg.Done()
			_, itemList, err := service.CommonTestCaseItemService.PaginateBetweenTimes(&model.CommonTestCaseItem{
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

func (c *commonTestCaseController) ListWithItemOnlyShow(ctx *fiber.Ctx) error {
	list, err := service.CommonTestCaseService.ListAllOnlyShow()
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.CommonTestCaseWithItem
	var wg sync.WaitGroup
	for _, item := range list {
		wg.Add(1)
		result := &domain.CommonTestCaseWithItem{
			CommonTestCase: *item,
		}
		resultList = append(resultList, result)
		go func(ctcwi *domain.CommonTestCaseWithItem) {
			defer wg.Done()
			itemList, err := service.CommonTestCaseItemService.ListAllOnlyShow(&model.CommonTestCaseItem{
				TestCaseID: ctcwi.ID,
			})
			if err != nil {
				logger.Errorln(err)
				return
			}
			ctcwi.Items = itemList
		}(result)
	}
	wg.Wait()

	return ctx.JSON(&server.CommonResponse{
		Data: resultList,
	})
}
