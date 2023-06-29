package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	"github.com/yockii/celestial/internal/module/uc/domain"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/celestial/pkg/crypto"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/server"
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
		uidStr, _ := ctx.Locals(constant.JwtClaimUserId).(string)
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

	// 数据脱敏
	for _, user := range list {
		// 确保用户密码不被传递
		user.Password = ""
		// 确保用户手机隐藏中间4位
		if user.Mobile != "" {
			if ml := len(user.Mobile); ml <= 10 {
				if ml <= 3 {
					user.Mobile = "****"
				} else {
					user.Mobile = user.Mobile[:3] + "****"
				}
			} else {
				user.Mobile = user.Mobile[:3] + "****" + user.Mobile[ml-4:]
			}
		}
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
	sessionKey := fmt.Sprintf("%s:%s", constant.RedisSessionIdKey, sid)

	_, err := conn.Do("SETEX", sessionKey, config.GetInt("userTokenExpire"), userId)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	claims[constant.JwtClaimUserId] = userId
	claims[constant.JwtClaimTenantId] = tenantId
	claims[constant.JwtClaimSessionId] = sid

	t, err := token.SignedString([]byte(constant.JwtSecret))
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	return t, nil
}

// AssignRole 给用户分配角色
func (c *userController) AssignRole(ctx *fiber.Ctx) error {
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
		key := fmt.Sprintf("%s:%d", constant.RedisKeyUserRoles, instance.UserID)
		_, err = conn.Do("DEL", key)
		if err != nil {
			logger.Errorln(err)
		}
		// 删除即可，等待中间件重新加载
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

func (c *userController) UserRoleIdList(ctx *fiber.Ctx) error {
	userIdStr := ctx.Query("id")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError + " id",
		})
	}
	if userId == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	roleList, err := service.UserService.Roles(userId)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	var roleIdList []string
	for _, role := range roleList {
		roleIdList = append(roleIdList, strconv.FormatUint(role.ID, 10))
	}
	return ctx.JSON(&server.CommonResponse{
		Data: roleIdList,
	})
}

// UserPermissions 获取用户权限
func (c *userController) UserPermissions(ctx *fiber.Ctx) error {
	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError + err.Error(),
		})
	}
	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	// 获取用户数据权限
	userDataPerm := 0
	// 获取用户角色
	userRolesKey := fmt.Sprintf("%s:%d", constant.RedisKeyUserRoles, uid)
	roleIds, _ := redis.Uint64s(conn.Do("SMEMBERS", userRolesKey))
	if len(roleIds) == 0 {
		// 获取该用户的角色id存入缓存
		var roles []*model.Role
		roles, err = service.UserService.Roles(uid)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		for _, role := range roles {
			roleIds = append(roleIds, role.ID)
			_, _ = conn.Do("SADD", userRolesKey, role.ID)
			if userDataPerm == 0 || role.DataPermission < userDataPerm {
				userDataPerm = role.DataPermission
			}
		}
	}
	_, _ = conn.Do("EXPIRE", userRolesKey, 3*24*60*60)

	// 获取用户资源编码列表
	resourceCodes := make(map[string]struct{})
	isSuperAdmin := false
	for _, roleId := range roleIds {
		if roleId == constant.SuperAdminRoleId {
			isSuperAdmin = true
			break
		}

		roleResourcesKey := fmt.Sprintf("%s:%d", constant.RedisKeyRoleResourceCode, roleId)
		codes, _ := redis.Strings(conn.Do("SMEMBERS", roleResourcesKey))
		if len(codes) == 0 {
			// 缓存没有，那么就去数据库取出来放进去
			codes, err = service.RoleService.ResourceCodes(roleId)
			if err != nil {
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase + err.Error(),
				})
			}
			for _, resourceCode := range codes {
				rc := resourceCode
				_, _ = conn.Do("SADD", roleResourcesKey, rc)
				if _, ok := resourceCodes[rc]; !ok {
					resourceCodes[rc] = struct{}{}
				}
			}
		}
		_, _ = conn.Do("EXPIRE", roleResourcesKey, 3*24*60*60)
	}
	result := new(domain.UserResourceCodesResponse)
	if isSuperAdmin {
		result.IsSuperAdmin = true
	} else {
		result.ResourceCodeList = make([]string, 0, len(resourceCodes))
		for code := range resourceCodes {
			result.ResourceCodeList = append(result.ResourceCodeList, code)
		}
	}
	result.DataPermission = userDataPerm
	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}

// ResetUserPassword 重置用户密码
func (c *userController) ResetUserPassword(ctx *fiber.Ctx) error {
	instance := new(model.User)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	if instance.ID == 0 || instance.Password == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
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

	// 重置密码
	success, err := service.UserService.UpdatePassword(instance)
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
