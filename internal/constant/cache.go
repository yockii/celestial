package constant

import "github.com/yockii/ruomu-core/config"

var (
	RedisKeyDingtalkAccessToken = "dingtalk:access_token"

	RedisSessionIdKey = "sessionId"

	RedisKeyUserRoles        = "userRole"
	RedisKeyRoleResourceCode = "roleResourceCode"
	RedisKeyRoleDataPerm     = "roleDataPerm"

	RedisKeyUserRolesInProject = "userRolesInProject"

	RedisKeySysConfig = "sysConfig"
)

var (
	// JwtSecret JWT签名密钥。生产环境务必通过配置 jwt.secret 覆盖为随机长字符串，
	// 默认值仅为向后兼容保留。
	JwtSecret = "yyyooccckkiiiiiiii"
)

const (
	JwtClaimUserId       = "uid"
	JwtClaimTenantId     = "tid"
	JwtClaimSessionId    = "sid"
	JwtClaimUserDataPerm = "dataPerm"

	SuperAdmin = "superAdmin"
)

func init() {
	if s := config.GetString("jwt.secret"); s != "" {
		JwtSecret = s
	}

	appName := config.GetString("redis.app") + ":"

	RedisKeyDingtalkAccessToken = appName + RedisKeyDingtalkAccessToken
	RedisSessionIdKey = appName + RedisSessionIdKey
	RedisKeyUserRoles = appName + RedisKeyUserRoles
	RedisKeyRoleResourceCode = appName + RedisKeyRoleResourceCode
	RedisKeyRoleDataPerm = appName + RedisKeyRoleDataPerm
	RedisKeyUserRolesInProject = appName + RedisKeyUserRolesInProject
	RedisKeySysConfig = appName + RedisKeySysConfig
}
