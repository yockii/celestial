package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/log/constant"
	"github.com/yockii/celestial/internal/module/log/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

func InitService() {
	_ = database.AutoMigrate(model.Models...)
	// 初始化项目的资源
	var resources []*ucModel.Resource
	// 日志
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "日志",
			ResourceCode: constant.ResourceLog,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "日志添加",
			ResourceCode: constant.ResourceLogAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "日志删除",
			ResourceCode: constant.ResourceLogDelete,
			HttpMethod:   "DELETE|POST",
		}, &ucModel.Resource{
			ResourceName: "日志修改",
			ResourceCode: constant.ResourceLogUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "日志列表",
			ResourceCode: constant.ResourceLogList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "日志详情",
			ResourceCode: constant.ResourceLogInstance,
			HttpMethod:   "GET",
		})
	}
	// 工时
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "工时",
			ResourceCode: constant.ResourceWorkTime,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "工时添加",
			ResourceCode: constant.ResourceWorkTimeAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "工时删除",
			ResourceCode: constant.ResourceWorkTimeDelete,
			HttpMethod:   "DELETE|POST",
		}, &ucModel.Resource{
			ResourceName: "工时修改",
			ResourceCode: constant.ResourceWorkTimeUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "工时列表",
			ResourceCode: constant.ResourceWorkTimeList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "工时详情",
			ResourceCode: constant.ResourceWorkTimeInstance,
			HttpMethod:   "GET",
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
