package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	// 项目测试封版
	{
		projectTest := server.Group("/api/v1/project/test")
		projectTest.Post("/add", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Add)
		projectTest.Delete("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Delete)
		projectTest.Put("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Update)
		projectTest.Put("/close", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Close)
		projectTest.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestList, constant.ResourceAllProjectDetail), ProjectTestController.List)
		projectTest.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestInstance, constant.ResourceAllProjectDetail), ProjectTestController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTest.Post("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Delete)
		projectTest.Post("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Update)
		projectTest.Post("/close", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestController.Close)
	}
	// 项目测试用例
	{
		projectTest := server.Group("/api/v1/project/testCase")
		projectTest.Post("/add", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.Add)
		projectTest.Post("/batchAdd", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.BatchAdd)
		projectTest.Delete("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.Delete)
		projectTest.Put("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.Update)
		projectTest.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseList, constant.ResourceAllProjectDetail), ProjectTestCaseController.List)
		projectTest.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseInstance, constant.ResourceAllProjectDetail), ProjectTestCaseController.Instance)
		projectTest.Get("/listWithItems", middleware.NeedAuthorization(constant.ResourceProjectTestCaseList, constant.ResourceAllProjectDetail), ProjectTestCaseController.ListWithItems)

		// 对于禁用put和delete方法时的处理
		projectTest.Post("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.Delete)
		projectTest.Post("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseController.Update)
	}

	// 项目测试用例项
	{
		projectTestCase := server.Group("/api/v1/project/testCaseItem")
		projectTestCase.Post("/add", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.Add)
		projectTestCase.Post("/batchAdd", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.BatchAdd)
		projectTestCase.Delete("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.Delete)
		projectTestCase.Put("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.Update)
		projectTestCase.Put("/updateStatus", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.UpdateStatus)
		projectTestCase.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemList, constant.ResourceAllProjectDetail), ProjectTestCaseItemController.List)
		projectTestCase.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemInstance, constant.ResourceAllProjectDetail), ProjectTestCaseItemController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCase.Post("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.Delete)
		projectTestCase.Post("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemController.Update)
	}

	// 项目测试用例项步骤
	{
		projectTestCaseStep := server.Group("/api/v1/project/testCaseItemStep")
		projectTestCaseStep.Post("/add", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemStepController.Add)
		projectTestCaseStep.Delete("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemStepController.Delete)
		projectTestCaseStep.Put("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemStepController.Update)
		projectTestCaseStep.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepList, constant.ResourceAllProjectDetail), ProjectTestCaseItemStepController.List)
		projectTestCaseStep.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepInstance, constant.ResourceAllProjectDetail), ProjectTestCaseItemStepController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCaseStep.Post("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemStepController.Delete)
		projectTestCaseStep.Post("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTestCaseItemStepController.Update)
	}
}
