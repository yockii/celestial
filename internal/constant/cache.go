package constant

import "github.com/yockii/ruomu-core/config"

var (
	RedisKeyDingtalkAccessToken = "dingtalk:access_token"

	RedisSessionIdKey = "sessionId"

	RedisKeyUserRoles        = "userRole"
	RedisKeyRoleResourceCode = "roleResourceCode"
	RedisKeyRoleDataPerm     = "roleDataPerm"

	RedisKeyUserRolesInProject = "userRolesInProject"
)

const (
	JwtSecret            = "yyyooccckkiiiiiiii"
	JwtClaimUserId       = "uid"
	JwtClaimTenantId     = "tid"
	JwtClaimSessionId    = "sid"
	JwtClaimUserDataPerm = "dataPerm"

	SuperAdmin = "superAdmin"
)

func init() {
	appName := config.GetString("redis.app") + ":"

	RedisKeyDingtalkAccessToken = appName + RedisKeyDingtalkAccessToken
	RedisSessionIdKey = appName + RedisSessionIdKey
	RedisKeyUserRoles = appName + RedisKeyUserRoles
	RedisKeyRoleResourceCode = appName + RedisKeyRoleResourceCode
	RedisKeyRoleDataPerm = appName + RedisKeyRoleDataPerm
	RedisKeyUserRolesInProject = appName + RedisKeyUserRolesInProject
}
