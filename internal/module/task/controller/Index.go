package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	taskConstant "github.com/yockii/celestial/internal/module/task/constant"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {

	// 项目任务
	{
		projectTask := server.Group("/api/v1/projectTask")
		projectTask.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectTaskAdd), ProjectTaskController.Add)
		projectTask.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskDelete), ProjectTaskController.Delete)
		projectTask.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskUpdate), ProjectTaskController.Update)
		projectTask.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectTaskList), ProjectTaskController.List)
		projectTask.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectTaskInstance), ProjectTaskController.Instance)
		projectTask.Get("/statisticsByProject", middleware.NeedAuthorization(constant.ResourceProjectInstance), ProjectTaskController.TaskDurationByProject)
		projectTask.Put("/cancel", middleware.NeedAuthorization(constant.ResourceProjectTaskCancel), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusCancel))
		projectTask.Put("/confirm", middleware.NeedAuthorization(constant.ResourceProjectTaskConfirm), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusConfirmed))
		projectTask.Put("/start", middleware.NeedAuthorization(constant.ResourceProjectTaskStart), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDoing))
		projectTask.Put("/done", middleware.NeedAuthorization(constant.ResourceProjectTaskDone), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDone))
		projectTask.Put("/restart", middleware.NeedAuthorization(constant.ResourceProjectTaskRestart), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusNotStart))

		// 对于禁用put和delete方法时的处理
		projectTask.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectTaskDelete), ProjectTaskController.Delete)
		projectTask.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectTaskUpdate), ProjectTaskController.Update)
		projectTask.Post("/cancel", middleware.NeedAuthorization(constant.ResourceProjectTaskCancel), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusCancel))
		projectTask.Post("/confirm", middleware.NeedAuthorization(constant.ResourceProjectTaskConfirm), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusConfirmed))
		projectTask.Post("/start", middleware.NeedAuthorization(constant.ResourceProjectTaskStart), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDoing))
		projectTask.Post("/done", middleware.NeedAuthorization(constant.ResourceProjectTaskDone), ProjectTaskController.MemberUpdateStatus(model.ProjectTaskStatusDone))
		projectTask.Post("/restart", middleware.NeedAuthorization(constant.ResourceProjectTaskRestart), ProjectTaskController.UpdateStatus(model.ProjectTaskStatusNotStart))
	}

	// 项目任务成员
	{
		projectTaskMember := server.Group("/api/v1/projectTaskMember")
		projectTaskMember.Post("/add", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberAdd), ProjectTaskMemberController.Add)
		projectTaskMember.Delete("/delete", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		projectTaskMember.Put("/update", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
		projectTaskMember.Get("/list", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberList), ProjectTaskMemberController.List)
		projectTaskMember.Get("/instance", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberInstance), ProjectTaskMemberController.Instance)

		// 对于禁用put和delete方法时的处理
		projectTaskMember.Post("/delete", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberDelete), ProjectTaskMemberController.Delete)
		projectTaskMember.Post("/update", middleware.NeedAuthorization(taskConstant.ResourceProjectTaskMemberUpdate), ProjectTaskMemberController.Update)
	}
}
