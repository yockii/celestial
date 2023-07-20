package main

import (
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/data"
	assetController "github.com/yockii/celestial/internal/module/asset/controller"
	logController "github.com/yockii/celestial/internal/module/log/controller"
	luceneController "github.com/yockii/celestial/internal/module/lucene/controller"
	meetingRoomController "github.com/yockii/celestial/internal/module/meeting/controller"
	"github.com/yockii/celestial/internal/module/message"
	onlyofficeController "github.com/yockii/celestial/internal/module/onlyoffice/controller"
	projectController "github.com/yockii/celestial/internal/module/project/controller"
	taskController "github.com/yockii/celestial/internal/module/task/controller"
	testController "github.com/yockii/celestial/internal/module/test/controller"
	ucController "github.com/yockii/celestial/internal/module/uc/controller"
	"github.com/yockii/celestial/pkg/search"
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

	// 初始化搜索引擎
	if err := search.InitMeiliSearch(config.GetString("meilisearch.host"), config.GetString("meilisearch.apiKey"), config.GetString("meilisearch.index")); err != nil {
		logger.Warnln(err)
	}

	// 初始化数据
	data.InitData()

	// 订阅消息
	message.InitDingtalkMessageAdapter()

	// 开启全文检索
	luceneController.InitRouter()
	// 统一用户中心模块
	ucController.InitRouter()
	// 项目管理模块
	projectController.InitRouter()
	// 任务管理模块
	taskController.InitRouter()
	// 测试管理模块
	testController.InitRouter()
	// 日志管理模块
	logController.InitRouter()
	// 资产管理模块
	assetController.InitRouter()
	// OnlyOffice在线文档编辑模块
	onlyofficeController.InitRouter()
	// 会议室管理
	meetingRoomController.InitRouter()

	for {
		err := server.Start()
		if err != nil {
			logger.Errorln(err)
		}
	}
}
