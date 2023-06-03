package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	"github.com/yockii/celestial/internal/module/uc/domain"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/celestial/pkg/crypto"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/shared"
	"github.com/yockii/ruomu-core/util"
)

var UserController = new(userController)

type userController struct{}

func (*userController) Register(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	instance.Username = strings.TrimSpace(instance.Username)
	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " username or password",
		})
	}
	isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
	if !isStrong {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodePasswordStrengthInvalid,
			Msg:  server.ResponseMsgPasswordStrengthInvalid,
		})
	}

	// 创建用户
	duplicated, success, err := service.UserService.Add(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if duplicated {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		})
	}

	if !success {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
	}

	instance.Password = ""
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *userController) Login(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Username == "" || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": 用户名及密码",
		})
	}

	// 解析密码
	if pwd, err := crypto.Sm2Decrypt(instance.Password); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + "密码不准确",
		})
	} else {
		instance.Password = pwd
	}

	isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
	if !isStrong {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodePasswordStrengthInvalid,
			Msg:  server.ResponseMsgPasswordStrengthInvalid,
		})
	}

	user, notMatch, err := service.UserService.LoginWithUsernameAndPassword(instance.Username, instance.Password)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if notMatch {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  "用户名与密码" + server.ResponseMsgDataNotMatch,
		})
	}
	return c.generateLoginResponse(user, ctx)
}

func (c *userController) GetUserRoleIds(ctx *fiber.Ctx) error {
	uid, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if uid == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user id",
		})
	}
	// 获取用户对应的权限和角色
	roles, err := service.UserService.Roles(uid)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	var roleIds []string
	for _, role := range roles {
		roleIds = append(roleIds, strconv.FormatUint(role.ID, 10))
	}
	return ctx.JSON(&server.CommonResponse{
		Data: roleIds,
	})
}

func (c *userController) Add(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Username == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " username",
		})
	}

	if instance.Password != "" {
		isStrong := util.PasswordStrengthCheck(8, 50, 4, instance.Password)
		if !isStrong {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodePasswordStrengthInvalid,
				Msg:  server.ResponseMsgPasswordStrengthInvalid,
			})
		}
	}

	duplicated, success, err := service.UserService.Add(instance)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if duplicated {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		})
	}
	if !success {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *userController) UpdateUser(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	// 处理必填
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user id",
		})
	}

	success, err := service.UserService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	// 处理必填
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user id",
		})
	}

	success, err := service.UserService.Delete(instance.ID)

	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *userController) Instance(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 {
		// 获取当前登录的用户
		uidStr, _ := ctx.Locals(shared.JwtClaimUserId).(string)
		if uidStr == "" {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " user id",
			})
		}
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamParseError,
				Msg:  server.ResponseMsgParamParseError,
			})
		}
		instance.ID = uid
	}
	user, err := service.UserService.Instance(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: user,
	})
}

func (c *userController) List(ctx *fiber.Ctx) error {
	instance := new(domain.UserListRequest)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	paginate := new(server.Paginate)
	if err := ctx.QueryParser(paginate); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if paginate.Limit == 0 {
		paginate.Limit = 10
	}

	tcList := make(map[string]*server.TimeCondition)
	if instance.CreateTimeCondition != nil {
		tcList["create_time"] = instance.CreateTimeCondition
	}
	if instance.UpdateTimeCondition != nil {
		tcList["update_time"] = instance.UpdateTimeCondition
	}

	total, list, err := service.UserService.PaginateBetweenTimes(&instance.User, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Items:  list,
			Total:  total,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

func generateJwtToken(userId, tenantId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	sid := util.GenerateXid()

	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	_, err := conn.Do("SETEX", shared.RedisSessionIdKey+sid, config.GetInt("userTokenExpire"), userId)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	claims[shared.JwtClaimUserId] = userId
	claims[shared.JwtClaimTenantId] = tenantId
	claims[shared.JwtClaimSessionId] = sid

	t, err := token.SignedString([]byte(shared.JwtSecret))
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	return t, nil
}

// DispatchRoles 给用户分配角色
func (c *userController) DispatchRoles(ctx *fiber.Ctx) error {
	instance := new(domain.UserDispatchRolesRequest)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	// 处理必填
	if instance.UserID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user id",
		})
	}

	success, err := service.UserService.DispatchRoles(instance.UserID, instance.RoleIDList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		// 成功处理后，清除用户的权限缓存
		conn := cache.Get()
		defer func(conn redis.Conn) {
			_ = conn.Close()
		}(conn)
		key := shared.RedisKeyUserRoles + strconv.FormatUint(instance.UserID, 10)
		_, err = conn.Do("DEL", key)
		if err != nil {
			logger.Errorln(err)
		}
		// 将新的用户角色ID缓存到redis
		args := make([]interface{}, len(instance.RoleIDList)+1)
		args[0] = key
		for i, v := range instance.RoleIDList {
			args[i+1] = v
		}
		_, err = conn.Do("SADD", args...)
		if err != nil {
			logger.Errorln(err)
		}
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

// UpdateSelf 更新自己的信息
func (c *userController) UpdateSelf(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	// 获取当前登录的用户ID
	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	instance.ID = uid
	success, err := service.UserService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *userController) LoginByDingtalkCode(ctx *fiber.Ctx) error {
	req := new(domain.LoginByDingtalkCodeRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if req.Code == "" || req.ThirdSourceID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " code or third source id",
		})
	}
	if thirdSource, err := service.ThirdSourceService.Instance(&model.ThirdSource{ID: req.ThirdSourceID}); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	} else if thirdSource == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	} else {
		var user *model.User
		user, err = service.DingtalkService.SyncDingUserByThirdSourceOutsideDingtalk(thirdSource, req.Code, true)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		return c.generateLoginResponse(user, ctx)
	}
}

func (c *userController) LoginInDingtalk(ctx *fiber.Ctx) error {
	req := new(domain.LoginByDingtalkCodeRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if req.Code == "" || req.ThirdSourceID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " code or third source id",
		})
	}
	if thirdSource, err := service.ThirdSourceService.Instance(&model.ThirdSource{ID: req.ThirdSourceID}); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	} else if thirdSource == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	} else {
		staffId := ""
		staffId, err = dingtalk.GetStaffIdByCode(thirdSource, req.Code)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeUnknownError,
				Msg:  server.ResponseMsgUnknownError + err.Error(),
			})
		}
		var user *model.User
		user, err = service.DingtalkService.SyncDingUserByThirdSource(thirdSource, staffId, true)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}

		return c.generateLoginResponse(user, ctx)
	}
}

func (c *userController) generateLoginResponse(user *model.User, ctx *fiber.Ctx) error {
	jwtToken, err := generateJwtToken(strconv.FormatUint(user.ID, 10), "")
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  server.ResponseMsgGeneration + err.Error(),
		})
	}
	user.Password = ""
	return ctx.JSON(&server.CommonResponse{
		Data: map[string]interface{}{
			"token": jwtToken,
			"user":  user,
		},
	})
}
