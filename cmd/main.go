package main

import (
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	projectController "github.com/yockii/celestial/internal/module/project/controller"
	ucController "github.com/yockii/celestial/internal/module/uc/controller"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/server"

	_ "github.com/yockii/celestial/initial"
	"github.com/yockii/ruomu-core/database"
)

func main() {
	defer ants.Release()
	config.InitialLogger()
	database.Initial()
	defer database.Close()

	//logger.Infoln("微核初始化完成")
	//
	//logger.Infoln("开始加载模块....")
	//logger.Infoln("加载模块管理")
	//_ = ruomu_module.Initial()
	//defer ruomu_module.Destroy()
	//logger.Infoln("模块管理加载完毕")

	// 统一用户中心模块
	ucController.InitRouter()
	// 项目管理模块
	projectController.InitRouter()

	for {
		err := server.Start()
		if err != nil {
			logger.Errorln(err)
		}
	}
}
