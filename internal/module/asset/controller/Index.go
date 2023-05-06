package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/asset/service"
	"github.com/yockii/celestial/internal/module/project/constant"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	service.InitService()

	// 资产分类
	{
		assetCategory := server.Group("/api/v1/assetCategory")
		assetCategory.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectPlanAdd), AssetCategoryController.Add)
		assetCategory.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), AssetCategoryController.Delete)
		assetCategory.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), AssetCategoryController.Update)
		assetCategory.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectPlanList), AssetCategoryController.List)
		assetCategory.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectPlanInstance), AssetCategoryController.Instance)

		// 对于禁用put和delete方法时的处理
		assetCategory.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), AssetCategoryController.Delete)
		assetCategory.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), AssetCategoryController.Update)
	}

	// 资产文件
	{
		assetFile := server.Group("/api/v1/assetFile")
		assetFile.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectPlanAdd), AssetFileController.Add)
		assetFile.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), AssetFileController.Delete)
		assetFile.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), AssetFileController.Update)
		assetFile.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectPlanList), AssetFileController.List)
		assetFile.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectPlanInstance), AssetFileController.Instance)

		// 对于禁用put和delete方法时的处理
		assetFile.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), AssetFileController.Delete)
		assetFile.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), AssetFileController.Update)
	}

	// 对象存储配置
	{
		ossConfig := server.Group("/api/v1/ossConfig")
		ossConfig.Post("/add", middleware.NeedAuthorization(constant.ResourceProjectPlanAdd), OssConfigController.Add)
		ossConfig.Delete("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), OssConfigController.Delete)
		ossConfig.Put("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), OssConfigController.Update)
		ossConfig.Get("/list", middleware.NeedAuthorization(constant.ResourceProjectPlanList), OssConfigController.List)
		ossConfig.Get("/instance", middleware.NeedAuthorization(constant.ResourceProjectPlanInstance), OssConfigController.Instance)

		// 对于禁用put和delete方法时的处理
		ossConfig.Post("/delete", middleware.NeedAuthorization(constant.ResourceProjectPlanDelete), OssConfigController.Delete)
		ossConfig.Post("/update", middleware.NeedAuthorization(constant.ResourceProjectPlanUpdate), OssConfigController.Update)
	}
}
