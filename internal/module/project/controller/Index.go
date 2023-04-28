package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/project/constant"
	"github.com/yockii/celestial/internal/module/project/service"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	service.InitService()

	// 项目阶段信息
	{
		stage := server.Group("/api/v1/stage")
		stage.Post("/add", middleware.NeedAuthorization(constant.ResourceStageAdd), StageController.Add)
		stage.Delete("/delete", middleware.NeedAuthorization(constant.ResourceStageDelete), StageController.Delete)
		stage.Put("/update", middleware.NeedAuthorization(constant.ResourceStageUpdate), StageController.Update)
		stage.Get("/list", middleware.NeedAuthorization(constant.ResourceStageList), StageController.List)
		stage.Get("/instance", middleware.NeedAuthorization(constant.ResourceStageInstance), StageController.Instance)

		// 对于禁用put和delete方法时的处理
		stage.Post("/delete", middleware.NeedAuthorization(constant.ResourceStageDelete), StageController.Delete)
		stage.Post("/update", middleware.NeedAuthorization(constant.ResourceStageUpdate), StageController.Update)
	}

	// 项目基础信息
	{
		project := server.Group("/api/v1/project")
		project.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectAdd), ProjectController.Add)
		project.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectDelete), ProjectController.Delete)
		project.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectUpdate), ProjectController.Update)
		project.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectList), ProjectController.List)
		project.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectInstance), ProjectController.Instance)

		// 对于禁用put和delete方法时的处理
		project.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectDelete), ProjectController.Delete)
		project.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectUpdate), ProjectController.Update)
	}

	// 项目成员
	{
		projectMember := server.Group("/api/v1/projectMember")
		projectMember.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectMemberAdd), ProjectMemberController.Add)
		projectMember.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectMemberDelete), ProjectMemberController.Delete)
		projectMember.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectMemberUpdate), ProjectMemberController.Update)
		projectMember.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectMemberList), ProjectMemberController.List)
		projectMember.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectMemberInstance), ProjectMemberController.Instance)

		// 对于禁用put和delete方法时的处理
		projectMember.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectMemberDelete), ProjectMemberController.Delete)
		projectMember.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectMemberUpdate), ProjectMemberController.Update)
	}

	// 项目计划
	{
		projectPlan := server.Group("/api/v1/projectPlan")
		projectPlan.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectPlanAdd), ProjectPlanController.Add)
		projectPlan.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), ProjectPlanController.Delete)
		projectPlan.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), ProjectPlanController.Update)
		projectPlan.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectPlanList), ProjectPlanController.List)
		projectPlan.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectPlanInstance), ProjectPlanController.Instance)

		// 对于禁用put和delete方法时的处理
		projectPlan.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), ProjectPlanController.Delete)
		projectPlan.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), ProjectPlanController.Update)
	}

	// 项目需求
	{
		projectRequirement := server.Group("/api/v1/projectRequirement")
		projectRequirement.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectRequirementAdd), ProjectRequirementController.Add)
		projectRequirement.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectRequirementDelete), ProjectRequirementController.Delete)
		projectRequirement.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectRequirementUpdate), ProjectRequirementController.Update)
		projectRequirement.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectRequirementList), ProjectRequirementController.List)
		projectRequirement.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectRequirementInstance), ProjectRequirementController.Instance)

		// 对于禁用put和delete方法时的处理
		projectRequirement.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectRequirementDelete), ProjectRequirementController.Delete)
		projectRequirement.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectRequirementUpdate), ProjectRequirementController.Update)
	}

	// 项目任务
	{
		projectTask := server.Group("/api/v1/projectTask")
		projectTask.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTaskAdd), ProjectTaskController.Add)
		projectTask.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskDelete), ProjectTaskController.Delete)
		projectTask.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskUpdate), ProjectTaskController.Update)
		projectTask.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTaskList), ProjectTaskController.List)
		projectTask.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTaskInstance), ProjectTaskController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTask.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskDelete), ProjectTaskController.Delete)
		projectTask.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskUpdate), ProjectTaskController.Update)
	}

	// 项目任务成员
	{
		projectTaskMember := server.Group("/api/v1/projectTaskMember")
		projectTaskMember.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberAdd), ProjectTaskMemberController.Add)
		projectTaskMember.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		projectTaskMember.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
		projectTaskMember.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberList), ProjectTaskMemberController.List)
		projectTaskMember.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberInstance), ProjectTaskMemberController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTaskMember.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		projectTaskMember.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
	}

	// 项目变更
	{
		projectChange := server.Group("/api/v1/projectChange")
		projectChange.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectChangeAdd), ProjectChangeController.Add)
		projectChange.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectChangeDelete), ProjectChangeController.Delete)
		projectChange.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectChangeUpdate), ProjectChangeController.Update)
		projectChange.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectChangeList), ProjectChangeController.List)
		projectChange.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectChangeInstance), ProjectChangeController.Instance)

		// 对于禁用put和delete方法时的处理
		projectChange.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectChangeDelete), ProjectChangeController.Delete)
		projectChange.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectChangeUpdate), ProjectChangeController.Update)
	}

	// 项目问题
	{
		projectIssue := server.Group("/api/v1/projectIssue")
		projectIssue.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectIssueAdd), ProjectIssueController.Add)
		projectIssue.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectIssueDelete), ProjectIssueController.Delete)
		projectIssue.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectIssueUpdate), ProjectIssueController.Update)
		projectIssue.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectIssueList), ProjectIssueController.List)
		projectIssue.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectIssueInstance), ProjectIssueController.Instance)

		// 对于禁用put和delete方法时的处理
		projectIssue.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectIssueDelete), ProjectIssueController.Delete)
		projectIssue.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectIssueUpdate), ProjectIssueController.Update)
	}

	// 项目风险
	{
		projectRisk := server.Group("/api/v1/projectRisk")
		projectRisk.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectRiskAdd), ProjectRiskController.Add)
		projectRisk.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectRiskDelete), ProjectRiskController.Delete)
		projectRisk.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectRiskUpdate), ProjectRiskController.Update)
		projectRisk.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectRiskList), ProjectRiskController.List)
		projectRisk.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectRiskInstance), ProjectRiskController.Instance)

		// 对于禁用put和delete方法时的处理
		projectRisk.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectRiskDelete), ProjectRiskController.Delete)
		projectRisk.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectRiskUpdate), ProjectRiskController.Update)
	}

	// 项目资产
	{
		projectAsset := server.Group("/api/v1/projectAsset")
		projectAsset.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectAssetAdd), ProjectAssetController.Add)
		projectAsset.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectAssetDelete), ProjectAssetController.Delete)
		projectAsset.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectAssetUpdate), ProjectAssetController.Update)
		projectAsset.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectAssetList), ProjectAssetController.List)
		projectAsset.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectAssetInstance), ProjectAssetController.Instance)

		// 对于禁用put和delete方法时的处理
		projectAsset.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectAssetDelete), ProjectAssetController.Delete)
		projectAsset.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectAssetUpdate), ProjectAssetController.Update)
	}

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
