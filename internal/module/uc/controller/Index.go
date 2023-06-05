package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/celestial/internal/module/uc/constant"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/server"
)

// single直接注册路由

func InitRouter() {
	service.InitService()

	// 注册
	server.Post("/api/v1/register", UserController.Register)
	// 登录
	{
		server.Post("/api/v1/login", UserController.Login)
		server.Post("/api/v1/loginByDingtalkCode", UserController.LoginByDingtalkCode)
		server.Post("/api/v1/loginInDingtalk", UserController.LoginInDingtalk)
	}
	// 用户
	{
		user := server.Group("/api/v1/user")
		user.Post("/add", middleware.NeedAuthorization(constant.ResourceUserAdd), UserController.Add)
		user.Delete("/delete", middleware.NeedAuthorization(constant.ResourceUserDelete), UserController.Delete)
		user.Put("/updateUser", middleware.NeedAuthorization(constant.ResourceUserUpdateUser), UserController.UpdateUser)
		user.Put("/update", middleware.NeedAuthorization(constant.ResourceUserUpdate), UserController.UpdateSelf)
		user.Get("/list", middleware.NeedAuthorization(constant.ResourceUserList), UserController.List)
		user.Get("/instance", middleware.NeedAuthorization(constant.ResourceUserInstance), UserController.Instance)
		user.Post("/dispatchRoles", middleware.NeedAuthorization(constant.ResourceUserDispatchRoles), UserController.DispatchRoles)

		// 对于禁用put和delete方法时的处理
		user.Post("/delete", middleware.NeedAuthorization(constant.ResourceUserDelete), UserController.Delete)
		user.Post("/updateUser", middleware.NeedAuthorization(constant.ResourceUserUpdateUser), UserController.UpdateUser)
		user.Post("/update", middleware.NeedAuthorization(constant.ResourceUserUpdate), UserController.UpdateSelf)
	}

	// 部门
	{
		department := server.Group("/api/v1/department")
		department.Post("/add", middleware.NeedAuthorization(constant.ResourceDepartmentAdd), DepartmentController.Add)
		department.Delete("/delete", middleware.NeedAuthorization(constant.ResourceDepartmentDelete), DepartmentController.Delete)
		department.Put("/update", middleware.NeedAuthorization(constant.ResourceDepartmentUpdate), DepartmentController.Update)
		department.Put("/updateName", middleware.NeedAuthorization(constant.ResourceDepartmentUpdateName), DepartmentController.UpdateName)
		department.Put("/changeParent", middleware.NeedAuthorization(constant.ResourceDepartmentChangeParent), DepartmentController.ChangeParent)
		department.Get("/list", middleware.NeedAuthorization(constant.ResourceDepartmentList), DepartmentController.List)
		department.Get("/instance", middleware.NeedAuthorization(constant.ResourceDepartmentInstance), DepartmentController.Instance)

		department.Post("/addUser", middleware.NeedAuthorization(constant.ResourceDepartmentAddUser), DepartmentController.AddUser)
		department.Post("/removeUser", middleware.NeedAuthorization(constant.ResourceDepartmentRemoveUser), DepartmentController.RemoveUser)

		// 对于禁用put和delete方法时的处理
		department.Post("/delete", middleware.NeedAuthorization(constant.ResourceDepartmentDelete), DepartmentController.Delete)
		department.Post("/update", middleware.NeedAuthorization(constant.ResourceDepartmentUpdate), DepartmentController.Update)
		department.Post("/updateName", middleware.NeedAuthorization(constant.ResourceDepartmentUpdateName), DepartmentController.UpdateName)
		department.Post("/changeParent", middleware.NeedAuthorization(constant.ResourceDepartmentChangeParent), DepartmentController.ChangeParent)
	}

	// 角色
	{
		role := server.Group("/api/v1/role")
		role.Post("/add", middleware.NeedAuthorization(constant.ResourceRoleAdd), RoleController.Add)
		role.Delete("/delete", middleware.NeedAuthorization(constant.ResourceRoleDelete), RoleController.Delete)
		role.Put("/update", middleware.NeedAuthorization(constant.ResourceRoleUpdate), RoleController.Update)
		role.Get("/list", middleware.NeedAuthorization(constant.ResourceRoleList), RoleController.List)
		role.Get("/instance", middleware.NeedAuthorization(constant.ResourceRoleInstance), RoleController.Instance)
		role.Post("/dispatchResources", middleware.NeedAuthorization(constant.ResourceRoleDispatchResources), RoleController.DispatchResources)
		role.Put("/setDefaultRole", middleware.NeedAuthorization(constant.ResourceRoleUpdate), RoleController.SetDefaultRole)

		// 对于禁用put和delete方法时的处理
		role.Post("/delete", middleware.NeedAuthorization(constant.ResourceRoleDelete), RoleController.Delete)
		role.Post("/update", middleware.NeedAuthorization(constant.ResourceRoleUpdate), RoleController.Update)
	}

	// 资源
	{
		resource := server.Group("/api/v1/resource")
		resource.Post("/add", middleware.NeedAuthorization(constant.ResourceResourceAdd), ResourceController.Add)
		resource.Delete("/delete", middleware.NeedAuthorization(constant.ResourceResourceDelete), ResourceController.Delete)
		//resource.Put("/update", middleware.NeedAuthorization("resource:update"), ResourceController.Update)
		resource.Get("/list", middleware.NeedAuthorization(constant.ResourceResourceList), ResourceController.List)

		// 对于禁用put和delete方法时的处理
		resource.Post("/delete", middleware.NeedAuthorization("resource:delete"), ResourceController.Delete)
	}

	// 三方登录源
	{
		thirdSource := server.Group("/api/v1/thirdSource")
		thirdSource.Post("/add", middleware.NeedAuthorization(constant.ResourceThirdSourceAdd), ThirdSourceController.Add)
		thirdSource.Delete("/delete", middleware.NeedAuthorization(constant.ResourceThirdSourceDelete), ThirdSourceController.Delete)
		thirdSource.Put("/update", middleware.NeedAuthorization(constant.ResourceThirdSourceUpdate), ThirdSourceController.Update)
		thirdSource.Get("/list", middleware.NeedAuthorization(constant.ResourceThirdSourceList), ThirdSourceController.List)
		thirdSource.Get("/instance", middleware.NeedAuthorization(constant.ResourceThirdSourceInstance), ThirdSourceController.Instance)
		thirdSource.Get("/publicList", ThirdSourceController.PublicList)
		// 对于禁用put和delete方法时的处理
		thirdSource.Post("/delete", middleware.NeedAuthorization(constant.ResourceThirdSourceDelete), ThirdSourceController.Delete)
		thirdSource.Post("/update", middleware.NeedAuthorization(constant.ResourceThirdSourceUpdate), ThirdSourceController.Update)
	}

	// 三方登录
	{
		// 钉钉回调
		server.Post("/api/v1/dingtalk/:id", ThirdLoginController.DingtalkCallback)
	}

}

