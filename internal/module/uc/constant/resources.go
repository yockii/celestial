package constant

// 用户资源
const (
	ResourceUser              = "user"
	ResourceUserAdd           = "user:add"
	ResourceUserDelete        = "user:delete"
	ResourceUserUpdateUser    = "user:updateUser"
	ResourceUserUpdate        = "user:update"
	ResourceUserList          = "user:list"
	ResourceUserInstance      = "user:instance"
	ResourceUserDispatchRoles = "user:dispatchRoles"
)

// 部门资源
const (
	ResourceDepartment             = "department"
	ResourceDepartmentAdd          = "department:add"
	ResourceDepartmentDelete       = "department:delete"
	ResourceDepartmentUpdate       = "department:update"
	ResourceDepartmentUpdateName   = "department:updateName"
	ResourceDepartmentChangeParent = "department:changeParent"
	ResourceDepartmentList         = "department:list"
	ResourceDepartmentInstance     = "department:instance"
	ResourceDepartmentAddUser      = "department:addUser"
	ResourceDepartmentRemoveUser   = "department:removeUser"
)

// 角色资源
const (
	ResourceRole                  = "role"
	ResourceRoleAdd               = "role:add"
	ResourceRoleDelete            = "role:delete"
	ResourceRoleUpdate            = "role:update"
	ResourceRoleList              = "role:list"
	ResourceRoleInstance          = "role:instance"
	ResourceRoleDispatchResources = "role:dispatchResources"
)

// 资源资源
const (
	ResourceResource       = "resource"
	ResourceResourceAdd    = "resource:add"
	ResourceResourceDelete = "resource:delete"
	ResourceResourceList   = "resource:list"
	//ResourceResourceUpdate   = "resource:update"
	//ResourceResourceInstance = "resource:instance"
)

// 三方登录源资源
const (
	ResourceThirdSource         = "thirdSource"
	ResourceThirdSourceAdd      = "thirdSource:add"
	ResourceThirdSourceDelete   = "thirdSource:delete"
	ResourceThirdSourceUpdate   = "thirdSource:update"
	ResourceThirdSourceList     = "thirdSource:list"
	ResourceThirdSourceInstance = "thirdSource:instance"
)
