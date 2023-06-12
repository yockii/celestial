package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {

	// 项目阶段信息
	{
		stage := server.Group("/api/v1/stage")
		stage.Post("/add", middleware.NeedAuthorization(constant.ResourceStageAdd), StageController.Add)
		stage.Delete("/delete", middleware.NeedAuthorization(constant.ResourceStageDelete), StageController.Delete)
		stage.Put("/update", middleware.NeedAuthorization(constant.ResourceStageUpdate), StageController.Update)
		stage.Get("/list", middleware.NeedAuthorization("user"), StageController.List)
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
		project.Get("/statisticsByStage", middleware.NeedAuthorization(constant.ResourceProjectList), ProjectController.StatisticsByStage)

		// 对于禁用put和delete方法时的处理
		project.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectDelete), ProjectController.Delete)
		project.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectUpdate), ProjectController.Update)
	}

	// 项目成员
	{
		projectMember := server.Group("/api/v1/projectMember")
		projectMember.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectMemberAdd), ProjectMemberController.Add)
		projectMember.Post("/batchAdd", middleware.NeedAuthorization(constant.ResourceProjectAdd), ProjectMemberController.BatchAdd)
		projectMember.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectMemberDelete), ProjectMemberController.Delete)
		projectMember.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectMemberUpdate), ProjectMemberController.Update)
		projectMember.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectMemberList), ProjectMemberController.List)
		projectMember.Get("/listByProject", middleware.NeedAuthorization(constant.ResourceProjectMemberList), ProjectMemberController.ListLiteByProjectId)
		//projectMember.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectMemberInstance), ProjectMemberController.Instance)

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
		projectPlan.Get("/executing", middleware.NeedAuthorization(constant.ResourceProjectPlanInstance), ProjectPlanController.ExecutingPlanByProject)

		// 对于禁用put和delete方法时的处理
		projectPlan.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), ProjectPlanController.Delete)
		projectPlan.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), ProjectPlanController.Update)
	}

	// 项目模块
	{
		projectModule := server.Group("/api/v1/projectModule")
		projectModule.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectModuleAdd), ProjectModuleController.Add)
		projectModule.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectModuleDelete), ProjectModuleController.Delete)
		projectModule.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectModuleUpdate), ProjectModuleController.Update)
		projectModule.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectModuleList), ProjectModuleController.List)
		//projectModule.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectModuleInstance), ProjectModuleController.Instance)
		// 评审状态修改
		projectModule.Put("/review", middleware.NeedAuthorization(constant.ResourceProjectModuleReview), ProjectModuleController.Review)

		// 对于禁用put和delete方法时的处理
		projectModule.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectModuleDelete), ProjectModuleController.Delete)
		projectModule.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectModuleUpdate), ProjectModuleController.Update)
		projectModule.Post("/review", middleware.NeedAuthorization(constant.ResourceProjectModuleReview), ProjectModuleController.Review)
	}

	// 项目需求
	{
		projectRequirement := server.Group("/api/v1/projectRequirement")
		projectRequirement.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectRequirementAdd), ProjectRequirementController.Add)
		projectRequirement.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectRequirementDelete), ProjectRequirementController.Delete)
		projectRequirement.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectRequirementUpdate), ProjectRequirementController.Update)
		projectRequirement.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectRequirementList), ProjectRequirementController.List)
		projectRequirement.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectRequirementInstance), ProjectRequirementController.Instance)
		// 3种状态修改
		projectRequirement.Put("/designed", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusDesign), ProjectRequirementController.StatusDesigned)
		projectRequirement.Put("/review", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusReview), ProjectRequirementController.StatusReview)
		projectRequirement.Put("/completed", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusCompleted), ProjectRequirementController.StatusCompleted)

		// 对于禁用put和delete方法时的处理
		projectRequirement.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectRequirementDelete), ProjectRequirementController.Delete)
		projectRequirement.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectRequirementUpdate), ProjectRequirementController.Update)
		projectRequirement.Post("/designed", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusDesign), ProjectRequirementController.StatusDesigned)
		projectRequirement.Post("/review", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusReview), ProjectRequirementController.StatusReview)
		projectRequirement.Post("/completed", middleware.NeedAuthorization(constant.ResourceProjectRequirementStatusCompleted), ProjectRequirementController.StatusCompleted)

	}

	//// 项目变更
	//{
	//	projectChange := server.Group("/api/v1/projectChange")
	//	projectChange.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectChangeAdd), ProjectChangeController.Add)
	//	projectChange.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectChangeDelete), ProjectChangeController.Delete)
	//	projectChange.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectChangeUpdate), ProjectChangeController.Update)
	//	projectChange.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectChangeList), ProjectChangeController.List)
	//	projectChange.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectChangeInstance), ProjectChangeController.Instance)
	//
	//	// 对于禁用put和delete方法时的处理
	//	projectChange.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectChangeDelete), ProjectChangeController.Delete)
	//	projectChange.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectChangeUpdate), ProjectChangeController.Update)
	//}
	//
	//// 项目问题
	//{
	//	projectIssue := server.Group("/api/v1/projectIssue")
	//	projectIssue.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectIssueAdd), ProjectIssueController.Add)
	//	projectIssue.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectIssueDelete), ProjectIssueController.Delete)
	//	projectIssue.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectIssueUpdate), ProjectIssueController.Update)
	//	projectIssue.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectIssueList), ProjectIssueController.List)
	//	projectIssue.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectIssueInstance), ProjectIssueController.Instance)
	//
	//	// 对于禁用put和delete方法时的处理
	//	projectIssue.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectIssueDelete), ProjectIssueController.Delete)
	//	projectIssue.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectIssueUpdate), ProjectIssueController.Update)
	//}
	//
	// 项目风险
	{
		projectRisk := server.Group("/api/v1/projectRisk")
		projectRisk.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectRiskAdd), ProjectRiskController.Add)
		projectRisk.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectRiskDelete), ProjectRiskController.Delete)
		projectRisk.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectRiskUpdate), ProjectRiskController.Update)
		projectRisk.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectRiskList), ProjectRiskController.List)
		projectRisk.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectRiskInstance), ProjectRiskController.Instance)
		projectRisk.Get("/coefficient", middleware.NeedAuthorization(constant.ResourceProjectInstance), ProjectRiskController.CalculateRiskByProject)

		// 对于禁用put和delete方法时的处理
		projectRisk.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectRiskDelete), ProjectRiskController.Delete)
		projectRisk.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectRiskUpdate), ProjectRiskController.Update)
	}
	//
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

}
