package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	// 日志信息
	{
		//log := server.Group("/api/v1/log")
		//log.Get("/list", middleware.NeedAuthorization(constant.ResourceLog), LogController.List)
		//log.Get("/instance", middleware.NeedAuthorization(constant.ResourceLogInstance), LogController.Instance)
	}

	// 工时信息
	{
		workTime := server.Group("/api/v1/workTime")
		workTime.Post("/add", middleware.NeedAuthorization(constant.ResourceWorkTimeAdd), WorkTimeController.Add)
		workTime.Delete("/delete", middleware.NeedAuthorization(constant.ResourceWorkTimeDelete), WorkTimeController.Delete)
		workTime.Put("/update", middleware.NeedAuthorization(constant.ResourceWorkTimeUpdate), WorkTimeController.Update)
		workTime.Get("/list", middleware.NeedAuthorization(constant.ResourceWorkTime), WorkTimeController.List)
		//workTime.Get("/instance", middleware.NeedAuthorization(constant.ResourceWorkTimeInstance), WorkTimeController.Instance)
		workTime.Get("/statistics", middleware.NeedAuthorization(constant.ResourceWorkTimeStatistics), WorkTimeController.Statistics)

		workTime.Get("/mine", middleware.NeedAuthorization(constant.ResourceWorkTime), WorkTimeController.Mine)

		// 对于禁用put和delete方法时的处理
		workTime.Post("/delete", middleware.NeedAuthorization(constant.ResourceWorkTimeDelete), WorkTimeController.Delete)
		workTime.Post("/update", middleware.NeedAuthorization(constant.ResourceWorkTimeUpdate), WorkTimeController.Update)
	}
}
