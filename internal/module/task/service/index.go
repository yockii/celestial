package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/task/constant"
	"github.com/yockii/celestial/internal/module/task/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

func InitService() {
	_ = database.AutoMigrate(model.Models...)
	// 初始化项目的资源
	var resources []*ucModel.Resource

	// 项目任务
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目任务",
			ResourceCode: constant.ResourceProjectTask,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目任务列表",
			ResourceCode: constant.ResourceProjectTaskList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目任务详情",
			ResourceCode: constant.ResourceProjectTaskInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目任务",
			ResourceCode: constant.ResourceProjectTaskAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目任务",
			ResourceCode: constant.ResourceProjectTaskUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目任务",
			ResourceCode: constant.ResourceProjectTaskDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目任务成员
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目任务成员",
			ResourceCode: constant.ResourceProjectTaskMember,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目任务成员列表",
			ResourceCode: constant.ResourceProjectTaskMemberList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目任务成员详情",
			ResourceCode: constant.ResourceProjectTaskMemberInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目任务成员",
			ResourceCode: constant.ResourceProjectTaskMemberAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目任务成员",
			ResourceCode: constant.ResourceProjectTaskMemberUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目任务成员",
			ResourceCode: constant.ResourceProjectTaskMemberDelete,
			HttpMethod:   "DELETE|POST",
		})
	}

	for _, resource := range resources {
		//没有就添加资源
		if err := database.DB.Where(resource).Attrs(&ucModel.Resource{
			ID: util.SnowflakeId(),
		}).FirstOrCreate(resource).Error; err != nil {
			logger.Errorln(err)
		}
	}
}
