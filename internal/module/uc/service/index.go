package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"
)

func InitService() {
	// 初始化一些数据
	_ = database.AutoMigrate(model.Models...)
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
	resources := []*model.Resource{
		// 用户
		{
			ResourceName: "用户",
			ResourceCode: constant.ResourceUser,
			HttpMethod:   "ALL",
		},
		{
			ResourceName: "添加用户",
			ResourceCode: constant.ResourceUserAdd,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "删除用户",
			ResourceCode: constant.ResourceUserDelete,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "更新用户",
			ResourceCode: constant.ResourceUserUpdateUser,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "更新自己",
			ResourceCode: constant.ResourceUserUpdate,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "用户列表",
			ResourceCode: constant.ResourceUserList,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "用户详情",
			ResourceCode: constant.ResourceUserInstance,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "用户分配角色",
			ResourceCode: constant.ResourceUserDispatchRoles,
			HttpMethod:   "POST",
		},
		// 角色
		{
			ResourceName: "角色",
			ResourceCode: constant.ResourceRole,
			HttpMethod:   "ALL",
		},
		{
			ResourceName: "添加角色",
			ResourceCode: constant.ResourceRoleAdd,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "删除角色",
			ResourceCode: constant.ResourceRoleDelete,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "更新角色",
			ResourceCode: constant.ResourceRoleUpdate,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "角色列表",
			ResourceCode: constant.ResourceRoleList,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "角色详情",
			ResourceCode: constant.ResourceRoleInstance,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "角色分配资源",
			ResourceCode: constant.ResourceRoleDispatchResources,
			HttpMethod:   "POST",
		},
		// 资源
		{
			ResourceName: "资源",
			ResourceCode: constant.ResourceResource,
			HttpMethod:   "ALL",
		},
		{
			ResourceName: "添加资源",
			ResourceCode: constant.ResourceResourceAdd,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "删除资源",
			ResourceCode: constant.ResourceResourceDelete,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "资源列表",
			ResourceCode: constant.ResourceResourceList,
			HttpMethod:   "GET",
		},
		// 部门
		{
			ResourceName: "部门",
			ResourceCode: constant.ResourceDepartment,
			HttpMethod:   "ALL",
		},
		{
			ResourceName: "添加部门",
			ResourceCode: constant.ResourceDepartmentAdd,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "删除部门",
			ResourceCode: constant.ResourceDepartmentDelete,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "更新部门",
			ResourceCode: constant.ResourceDepartmentUpdate,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "部门列表",
			ResourceCode: constant.ResourceDepartmentList,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "部门详情",
			ResourceCode: constant.ResourceDepartmentInstance,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "部门添加用户",
			ResourceCode: constant.ResourceDepartmentAddUser,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "部门删除用户",
			ResourceCode: constant.ResourceDepartmentRemoveUser,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "更新部门名称",
			ResourceCode: constant.ResourceDepartmentUpdateName,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "变更父级部门",
			ResourceCode: constant.ResourceDepartmentChangeParent,
			HttpMethod:   "PUT|POST",
		},
		// 三方登录源
		{
			ResourceName: "三方登录源",
			ResourceCode: constant.ResourceThirdSource,
			HttpMethod:   "ALL",
		},
		{
			ResourceName: "添加三方登录源",
			ResourceCode: constant.ResourceThirdSourceAdd,
			HttpMethod:   "POST",
		},
		{
			ResourceName: "删除三方登录源",
			ResourceCode: constant.ResourceThirdSourceDelete,
			HttpMethod:   "DELETE|POST",
		},
		{
			ResourceName: "更新三方登录源",
			ResourceCode: constant.ResourceThirdSourceUpdate,
			HttpMethod:   "PUT|POST",
		},
		{
			ResourceName: "三方登录源列表",
			ResourceCode: constant.ResourceThirdSourceList,
			HttpMethod:   "GET",
		},
		{
			ResourceName: "三方登录源详情",
			ResourceCode: constant.ResourceThirdSourceInstance,
			HttpMethod:   "GET",
		},
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
