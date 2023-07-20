package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	// OnlyOffice的路由
	// 下载
	//server.Get("/api/v1/office/download", middleware.NeedAuthorization(constant.NeedLogin), OnlyOfficeController.Download)

	// 保存
	//server.Post("/api/v1/office/save", middleware.NeedAuthorization(constant.NeedLogin), OnlyOfficeController.Save)

	server.Get("/api/v1/office/config", middleware.NeedAuthorization(constant.NeedLogin), OnlyOfficeController.GetConfig)
	server.Get("/api/v1/office/download", OnlyOfficeController.Download)
	server.Post("/api/v1/office/callback", OnlyOfficeController.Callback)

}
