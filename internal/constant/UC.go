package constant

import "github.com/yockii/ruomu-core/config"

var (
	// AdminDefaultPassword 初始admin账号的默认密码。生产环境务必通过配置 admin.initPassword
	// 覆盖，并在首次登录后修改，默认值仅为向后兼容保留。
	AdminDefaultPassword = "Admin123!@#"
)

const (
	SuperAdminRoleId = 1

	// MaskedSecret 敏感信息（如密钥）对外展示时使用的掩码
	MaskedSecret = "*****"

	// SysConfigKeyUsernamePasswordEnabled 系统配置表中的登录开关键
	SysConfigKeyUsernamePasswordEnabled = "login.usernamePasswordEnabled"
)

var AdminUserID uint64

func init() {
	if s := config.GetString("admin.initPassword"); s != "" {
		AdminDefaultPassword = s
	}
}
