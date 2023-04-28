package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/constant"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/shared"
	"strconv"
	"strings"
)

func NeedAuthorization(code string) fiber.Handler {
	code = strings.ToLower(code)
	if code == "" || code == "anon" {
		return func(ctx *fiber.Ctx) error {
			return ctx.Next()
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(shared.JwtSecret),
		ContextKey:    "jwt-subject",
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization,cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).SendString("无效的token信息")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Authorization Token")
			}
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			jwtToken := c.Locals("jwt-subject").(*jwt.Token)
			claims := jwtToken.Claims.(jwt.MapClaims)
			uid := claims[shared.JwtClaimUserId].(string)
			sid := claims[shared.JwtClaimSessionId].(string)
			tenantId, hasTenantId := claims[shared.JwtClaimTenantId].(string)

			conn := cache.Get()
			defer func(conn redis.Conn) {
				_ = conn.Close()
			}(conn)
			cachedUid, err := redis.String(conn.Do("GET", shared.RedisSessionIdKey+sid))
			if err != nil {
				if err != redis.ErrNil {
					logrus.Errorln(err)
				}
				return c.Status(fiber.StatusUnauthorized).SendString("token信息已失效")
			}
			if cachedUid != uid {
				return c.Status(fiber.StatusUnauthorized).SendString("token信息不正确")
			}

			// 判断是否有权限 1、读取用户的权限信息 2、判断是否有权限
			// 获取用户角色
			roleIds, _ := redis.Uint64s(conn.Do("SMEMBERS", shared.RedisKeyUserRoles+uid))
			if len(roleIds) == 0 {
				// 获取该用户的角色id存入缓存
				userId, _ := strconv.ParseUint(uid, 10, 64)
				if userId == 0 {
					return c.Status(fiber.StatusUnauthorized).SendString("token信息已失效")
				}
				roleIds, err = service.UserService.RoleIds(userId)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString("系统错误")
				}
				for _, roleId := range roleIds {
					_, _ = conn.Do("SADD", shared.RedisKeyUserRoles+uid, roleId)
				}
			}
			_, _ = conn.Do("EXPIRE", shared.RedisKeyUserRoles+uid, 3*24*60*60)

			hasAuth := false
			for _, roleId := range roleIds {
				if roleId == constant.SuperAdminRoleId {
					hasAuth = true
					break
				} else {
					roleIdStr := strconv.FormatUint(roleId, 10)
					codes, _ := redis.Strings(conn.Do("GET", shared.RedisKeyRoleResourceCode+roleIdStr))
					if len(codes) == 0 {
						// 缓存没有，那么就去数据库取出来放进去
						codes, err = service.RoleService.ResourceCodes(roleId)
						for _, resourceCode := range codes {
							rc := resourceCode
							_, _ = conn.Do("SADD", shared.RedisKeyRoleResourceCode+roleIdStr, rc)
							if code == rc {
								hasAuth = true
							} else if strings.HasPrefix(code, rc+":") {
								hasAuth = true
							}
						}
					}
					_, _ = conn.Do("EXPIRE", shared.RedisKeyRoleResourceCode+roleIdStr, 3*24*60*60)
					if hasAuth {
						break
					}
					for _, resourceCode := range codes {
						rc := resourceCode
						if code == rc {
							hasAuth = true
							break
						} else if strings.HasPrefix(code, rc+":") {
							hasAuth = true
							break
						}
					}
				}
				if hasAuth {
					break
				}
			}

			if !hasAuth {
				return c.Status(fiber.StatusUnauthorized).SendString("无权限")
			}

			// 有权限，那么就把用户信息放到上下文中
			c.Locals(shared.JwtClaimUserId, uid)
			// 如果有租户，则租户信息也放入
			if hasTenantId {
				c.Locals(shared.JwtClaimTenantId, tenantId)
			}
			// token续期
			_, _ = conn.Do("EXPIRE", shared.RedisSessionIdKey+sid, config.GetInt("userTokenExpire"))
			return c.Next()
		},
	})
}
