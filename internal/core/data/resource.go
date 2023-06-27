package data

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"
)

func InitData() {
	// 初始化一些数据
	_ = database.AutoMigrate(constant.Models...)
	// 初始化一个admin用户
	adminUser := &model.User{
		Username: "admin",
	}
	{
		pwd, _ := bcrypt.GenerateFromPassword([]byte(constant.AdminDefaultPassword), bcrypt.DefaultCost)
		if err := database.DB.Where(adminUser).Attrs(&model.User{
			ID:       util.SnowflakeId(),
			RealName: "管理员",
			Status:   model.UserStatusNormal,
			Password: string(pwd),
		}).FirstOrCreate(adminUser).Error; err != nil {
			logger.Errorln(err)
		}
	}

	// 初始化一个超级管理员角色
	superAdminRole := &model.Role{
		ID:             constant.SuperAdminRoleId,
		Type:           model.RoleTypeSuperAdmin,
		DataPermission: model.RoleDataPermissionAll,
		Status:         model.RoleStatusNormal,
	}
	{
		if err := database.DB.Where(superAdminRole).Attrs(&model.Role{
			Name: "超级管理员",
		}).FirstOrCreate(superAdminRole).Error; err != nil {
			logger.Errorln(err)
		}
	}

	// 关联admin和超级管理员角色
	{
		userRole := &model.UserRole{
			UserID: adminUser.ID,
			RoleID: superAdminRole.ID,
		}
		if err := database.DB.Where(userRole).Attrs(&model.UserRole{
			ID: util.SnowflakeId(),
		}).FirstOrCreate(userRole).Error; err != nil {
			logger.Errorln(err)
		}
	}

	// 初始化用户中心的资源
	var resources []*model.Resource

	// 首页
	{
		resources = append(resources,
			&model.Resource{
				ResourceName: "首页",
				ResourceCode: constant.ResourceHome,
				Type:         1,
			})

		// 仪表盘
		{
			resources = append(resources, &model.Resource{
				ResourceName: "仪表盘",
				ResourceCode: constant.ResourceDashboard,
				Type:         1,
			})
		}
	}
	// 项目
	{
		resources = append(resources,
			&model.Resource{
				ResourceName: "项目",
				ResourceCode: constant.ResourceProject,
				Type:         1,
			})
		resources = append(resources, &model.Resource{
			ResourceName: "项目列表",
			ResourceCode: constant.ResourceProjectList,
			Type:         1,
		})
		resources = append(resources, &model.Resource{
			ResourceName: "添加项目",
			ResourceCode: constant.ResourceProjectAdd,
			Type:         1,
		})
		// 项目详情
		{
			resources = append(resources, &model.Resource{
				ResourceName: "项目详情",
				ResourceCode: constant.ResourceProjectInstance,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新项目",
				ResourceCode: constant.ResourceProjectUpdate,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除项目",
				ResourceCode: constant.ResourceProjectDelete,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "项目成员",
				ResourceCode: constant.ResourceProjectMember,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加项目成员",
				ResourceCode: constant.ResourceProjectMemberAdd,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除项目成员",
				ResourceCode: constant.ResourceProjectMemberDelete,
				Type:         2,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新项目成员",
				ResourceCode: constant.ResourceProjectMemberUpdate,
				Type:         2,
			})
			// 项目计划
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目计划",
					ResourceCode: constant.ResourceProjectPlan,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目计划",
					ResourceCode: constant.ResourceProjectPlanAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目计划",
					ResourceCode: constant.ResourceProjectPlanUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目计划",
					ResourceCode: constant.ResourceProjectPlanDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目计划详情",
					ResourceCode: constant.ResourceProjectPlanInstance,
					Type:         2,
				})
			}
			// 功能模块
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目模块",
					ResourceCode: constant.ResourceProjectModule,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目模块",
					ResourceCode: constant.ResourceProjectModuleAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目模块",
					ResourceCode: constant.ResourceProjectModuleUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目模块",
					ResourceCode: constant.ResourceProjectModuleDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目模块评审",
					ResourceCode: constant.ResourceProjectModuleReview,
					Type:         2,
				})
			}
			// 项目需求
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目需求",
					ResourceCode: constant.ResourceProjectRequirement,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目需求",
					ResourceCode: constant.ResourceProjectRequirementAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目需求",
					ResourceCode: constant.ResourceProjectRequirementUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目需求",
					ResourceCode: constant.ResourceProjectRequirementDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目需求详情",
					ResourceCode: constant.ResourceProjectRequirementInstance,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目需求设计",
					ResourceCode: constant.ResourceProjectRequirementStatusDesign,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目需求评审",
					ResourceCode: constant.ResourceProjectRequirementStatusReview,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目需求状态置为已完成",
					ResourceCode: constant.ResourceProjectRequirementStatusCompleted,
					Type:         2,
				})
			}
			// 项目任务
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目任务",
					ResourceCode: constant.ResourceProjectTask,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目任务",
					ResourceCode: constant.ResourceProjectTaskAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目任务",
					ResourceCode: constant.ResourceProjectTaskUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目任务",
					ResourceCode: constant.ResourceProjectTaskDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目任务列表",
					ResourceCode: constant.ResourceProjectTaskList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目任务详情",
					ResourceCode: constant.ResourceProjectTaskInstance,
					Type:         2,
				})
				// 5个更新状态
				resources = append(resources, &model.Resource{
					ResourceName: "取消项目任务",
					ResourceCode: constant.ResourceProjectTaskCancel,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "任务成员确认项目任务",
					ResourceCode: constant.ResourceProjectTaskConfirm,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "任务成员开始项目任务",
					ResourceCode: constant.ResourceProjectTaskStart,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "任务成员完成项目任务",
					ResourceCode: constant.ResourceProjectTaskDone,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "重启项目任务",
					ResourceCode: constant.ResourceProjectTaskRestart,
					Type:         2,
				})
			}
			// 项目测试
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目测试",
					ResourceCode: constant.ResourceProjectTest,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "增加项目测试轮",
					ResourceCode: constant.ResourceProjectTestAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目测试轮",
					ResourceCode: constant.ResourceProjectTestUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目测试轮",
					ResourceCode: constant.ResourceProjectTestDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目测试轮封版",
					ResourceCode: constant.ResourceProjectTestClose,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目测试列表",
					ResourceCode: constant.ResourceProjectTestList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目测试详情",
					ResourceCode: constant.ResourceProjectTestInstance,
					Type:         2,
				})

				// 测试用例
				{
					resources = append(resources, &model.Resource{
						ResourceName: "项目测试用例",
						ResourceCode: constant.ResourceProjectTestCase,
						Type:         2,
					})
					resources = append(resources, &model.Resource{
						ResourceName: "添加项目测试用例",
						ResourceCode: constant.ResourceProjectTestCaseAdd,
						Type:         2,
					})
					resources = append(resources, &model.Resource{
						ResourceName: "更新项目测试用例",
						ResourceCode: constant.ResourceProjectTestCaseUpdate,
						Type:         2,
					})
					resources = append(resources, &model.Resource{
						ResourceName: "删除项目测试用例",
						ResourceCode: constant.ResourceProjectTestCaseDelete,
						Type:         2,
					})
					resources = append(resources, &model.Resource{
						ResourceName: "项目测试用例列表",
						ResourceCode: constant.ResourceProjectTestCaseList,
						Type:         2,
					})
					resources = append(resources, &model.Resource{
						ResourceName: "项目测试用例详情",
						ResourceCode: constant.ResourceProjectTestCaseInstance,
						Type:         2,
					})
					// 测试用例项
					{
						resources = append(resources, &model.Resource{
							ResourceName: "项目测试用例项",
							ResourceCode: constant.ResourceProjectTestCaseItem,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "添加项目测试用例项",
							ResourceCode: constant.ResourceProjectTestCaseItemAdd,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "更新项目测试用例项",
							ResourceCode: constant.ResourceProjectTestCaseItemUpdate,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "更新项目测试用例项状态",
							ResourceCode: constant.ResourceProjectTestCaseItemUpdateStatus,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "删除项目测试用例项",
							ResourceCode: constant.ResourceProjectTestCaseItemDelete,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "项目测试用例项列表",
							ResourceCode: constant.ResourceProjectTestCaseItemList,
							Type:         2,
						})
						resources = append(resources, &model.Resource{
							ResourceName: "项目测试用例项详情",
							ResourceCode: constant.ResourceProjectTestCaseItemInstance,
							Type:         2,
						})

						// 测试用例项步骤
						{
							resources = append(resources, &model.Resource{
								ResourceName: "项目测试用例项步骤",
								ResourceCode: constant.ResourceProjectTestCaseItemStep,
								Type:         2,
							})
							resources = append(resources, &model.Resource{
								ResourceName: "添加项目测试用例项步骤",
								ResourceCode: constant.ResourceProjectTestCaseItemStepAdd,
								Type:         2,
							})
							resources = append(resources, &model.Resource{
								ResourceName: "更新项目测试用例项步骤",
								ResourceCode: constant.ResourceProjectTestCaseItemStepUpdate,
								Type:         2,
							})
							resources = append(resources, &model.Resource{
								ResourceName: "删除项目测试用例项步骤",
								ResourceCode: constant.ResourceProjectTestCaseItemStepDelete,
								Type:         2,
							})
							resources = append(resources, &model.Resource{
								ResourceName: "项目测试用例项步骤列表",
								ResourceCode: constant.ResourceProjectTestCaseItemStepList,
								Type:         2,
							})
							resources = append(resources, &model.Resource{
								ResourceName: "项目测试用例项步骤详情",
								ResourceCode: constant.ResourceProjectTestCaseItemStepInstance,
								Type:         2,
							})
						}
					}
				}
			}

			// 项目缺陷
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目缺陷",
					ResourceCode: constant.ResourceProjectIssue,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目缺陷",
					ResourceCode: constant.ResourceProjectIssueAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目缺陷",
					ResourceCode: constant.ResourceProjectIssueUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目缺陷",
					ResourceCode: constant.ResourceProjectIssueDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目缺陷列表",
					ResourceCode: constant.ResourceProjectIssueList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目缺陷详情",
					ResourceCode: constant.ResourceProjectIssueInstance,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目缺陷指派处理人",
					ResourceCode: constant.ResourceProjectIssueAssign,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "开始处理项目缺陷",
					ResourceCode: constant.ResourceProjectIssueStart,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "完成项目缺陷待验证",
					ResourceCode: constant.ResourceProjectIssueDone,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "验证项目缺陷",
					ResourceCode: constant.ResourceProjectIssueVerify,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "关闭项目缺陷",
					ResourceCode: constant.ResourceProjectIssueClose,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "拒绝项目缺陷",
					ResourceCode: constant.ResourceProjectIssueReject,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "重新打开项目缺陷",
					ResourceCode: constant.ResourceProjectIssueReopen,
					Type:         2,
				})
			}

			// 项目变更
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目变更",
					ResourceCode: constant.ResourceProjectChange,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目变更",
					ResourceCode: constant.ResourceProjectChangeAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目变更",
					ResourceCode: constant.ResourceProjectChangeUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目变更",
					ResourceCode: constant.ResourceProjectChangeDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目变更列表",
					ResourceCode: constant.ResourceProjectChangeList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目变更详情",
					ResourceCode: constant.ResourceProjectChangeInstance,
					Type:         2,
				})
			}

			// 项目风险
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目风险",
					ResourceCode: constant.ResourceProjectRisk,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目风险",
					ResourceCode: constant.ResourceProjectRiskAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目风险",
					ResourceCode: constant.ResourceProjectRiskUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目风险",
					ResourceCode: constant.ResourceProjectRiskDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目风险列表",
					ResourceCode: constant.ResourceProjectRiskList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目风险详情",
					ResourceCode: constant.ResourceProjectRiskInstance,
					Type:         2,
				})
			}

			// 项目资产
			{
				resources = append(resources, &model.Resource{
					ResourceName: "项目资产",
					ResourceCode: constant.ResourceProjectAsset,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "添加项目资产",
					ResourceCode: constant.ResourceProjectAssetAdd,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "更新项目资产",
					ResourceCode: constant.ResourceProjectAssetUpdate,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "删除项目资产",
					ResourceCode: constant.ResourceProjectAssetDelete,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目资产列表",
					ResourceCode: constant.ResourceProjectAssetList,
					Type:         2,
				})
				resources = append(resources, &model.Resource{
					ResourceName: "项目资产详情",
					ResourceCode: constant.ResourceProjectAssetInstance,
					Type:         2,
				})
			}
		}
	}

	// 任务
	{
		resources = append(resources, &model.Resource{
			ResourceName: "任务",
			ResourceCode: constant.ResourceTask,
			Type:         1,
		})
	}

	// 测试
	{
		resources = append(resources, &model.Resource{
			ResourceName: "测试",
			ResourceCode: constant.ResourceTest,
			Type:         1,
		})
	}

	// 资产
	{
		resources = append(resources, &model.Resource{
			ResourceName: "资产",
			ResourceCode: constant.ResourceAsset,
			Type:         1,
		})
		// 文件
		{
			resources = append(resources, &model.Resource{
				ResourceName: "文件",
				ResourceCode: constant.ResourceFile,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加文件",
				ResourceCode: constant.ResourceFileAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除文件",
				ResourceCode: constant.ResourceFileDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新文件",
				ResourceCode: constant.ResourceFileUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "文件详情",
				ResourceCode: constant.ResourceFileInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "文件列表",
				ResourceCode: constant.ResourceFileList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "文件下载",
				ResourceCode: constant.ResourceFileDownload,
				Type:         1,
			})
		}
		// 通用测试用例
		{
			resources = append(resources, &model.Resource{
				ResourceName: "通用测试用例",
				ResourceCode: constant.ResourceCommonTestCase,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加通用测试用例",
				ResourceCode: constant.ResourceCommonTestCaseAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除通用测试用例",
				ResourceCode: constant.ResourceCommonTestCaseDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新通用测试用例",
				ResourceCode: constant.ResourceCommonTestCaseUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "通用测试用例列表",
				ResourceCode: constant.ResourceCommonTestCaseList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "通用测试用例详情",
				ResourceCode: constant.ResourceCommonTestCaseInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加通用测试用例项",
				ResourceCode: constant.ResourceCommonTestCaseAddItem,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除通用测试用例项",
				ResourceCode: constant.ResourceCommonTestCaseDeleteItem,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新通用测试用例项",
				ResourceCode: constant.ResourceCommonTestCaseUpdateItem,
				Type:         1,
			})
		}
	}

	// 系统
	{
		resources = append(resources, &model.Resource{
			ResourceName: "系统",
			ResourceCode: constant.ResourceSystem,
			Type:         1,
		})
		// 资源列表
		{
			resources = append(resources, &model.Resource{
				ResourceName: "资源列表",
				ResourceCode: constant.ResourceResourceList,
				Type:         1,
			})
		}
		// 用户
		{
			resources = append(resources, &model.Resource{
				ResourceName: "用户",
				ResourceCode: constant.ResourceUser,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加用户",
				ResourceCode: constant.ResourceUserAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除用户",
				ResourceCode: constant.ResourceUserDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新用户",
				ResourceCode: constant.ResourceUserUpdateUser,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新自己",
				ResourceCode: constant.ResourceUserUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "用户列表",
				ResourceCode: constant.ResourceUserList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "用户详情",
				ResourceCode: constant.ResourceUserInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "用户分配角色",
				ResourceCode: constant.ResourceUserDispatchRoles,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "获取用户角色",
				ResourceCode: constant.ResourceUserRoles,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "重置用户密码",
				ResourceCode: constant.ResourceUserResetUserPassword,
				Type:         1,
			})
		}
		// 角色
		{
			resources = append(resources, &model.Resource{
				ResourceName: "角色",
				ResourceCode: constant.ResourceRole,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加角色",
				ResourceCode: constant.ResourceRoleAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除角色",
				ResourceCode: constant.ResourceRoleDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新角色",
				ResourceCode: constant.ResourceRoleUpdate,
				Type:         1,
			})
			//resources = append(resources, &model.Resource{
			//	ResourceName: "角色列表",
			//	ResourceCode: constant.ResourceRoleList,
			//	Type:         3,
			//})
			resources = append(resources, &model.Resource{
				ResourceName: "角色详情",
				ResourceCode: constant.ResourceRoleInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "角色分配资源",
				ResourceCode: constant.ResourceRoleDispatchResources,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "获取角色资源",
				ResourceCode: constant.ResourceRoleResources,
				Type:         1,
			})
		}
		// 阶段
		{
			resources = append(resources, &model.Resource{
				ResourceName: "阶段",
				ResourceCode: constant.ResourceStage,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加阶段",
				ResourceCode: constant.ResourceStageAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除阶段",
				ResourceCode: constant.ResourceStageDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新阶段",
				ResourceCode: constant.ResourceStageUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "阶段列表",
				ResourceCode: constant.ResourceStageList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "阶段详情",
				ResourceCode: constant.ResourceStageInstance,
				Type:         1,
			})
		}
		// 部门
		{
			resources = append(resources, &model.Resource{
				ResourceName: "部门",
				ResourceCode: constant.ResourceDepartment,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加部门",
				ResourceCode: constant.ResourceDepartmentAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除部门",
				ResourceCode: constant.ResourceDepartmentDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新部门",
				ResourceCode: constant.ResourceDepartmentUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "部门列表",
				ResourceCode: constant.ResourceDepartmentList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "部门详情",
				ResourceCode: constant.ResourceDepartmentInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "部门添加用户",
				ResourceCode: constant.ResourceDepartmentAddUser,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "部门删除用户",
				ResourceCode: constant.ResourceDepartmentRemoveUser,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新部门名称",
				ResourceCode: constant.ResourceDepartmentUpdateName,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "变更父级部门",
				ResourceCode: constant.ResourceDepartmentChangeParent,
				Type:         1,
			})
		}
		// 三方登录源
		{
			resources = append(resources, &model.Resource{
				ResourceName: "三方登录源",
				ResourceCode: constant.ResourceThirdSource,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加三方登录源",
				ResourceCode: constant.ResourceThirdSourceAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除三方登录源",
				ResourceCode: constant.ResourceThirdSourceDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新三方登录源",
				ResourceCode: constant.ResourceThirdSourceUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "三方登录源列表",
				ResourceCode: constant.ResourceThirdSourceList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "三方登录源详情",
				ResourceCode: constant.ResourceThirdSourceInstance,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "同步三方登录源数据",
				ResourceCode: constant.ResourceThirdSourceSync,
				Type:         1,
			})
		}
		// 资产目录
		{
			resources = append(resources, &model.Resource{
				ResourceName: "资产目录",
				ResourceCode: constant.ResourceAssetCategory,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加资产目录",
				ResourceCode: constant.ResourceAssetCategoryAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除资产目录",
				ResourceCode: constant.ResourceAssetCategoryDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新资产目录",
				ResourceCode: constant.ResourceAssetCategoryUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "资产目录列表",
				ResourceCode: constant.ResourceAssetCategoryList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "资产目录详情",
				ResourceCode: constant.ResourceAssetCategoryInstance,
				Type:         1,
			})
		}
		// oss配置
		{
			resources = append(resources, &model.Resource{
				ResourceName: "oss配置",
				ResourceCode: constant.ResourceOssConfig,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "添加oss配置",
				ResourceCode: constant.ResourceOssConfigAdd,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "删除oss配置",
				ResourceCode: constant.ResourceOssConfigDelete,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "更新oss配置",
				ResourceCode: constant.ResourceOssConfigUpdate,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "oss配置列表",
				ResourceCode: constant.ResourceOssConfigList,
				Type:         1,
			})
			resources = append(resources, &model.Resource{
				ResourceName: "oss配置详情",
				ResourceCode: constant.ResourceOssConfigInstance,
				Type:         1,
			})
		}
	}

	for _, resource := range resources {
		//没有就添加资源
		if err := database.DB.Where(resource).Attrs(&model.Resource{
			ID: util.SnowflakeId(),
		}).FirstOrCreate(resource).Error; err != nil {
			logger.Errorln(err)
		}
	}
}
