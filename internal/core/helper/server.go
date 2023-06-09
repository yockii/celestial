package helper

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/project/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/shared"
	"strconv"
)

func GetCurrentUserID(ctx *fiber.Ctx) (uint64, error) {
	// 获取当前登录的用户
	uidStr, _ := ctx.Locals(shared.JwtClaimUserId).(string)
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
	permit, err = redis.Int(conn.Do("HGET", constant.RedisKeyUserDataPermitInProject+strconv.FormatUint(uid, 10), strconv.FormatUint(projectID, 10)))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		logger.Errorln(err)
		return 0, err
	}
	// 如果不存在或者为0，那么就去数据库中查找
	if permit != 0 {
		return permit, nil
	}

	// 查询uid对应在项目中的角色id
	projectMember := &model.ProjectMember{ProjectID: projectID, UserID: uid}
	if err = database.DB.Where(projectMember).First(projectMember).Error; err != nil {
		logger.Errorln(err)
		return 0, err
	}
	// 从缓存中查询角色对应的数据权限
	permit, err = redis.Int(conn.Do("HGET", constant.RedisKeyRoleDataPerm, projectMember.RoleID))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		logger.Errorln(err)
		return 0, err
	}
	// 如果不存在或者为0，那么就去数据库中查找
	if permit != 0 {
		return permit, nil
	}
	// 查询角色对应的数据权限
	role := &ucModel.Role{ID: projectMember.RoleID}
	if err = database.DB.Where(role).First(role).Error; err != nil {
		logger.Errorln(err)
		return 0, err
	}
	permit = role.DataPermission
	// 存入缓存中
	_, _ = conn.Do("HSET", constant.RedisKeyRoleDataPerm, projectMember.RoleID, permit)
	// 设置失效时间
	_, _ = conn.Do("EXPIRE", constant.RedisKeyRoleDataPerm, 3*24*60*60)
	return permit, nil
}
