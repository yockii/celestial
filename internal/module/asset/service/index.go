package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/constant"
	"github.com/yockii/celestial/internal/module/asset/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

func InitService() {
	_ = database.AutoMigrate(model.Models...)
	// 初始化项目的资源
	var resources []*ucModel.Resource
	// 资产分类
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "资产分类",
			ResourceCode: constant.ResourceAssetCategory,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "资产分类添加",
			ResourceCode: constant.ResourceAssetCategoryAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "资产分类删除",
			ResourceCode: constant.ResourceAssetCategoryDelete,
			HttpMethod:   "DELETE|POST",
		}, &ucModel.Resource{
			ResourceName: "资产分类修改",
			ResourceCode: constant.ResourceAssetCategoryUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "资产分类列表",
			ResourceCode: constant.ResourceAssetCategoryList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "资产分类详情",
			ResourceCode: constant.ResourceAssetCategoryInstance,
			HttpMethod:   "GET",
		})
	}
	// 资产文件
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "资产文件",
			ResourceCode: constant.ResourceAssetFile,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "资产文件添加",
			ResourceCode: constant.ResourceAssetFileAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "资产文件删除",
			ResourceCode: constant.ResourceAssetFileDelete,
			HttpMethod:   "DELETE|POST",
		}, &ucModel.Resource{
			ResourceName: "资产文件修改",
			ResourceCode: constant.ResourceAssetFileUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "资产文件列表",
			ResourceCode: constant.ResourceAssetFileList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "资产文件详情",
			ResourceCode: constant.ResourceAssetFileInstance,
			HttpMethod:   "GET",
		})
	}

	// 对象存储配置
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "对象存储配置",
			ResourceCode: constant.ResourceOssConfig,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "对象存储配置添加",
			ResourceCode: constant.ResourceOssConfigAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "对象存储配置删除",
			ResourceCode: constant.ResourceOssConfigDelete,
			HttpMethod:   "DELETE|POST",
		}, &ucModel.Resource{
			ResourceName: "对象存储配置修改",
			ResourceCode: constant.ResourceOssConfigUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "对象存储配置列表",
			ResourceCode: constant.ResourceOssConfigList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "对象存储配置详情",
			ResourceCode: constant.ResourceOssConfigInstance,
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
