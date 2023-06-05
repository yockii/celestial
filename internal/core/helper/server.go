package helper

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
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
