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
		projectTest.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestAdd), ProjectTestController.Add)
		projectTest.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestDelete), ProjectTestController.Delete)
		projectTest.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestUpdate), ProjectTestController.Update)
		projectTest.Put("/close", middleware.NeedAuthorization(constant.ResourceProjectTestClose), ProjectTestController.Close)
		projectTest.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestList), ProjectTestController.List)
		projectTest.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestInstance), ProjectTestController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTest.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestDelete), ProjectTestController.Delete)
		projectTest.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestUpdate), ProjectTestController.Update)
		projectTest.Post("/close", middleware.NeedAuthorization(constant.ResourceProjectTestClose), ProjectTestController.Close)
	}
	// 项目测试用例
	{
		projectTest := server.Group("/api/v1/project/testCase")
		projectTest.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestCaseAdd), ProjectTestCaseController.Add)
		projectTest.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseDelete), ProjectTestCaseController.Delete)
		projectTest.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseUpdate), ProjectTestCaseController.Update)
		projectTest.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseList), ProjectTestCaseController.List)
		projectTest.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseInstance), ProjectTestCaseController.Instance)
		projectTest.Get("/listWithItems", middleware.NeedAuthorization(constant.ResourceProjectTestCaseList), ProjectTestCaseController.ListWithItems)

		// 对于禁用put和delete方法时的处理
		projectTest.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseDelete), ProjectTestCaseController.Delete)
		projectTest.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseUpdate), ProjectTestCaseController.Update)
	}

	// 项目测试用例项
	{
		projectTestCase := server.Group("/api/v1/project/testCaseItem")
		projectTestCase.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemAdd), ProjectTestCaseItemController.Add)
		projectTestCase.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemDelete), ProjectTestCaseItemController.Delete)
		projectTestCase.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemUpdate), ProjectTestCaseItemController.Update)
		projectTestCase.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemList), ProjectTestCaseItemController.List)
		projectTestCase.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemInstance), ProjectTestCaseItemController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCase.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemDelete), ProjectTestCaseItemController.Delete)
		projectTestCase.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemUpdate), ProjectTestCaseItemController.Update)
	}

	// 项目测试用例项步骤
	{
		projectTestCaseStep := server.Group("/api/v1/project/testCaseItemStep")
		projectTestCaseStep.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepAdd), ProjectTestCaseItemStepController.Add)
		projectTestCaseStep.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepDelete), ProjectTestCaseItemStepController.Delete)
		projectTestCaseStep.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepUpdate), ProjectTestCaseItemStepController.Update)
		projectTestCaseStep.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepList), ProjectTestCaseItemStepController.List)
		projectTestCaseStep.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepInstance), ProjectTestCaseItemStepController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCaseStep.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepDelete), ProjectTestCaseItemStepController.Delete)
		projectTestCaseStep.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseItemStepUpdate), ProjectTestCaseItemStepController.Update)
	}
}
