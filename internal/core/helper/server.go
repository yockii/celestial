package helper

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/project/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"strconv"
	"strings"
)

func GetCurrentUserID(ctx *fiber.Ctx) (uint64, error) {
	// 获取当前登录的用户
	uidStr, _ := ctx.Locals(constant.JwtClaimUserId).(string)
	if uidStr == "" {
		return 0, nil
	}
	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return 0, err
	}
	return uid, nil
}

func GetCurrentUserDataPermit(ctx *fiber.Ctx) (int, error) {
	// 获取当前登录用户的数据权限
	permit, _ := ctx.Locals(constant.JwtClaimUserDataPerm).(int)
	return permit, nil
}

func GetCurrentUserDataPermitInProject(ctx *fiber.Ctx, projectID uint64) (int, error) {
	uid, err := GetCurrentUserID(ctx)
	if err != nil {
		return 0, err
	}
	// 看看缓存中是否有当前用户在当前项目的数据权限
	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	permit := 0

	// 查询uid对应在项目中的角色id
	roleIdsInProject, err := GetUserRoleIdsInProject(uid, projectID)
	if err != nil {
		return 0, err
	}
	// 如果没有角色，那么就返回0
	if len(roleIdsInProject) == 0 {
		return 0, nil
	}

	for _, roleId := range roleIdsInProject {
		// 先从缓存中查询角色对应的数据权限
		thisRolePermit := 0
		thisRolePermit, err = redis.Int(conn.Do("HGET", constant.RedisKeyRoleDataPerm, roleId))
		if err != nil && !errors.Is(err, redis.ErrNil) {
			logger.Errorln(err)
			return 0, err
		}
		if thisRolePermit != 0 {
			if thisRolePermit < permit {
				permit = thisRolePermit
				continue
			}
		}
		// 如果不存在或者为0，那么就去数据库中查找
		role := &ucModel.Role{ID: roleId}
		if err = database.DB.Where(role).First(role).Error; err != nil {
			logger.Errorln(err)
			return 0, err
		}
		if permit == 0 || role.DataPermission < permit {
			permit = role.DataPermission
		}
		// 存入缓存中
		_, _ = conn.Do("HSET", constant.RedisKeyRoleDataPerm, roleId, role.DataPermission)
		// 设置失效时间
		_, _ = conn.Do("EXPIRE", constant.RedisKeyRoleDataPerm, 3*24*60*60)
	}
	return permit, nil
}

func GetUserRoleIdsInProject(uid, projectID uint64) (roleIDs []uint64, err error) {
	if projectID == 0 {
		return
	}

	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	roleIdKey := fmt.Sprintf(
		"%s:%d:%d",
		constant.RedisKeyUserRolesInProject,
		projectID,
		uid,
	)

	// 查询uid对应在项目中的角色id列表
	// 先从缓存中查找
	roleIDs, err = redis.Uint64s(conn.Do("SMEMBERS", roleIdKey))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		logger.Errorln(err)
		return
	}
	// 如果不存在或者为空，那么就去数据库中查找
	if len(roleIDs) == 0 {
		// 查询uid对应在项目中的角色id列表
		err = database.DB.Model(&model.ProjectMember{}).Where(&model.ProjectMember{
			ProjectID: projectID,
			UserID:    uid,
		}).Distinct().Pluck("role_id", &roleIDs).Error
		if err != nil {
			logger.Errorln(err)
			return
		}
		// 存入缓存中
		for _, roleID := range roleIDs {
			_, _ = conn.Do("SADD", roleIdKey, roleID)
		}
		// 设置失效时间
		_, _ = conn.Do("EXPIRE", roleIdKey, 3*24*60*60)
	}
	return
}

func HasResourceCodeInProject(uid, projectID uint64, codes ...string) (bool, error) {
	roleIDs, err := GetUserRoleIdsInProject(uid, projectID)
	if err != nil {
		return false, err
	}

	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	codeMap := make(map[string]struct{})
	for _, code := range codes {
		codeMap[code] = struct{}{}
	}

	hasAuth := false
	// 查询角色对应的资源权限
	for _, roleID := range roleIDs {
		if roleID == 0 {
			// 创建人
			hasAuth = true
			break
		}
		// 先从缓存中查找
		roleResourceKey := fmt.Sprintf("%s:%d", constant.RedisKeyRoleResourceCode, roleID)
		var resourceCodesInRole []string
		resourceCodesInRole, err = redis.Strings(conn.Do("SMEMBERS", roleResourceKey))
		if err != nil && !errors.Is(err, redis.ErrNil) {
			logger.Errorln(err)
			return false, err
		}
		dBCodeOp := false
		// 如果不存在或者为空，那么就去数据库中查找
		if len(resourceCodesInRole) == 0 {
			resourceCodesInRole, err = service.RoleService.ResourceCodes(roleID)
			if err != nil {
				logger.Errorln(err)
				return false, err
			}
			dBCodeOp = true
			// 存入缓存中
			args := []interface{}{roleResourceKey}
			for _, resourceCode := range resourceCodesInRole {
				args = append(args, resourceCode)
				if resourceCode == constant.ResourceAllProjectDetail {
					hasAuth = true
				} else if _, ok := codeMap[resourceCode]; ok {
					hasAuth = true
				} else {
					for _, code := range codes {
						if strings.HasPrefix(code, resourceCode+":") {
							hasAuth = true
							break
						}
					}
				}
			}
			_, _ = conn.Do("SADD", args...)
		}
		_, _ = conn.Do("EXPIRE", roleResourceKey, 3*24*60*60)
		if !dBCodeOp {
			for _, resourceCode := range resourceCodesInRole {
				if resourceCode == constant.ResourceAllProjectDetail {
					hasAuth = true
				} else if _, ok := codeMap[resourceCode]; ok {
					hasAuth = true
					break
				} else {
					for _, code := range codes {
						if strings.HasPrefix(code, resourceCode+":") {
							hasAuth = true
							break
						}
					}
				}
			}
		}
		if hasAuth {
			break
		}
	}
	return hasAuth, nil
}

func CheckResourceCodeInProject(ctx *fiber.Ctx, projectID uint64, codes ...string) (uid uint64, success bool, err error) {
	uid, err = GetCurrentUserID(ctx)
	if err != nil {
		logger.Errorln(err)
		err = ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
		return
	}
	var hasPermit bool
	hasPermit, err = HasResourceCodeInProject(uid, projectID, codes...)
	if err != nil {
		logger.Errorln(err)
		err = ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
		return
	}
	if !hasPermit {
		ctx.Status(fiber.StatusForbidden)
		err = ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch,
		})
		return
	}
	success = true
	return
}
