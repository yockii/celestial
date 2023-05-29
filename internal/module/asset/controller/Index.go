package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/asset/constant"
	"github.com/yockii/celestial/internal/module/asset/service"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	service.InitService()

	// 资产分类
	{
		assetCategory := server.Group("/api/v1/assetCategory")
		assetCategory.Post("/add", middleware.NeedAuthorization(constant.ResourceAssetCategory), AssetCategoryController.Add)
		assetCategory.Delete("/delete", middleware.NeedAuthorization(constant.ResourceAssetCategoryDelete), AssetCategoryController.Delete)
		assetCategory.Put("/update", middleware.NeedAuthorization(constant.ResourceAssetCategoryUpdate), AssetCategoryController.Update)
		assetCategory.Get("/list", middleware.NeedAuthorization(constant.ResourceAssetCategoryList), AssetCategoryController.List)
		assetCategory.Get("/instance", middleware.NeedAuthorization(constant.ResourceAssetCategoryInstance), AssetCategoryController.Instance)

		// 对于禁用put和delete方法时的处理
		assetCategory.Post("/delete", middleware.NeedAuthorization(constant.ResourceAssetCategoryDelete), AssetCategoryController.Delete)
		assetCategory.Post("/update", middleware.NeedAuthorization(constant.ResourceAssetCategoryUpdate), AssetCategoryController.Update)
	}

	// 资产文件
	{
		assetFile := server.Group("/api/v1/assetFile")
		assetFile.Post("/add", middleware.NeedAuthorization(constant.ResourceFileAdd), AssetFileController.Add)
		assetFile.Delete("/delete", middleware.NeedAuthorization(constant.ResourceFileDelete), AssetFileController.Delete)
		assetFile.Put("/update", middleware.NeedAuthorization(constant.ResourceFileUpdate), AssetFileController.Update)
		assetFile.Get("/list", middleware.NeedAuthorization(constant.ResourceFileList), AssetFileController.List)
		assetFile.Get("/instance", middleware.NeedAuthorization(constant.ResourceFileInstance), AssetFileController.Instance)

		// 对于禁用put和delete方法时的处理
		assetFile.Post("/delete", middleware.NeedAuthorization(constant.ResourceFileDelete), AssetFileController.Delete)
		assetFile.Post("/update", middleware.NeedAuthorization(constant.ResourceFileUpdate), AssetFileController.Update)
	}

	// 对象存储配置
	{
		ossConfig := server.Group("/api/v1/ossConfig")
		ossConfig.Post("/add", middleware.NeedAuthorization(constant.ResourceOssConfigAdd), OssConfigController.Add)
		ossConfig.Delete("/delete", middleware.NeedAuthorization(constant.ResourceOssConfigDelete), OssConfigController.Delete)
		ossConfig.Put("/update", middleware.NeedAuthorization(constant.ResourceOssConfigUpdate), OssConfigController.Update)
		ossConfig.Get("/list", middleware.NeedAuthorization(constant.ResourceOssConfigList), OssConfigController.List)
		ossConfig.Get("/instance", middleware.NeedAuthorization(constant.ResourceOssConfigInstance), OssConfigController.Instance)

		// 对于禁用put和delete方法时的处理
		ossConfig.Post("/delete", middleware.NeedAuthorization(constant.ResourceOssConfigDelete), OssConfigController.Delete)
		ossConfig.Post("/update", middleware.NeedAuthorization(constant.ResourceOssConfigUpdate), OssConfigController.Update)
	}

	// 通用测试用例
	{
		commonTestCase := server.Group("/api/v1/commonTestCase")
		commonTestCase.Post("/add", middleware.NeedAuthorization(constant.ResourceCommonTestCaseAdd), CommonTestCaseController.Add)
		commonTestCase.Delete("/delete", middleware.NeedAuthorization(constant.ResourceCommonTestCaseDelete), CommonTestCaseController.Delete)
		commonTestCase.Put("/update", middleware.NeedAuthorization(constant.ResourceCommonTestCaseUpdate), CommonTestCaseController.Update)
		commonTestCase.Get("/list", middleware.NeedAuthorization(constant.ResourceCommonTestCaseList), CommonTestCaseController.List)

		// 对于禁用put和delete方法时的处理
		commonTestCase.Post("/delete", middleware.NeedAuthorization(constant.ResourceCommonTestCaseDelete), CommonTestCaseController.Delete)
		commonTestCase.Post("/update", middleware.NeedAuthorization(constant.ResourceCommonTestCaseUpdate), CommonTestCaseController.Update)
	}

	// 通用测试用例项
	{
		commonTestCaseItem := server.Group("/api/v1/commonTestCaseItem")
		commonTestCaseItem.Post("/add", middleware.NeedAuthorization(constant.ResourceCommonTestCaseItemAdd), CommonTestCaseItemController.Add)
		commonTestCaseItem.Delete("/delete", middleware.NeedAuthorization(constant.ResourceCommonTestCaseItemDelete), CommonTestCaseItemController.Delete)
		commonTestCaseItem.Get("/list", middleware.NeedAuthorization(constant.ResourceCommonTestCaseItemList), CommonTestCaseItemController.List)

		// 对于禁用put和delete方法时的处理
		commonTestCaseItem.Post("/delete", middleware.NeedAuthorization(constant.ResourceCommonTestCaseItemDelete), CommonTestCaseItemController.Delete)
	}
}
