package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {

	// 项目任务
	{
		projectTask := server.Group("/api/v1/projectTask")
		projectTask.Post("/add", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.Add)
		projectTask.Delete("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.Delete)
		projectTask.Put("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.Update)
		projectTask.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTaskList), ProjectTaskController.List)
		projectTask.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTaskInstance), ProjectTaskController.Instance)
		projectTask.Get("/statisticsByProject", middleware.NeedAuthorization(constant.ResourceProjectInstance), ProjectTaskController.TaskDurationByProject)
		projectTask.Put("/cancel", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusCancel))
		projectTask.Put("/confirm", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusConfirmed))
		projectTask.Put("/start", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDoing))
		projectTask.Put("/devDone", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDevDone))
		projectTask.Put("/testReject", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTestReject))
		projectTask.Put("/testing", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTesting))
		projectTask.Put("/testPass", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTestPass))
		projectTask.Put("/done", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusDone))
		projectTask.Put("/restart", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusNotStart))
		projectTask.Get("/listMine", middleware.NeedAuthorization(constant.ResourceProjectTaskList), ProjectTaskController.ListMine)

		// 对于禁用put和delete方法时的处理
		projectTask.Post("/delete", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.Delete)
		projectTask.Post("/update", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.Update)
		projectTask.Post("/cancel", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusCancel))
		projectTask.Post("/confirm", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusConfirmed))
		projectTask.Post("/start", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDoing))
		projectTask.Post("/devDone", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDevDone))
		projectTask.Post("/testReject", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTestReject))
		projectTask.Post("/testing", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTesting))
		projectTask.Post("/testPass", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusTestPass))
		projectTask.Post("/done", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusDone))
		projectTask.Post("/restart", middleware.NeedAuthorization(constant.NeedLogin), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusNotStart))
	}

	// 项目任务成员
	{
		//projectTaskMember := server.Group("/api/v1/projectTaskMember")
		//projectTaskMember.Post("/add", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberAdd), ProjectTaskMemberController.Add)
		//projectTaskMember.Delete("/delete", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		//projectTaskMember.Put("/update", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
		//projectTaskMember.Get("/list", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberList), ProjectTaskMemberController.List)
		//projectTaskMember.Get("/instance", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberInstance), ProjectTaskMemberController.Instance)
		//
		//// 对于禁用put和delete方法时的处理
		//projectTaskMember.Post("/delete", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		//projectTaskMember.Post("/update", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
	}
}
