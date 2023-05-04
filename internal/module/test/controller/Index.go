package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/test/constant"
	"github.com/yockii/celestial/internal/module/test/service"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	service.InitService()

	// 项目测试
	{
		projectTest := server.Group("/api/v1/projectTest")
		projectTest.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestAdd), ProjectTestController.Add)
		projectTest.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestDelete), ProjectTestController.Delete)
		projectTest.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestUpdate), ProjectTestController.Update)
		projectTest.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestList), ProjectTestController.List)
		projectTest.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestInstance), ProjectTestController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTest.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestDelete), ProjectTestController.Delete)
		projectTest.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestUpdate), ProjectTestController.Update)
	}

	// 项目测试用例
	{
		projectTestCase := server.Group("/api/v1/projectTestCase")
		projectTestCase.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestCaseAdd), ProjectTestCaseController.Add)
		projectTestCase.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseDelete), ProjectTestCaseController.Delete)
		projectTestCase.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseUpdate), ProjectTestCaseController.Update)
		projectTestCase.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseList), ProjectTestCaseController.List)
		projectTestCase.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseInstance), ProjectTestCaseController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCase.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseDelete), ProjectTestCaseController.Delete)
		projectTestCase.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseUpdate), ProjectTestCaseController.Update)
	}

	// 项目测试用例步骤
	{
		projectTestCaseStep := server.Group("/api/v1/projectTestCaseStep")
		projectTestCaseStep.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepAdd), ProjectTestCaseStepController.Add)
		projectTestCaseStep.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepDelete), ProjectTestCaseStepController.Delete)
		projectTestCaseStep.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepUpdate), ProjectTestCaseStepController.Update)
		projectTestCaseStep.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepList), ProjectTestCaseStepController.List)
		projectTestCaseStep.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepInstance), ProjectTestCaseStepController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTestCaseStep.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepDelete), ProjectTestCaseStepController.Delete)
		projectTestCaseStep.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTestCaseStepUpdate), ProjectTestCaseStepController.Update)
	}
}
