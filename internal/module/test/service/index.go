package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/test/constant"
	"github.com/yockii/celestial/internal/module/test/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

func InitService() {
	_ = database.AutoMigrate(model.Models...)
	// 初始化项目的资源
	var resources []*ucModel.Resource
	// 项目测试
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目测试",
			ResourceCode: constant.ResourceProjectTest,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目测试列表",
			ResourceCode: constant.ResourceProjectTestList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目测试详情",
			ResourceCode: constant.ResourceProjectTestInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目测试",
			ResourceCode: constant.ResourceProjectTestAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目测试",
			ResourceCode: constant.ResourceProjectTestUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目测试",
			ResourceCode: constant.ResourceProjectTestDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目测试用例
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目测试用例",
			ResourceCode: constant.ResourceProjectTestCase,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目测试用例列表",
			ResourceCode: constant.ResourceProjectTestCaseList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目测试用例详情",
			ResourceCode: constant.ResourceProjectTestCaseInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目测试用例",
			ResourceCode: constant.ResourceProjectTestCaseAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目测试用例",
			ResourceCode: constant.ResourceProjectTestCaseUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目测试用例",
			ResourceCode: constant.ResourceProjectTestCaseDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目测试用例步骤
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目测试用例步骤",
			ResourceCode: constant.ResourceProjectTestCaseStep,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目测试用例步骤列表",
			ResourceCode: constant.ResourceProjectTestCaseStepList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目测试用例步骤详情",
			ResourceCode: constant.ResourceProjectTestCaseStepInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目测试用例步骤",
			ResourceCode: constant.ResourceProjectTestCaseStepAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目测试用例步骤",
			ResourceCode: constant.ResourceProjectTestCaseStepUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目测试用例步骤",
			ResourceCode: constant.ResourceProjectTestCaseStepDelete,
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
