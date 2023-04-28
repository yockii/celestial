package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/dingtalk"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/server"
	"strconv"
)

var ThirdLoginController = new(thirdLoginController)

type thirdLoginController struct{}

// DingtalkCallback 钉钉回调
func (c *thirdLoginController) DingtalkCallback(ctx *fiber.Ctx) error {
	sourceID, _ := strconv.ParseUint(ctx.Params("id"), 10, 64)

	source, err := service.ThirdSourceService.Instance(&model.ThirdSource{ID: sourceID})
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}

	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	body := ctx.Body()
	req := gjson.ParseBytes(body)
	if req.Get("encrypt").Exists() {
		encrypt := req.Get("encrypt").String()
		if encrypt != "" {
			decrypted := ""
			decrypted, err = dingtalk.GetDecryptMsg(
				source,
				signature,
				timestamp,
				nonce,
				encrypt,
			)
			if err != nil {
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeGeneration,
					Msg:  server.ResponseMsgGeneration,
				})
			}
			dj := gjson.Parse(decrypted)
			eventType := dj.Get("EventType").String()
			switch eventType {
			case "check_url":
				logger.Debug("测试回调url正确性")
			case "user_add_org": // 通讯录用户增加
				fallthrough
			case "user_modify_org": // 通讯录用户更改
				fallthrough
			case "user_active_org": // 加入企业后用户激活
				// 通讯录用户同步...
				uidList := dj.Get("UserId").Array()
				var uids []string
				for _, uid := range uidList {
					uids = append(uids, uid.String())
				}
				corpId := dj.Get("CorpId").String()
				err = userUpdated(corpId, uids)

			case "user_leave_org": // 通讯录用户离职
				// 离职
				uidList := dj.Get("UserId").Array()
				var uids []string
				for _, uid := range uidList {
					uids = append(uids, uid.String())
				}
				corpId := dj.Get("CorpId").String()
				err = userRemoved(corpId, uids)

			case "org_dept_create": // 通讯录部门创建
				fallthrough
			case "org_dept_modify": // 通讯录部门修改
				didList := dj.Get("DeptId").Array()
				var dids []string
				for _, uid := range didList {
					dids = append(dids, uid.String())
				}
				corpId := dj.Get("CorpId").String()
				err = deptUpdated(corpId, dids)
			case "org_dept_remove": // 通讯录部门删除
				didList := dj.Get("DeptId").Array()
				var dids []string
				for _, uid := range didList {
					dids = append(dids, uid.String())
				}
				err = deptRemoved(dids)
			}
			if err != nil {
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeGeneration,
					Msg:  server.ResponseMsgGeneration,
				})
			}
			m, err := dingtalk.GetEncryptMsg(source, "success")
			if err != nil {
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeGeneration,
					Msg:  server.ResponseMsgGeneration,
				})
			}
			return ctx.JSON(m)
		}
	}
	return ctx.JSON(&server.CommonResponse{
		Code: server.ResponseCodeUnknownError,
		Msg:  "未知错误",
	})
}

func deptRemoved(dids []string) error {
	for _, did := range dids {
		err := service.DepartmentService.DeptRemoved(did)
		if err != nil {
			return err
		}
	}
	return nil
}

func deptUpdated(corpId string, dids []string) error {
	for _, did := range dids {
		_, err := service.DingtalkService.SyncDingDept(corpId, did, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func userRemoved(corpId string, uids []string) error {
	for _, uid := range uids {
		err := service.UserService.UserLeaved(corpId, uid)
		if err != nil {
			return err
		}
	}
	return nil
}

func userUpdated(corpId string, uids []string) error {
	for _, uid := range uids {
		_, err := service.DingtalkService.SyncDingUser(corpId, uid, true)
		if err != nil {
			return err
		}
	}
	return nil
}