// module方式用以下代码注入

//
//func Dispatch(code string, headers map[string]string, value []byte) ([]byte, error) {
//	switch code {
//	// 代码注入点
//	case shared.InjectCodeAuthorizationInfoByUserId:
//		return wrapCall(value, UserController.GetUserRoleIds)
//	case shared.InjectCodeAuthorizationInfoByRoleId:
//		return wrapCall(value, RoleController.GetRoleResourceCodes)
//	case shared.InjectCodeAddResourceInfo:
//		return wrapCall(value, ResourceController.Add)
//
//	// HTTP 注入点
//
//	// 用户
//	case constant2.InjectCodeUserLogin:
//		return wrapCall(value, UserController.Login)
//	case constant2.InjectCodeUserAdd:
//		return wrapCall(value, UserController.Add)
//	case constant2.InjectCodeUserUpdate:
//		return wrapCall(value, UserController.Update)
//	case constant2.InjectCodeUserDelete:
//		return wrapCall(value, UserController.Delete)
//	case constant2.InjectCodeUserInstance:
//		return wrapCall(value, UserController.Instance)
//	case constant2.InjectCodeUserList:
//		return wrapCall(value, UserController.List)
//	case constant2.InjectCodeUserRoles:
//		return wrapCall(value, UserController.GetUserRoleIds)
//
//	// 角色
//	case constant2.InjectCodeRoleAdd:
//		return wrapCall(value, RoleController.Add)
//	case constant2.InjectCodeRoleUpdate:
//		return wrapCall(value, RoleController.Update)
//	case constant2.InjectCodeRoleDelete:
//		return wrapCall(value, RoleController.Delete)
//	case constant2.InjectCodeRoleInstance:
//		return wrapCall(value, RoleController.Instance)
//	case constant2.InjectCodeRoleList:
//		return wrapCall(value, RoleController.List)
//	case constant2.InjectCodeRoleResourceCodes:
//		return wrapCall(value, RoleController.GetRoleResourceCodes)
//
//	// 资源
//	case constant2.InjectCodeResourceAdd:
//		return wrapCall(value, ResourceController.Add)
//	case constant2.InjectCodeResourceDelete:
//		return wrapCall(value, ResourceController.Delete)
//	case constant2.InjectCodeResourceList:
//		return wrapCall(value, ResourceController.List)
//	}
//	return nil, nil
//}
//
//func wrapCall(v []byte, f func([]byte) (any, error)) ([]byte, error) {
//	r, err := f(v)
//	if err != nil {
//		return nil, err
//	}
//	bs, err := json.Marshal(r)
//	return bs, err
//}
