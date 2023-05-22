package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/constant"
	"github.com/yockii/celestial/internal/module/project/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

func InitService() {
	_ = database.AutoMigrate(model.Models...)
	// 初始化项目的资源
	var resources []*ucModel.Resource
	// 阶段
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "阶段",
			ResourceCode: constant.ResourceStage,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "阶段列表",
			ResourceCode: constant.ResourceStageList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "阶段详情",
			ResourceCode: constant.ResourceStageInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加阶段",
			ResourceCode: constant.ResourceStageAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改阶段",
			ResourceCode: constant.ResourceStageUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除阶段",
			ResourceCode: constant.ResourceStageDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目",
			ResourceCode: constant.ResourceProject,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目列表",
			ResourceCode: constant.ResourceProjectList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目详情",
			ResourceCode: constant.ResourceProjectInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目",
			ResourceCode: constant.ResourceProjectAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目",
			ResourceCode: constant.ResourceProjectUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目",
			ResourceCode: constant.ResourceProjectDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目成员
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目成员",
			ResourceCode: constant.ResourceProjectMember,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目成员列表",
			ResourceCode: constant.ResourceProjectMemberList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目成员详情",
			ResourceCode: constant.ResourceProjectMemberInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目成员",
			ResourceCode: constant.ResourceProjectMemberAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目成员",
			ResourceCode: constant.ResourceProjectMemberUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目成员",
			ResourceCode: constant.ResourceProjectMemberDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目变更
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目变更",
			ResourceCode: constant.ResourceProjectChange,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目变更列表",
			ResourceCode: constant.ResourceProjectChangeList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目变更详情",
			ResourceCode: constant.ResourceProjectChangeInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目变更",
			ResourceCode: constant.ResourceProjectChangeAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目变更",
			ResourceCode: constant.ResourceProjectChangeUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目变更",
			ResourceCode: constant.ResourceProjectChangeDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目问题
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目问题",
			ResourceCode: constant.ResourceProjectIssue,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目问题列表",
			ResourceCode: constant.ResourceProjectIssueList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目问题详情",
			ResourceCode: constant.ResourceProjectIssueInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目问题",
			ResourceCode: constant.ResourceProjectIssueAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目问题",
			ResourceCode: constant.ResourceProjectIssueUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目问题",
			ResourceCode: constant.ResourceProjectIssueDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目风险
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目风险",
			ResourceCode: constant.ResourceProjectRisk,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目风险列表",
			ResourceCode: constant.ResourceProjectRiskList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目风险详情",
			ResourceCode: constant.ResourceProjectRiskInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目风险",
			ResourceCode: constant.ResourceProjectRiskAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目风险",
			ResourceCode: constant.ResourceProjectRiskUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目风险",
			ResourceCode: constant.ResourceProjectRiskDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目资产
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目资产",
			ResourceCode: constant.ResourceProjectAsset,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目资产列表",
			ResourceCode: constant.ResourceProjectAssetList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目资产详情",
			ResourceCode: constant.ResourceProjectAssetInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目资产",
			ResourceCode: constant.ResourceProjectAssetAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目资产",
			ResourceCode: constant.ResourceProjectAssetUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目资产",
			ResourceCode: constant.ResourceProjectAssetDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目模块
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目模块",
			ResourceCode: constant.ResourceProjectModule,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目模块列表",
			ResourceCode: constant.ResourceProjectModuleList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目模块详情",
			ResourceCode: constant.ResourceProjectModuleInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目模块",
			ResourceCode: constant.ResourceProjectModuleAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目模块",
			ResourceCode: constant.ResourceProjectModuleUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目模块",
			ResourceCode: constant.ResourceProjectModuleDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目需求
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目需求",
			ResourceCode: constant.ResourceProjectRequirement,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目需求列表",
			ResourceCode: constant.ResourceProjectRequirementList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目需求详情",
			ResourceCode: constant.ResourceProjectRequirementInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目需求",
			ResourceCode: constant.ResourceProjectRequirementAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目需求",
			ResourceCode: constant.ResourceProjectRequirementUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目需求",
			ResourceCode: constant.ResourceProjectRequirementDelete,
			HttpMethod:   "DELETE|POST",
		})
	}
	// 项目计划
	{
		resources = append(resources, &ucModel.Resource{
			ResourceName: "项目计划",
			ResourceCode: constant.ResourceProjectPlan,
			HttpMethod:   "ALL",
		}, &ucModel.Resource{
			ResourceName: "项目计划列表",
			ResourceCode: constant.ResourceProjectPlanList,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "项目计划详情",
			ResourceCode: constant.ResourceProjectPlanInstance,
			HttpMethod:   "GET",
		}, &ucModel.Resource{
			ResourceName: "添加项目计划",
			ResourceCode: constant.ResourceProjectPlanAdd,
			HttpMethod:   "POST",
		}, &ucModel.Resource{
			ResourceName: "修改项目计划",
			ResourceCode: constant.ResourceProjectPlanUpdate,
			HttpMethod:   "PUT|POST",
		}, &ucModel.Resource{
			ResourceName: "删除项目计划",
			ResourceCode: constant.ResourceProjectPlanDelete,
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
