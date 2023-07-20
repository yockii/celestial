package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/core/helper"
	assetModel "github.com/yockii/celestial/internal/module/asset/model"
	assetService "github.com/yockii/celestial/internal/module/asset/service"
	"github.com/yockii/celestial/internal/module/onlyoffice/domain"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	ucService "github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strconv"
)

var OnlyOfficeController = new(onlyOfficeController)

type onlyOfficeController struct{}

func (c *onlyOfficeController) Callback(ctx *fiber.Ctx) error {
	j := gjson.ParseBytes(ctx.Body())
	status := j.Get("status").Int() // 1：编辑中 2：已准备好保存 3：保存发生错误 4：无改变关闭 6：正在编辑文档但执行了保存 7：强制保存时出错
	if status == 2 || status == 6 {
		// 正在编辑文档但保存了当前文档状态
		fileVersionId := j.Get("key").Uint()
		downloadUrl := j.Get("url").String()
		versionedFile, err := assetService.AssetFileService.GetFileVersion(&assetModel.FileVersion{
			ID: fileVersionId,
		})
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&domain.Response{
				Error: 7,
			})
		}
		file, err := assetService.AssetFileService.Instance(versionedFile.FileID)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&domain.Response{
				Error: 7,
			})
		}

		// 下载文档文件并保存
		resp, err := http.Get(downloadUrl)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&domain.Response{
				Error: 7,
			})
		}
		defer resp.Body.Close()
		// 保存文件
		file.Size = resp.ContentLength
		file.CreatorID = j.Get("users.0").Uint()
		var success bool
		success, err = assetService.AssetFileService.AddNewVersion(file, resp.Body)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&domain.Response{
				Error: 7,
			})
		}
		// 成功返回 error:0 出错的话，返回error: 7
		if !success {
			return ctx.JSON(&domain.Response{
				Error: 7,
			})
		}
	} else if status == 7 {
		// 强制保存文档出错
		return ctx.JSON(&domain.Response{
			Error: 7,
		})
	}
	return ctx.JSON(&domain.Response{})
}

func (c *onlyOfficeController) GetConfig(ctx *fiber.Ctx) error {
	var err error
	fileVersion := new(assetModel.FileVersion)
	if err = ctx.QueryParser(fileVersion); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if fileVersion.ID == 0 && fileVersion.FileID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": id",
		})
	}

	baseUrl := ctx.Query("baseUrl")
	if baseUrl == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + ": base url",
		})
	}
	var uid uint64
	uid, err = helper.GetCurrentUserID(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}
	var user *ucModel.User
	user, err = ucService.UserService.Instance(&ucModel.User{ID: uid})
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	if user == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	var file *assetModel.File
	file, err = assetService.AssetFileService.Instance(fileVersion.FileID)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	if file == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	// 获取file permission
	var filePermission *assetModel.FilePermission
	filePermission, err = assetService.AssetFilePermissionService.Instance(&assetModel.FilePermission{
		FileID: file.ID,
		UserID: uid,
	})
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	if filePermission == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	if fileVersion.ID == 0 {
		fileVersion, err = assetService.AssetFileService.LatestVersion(file.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果是未找到，则尝试将当前file信息放入version并存储
				// 添加FileVersion
				fileVersion = &assetModel.FileVersion{
					ID:          util.SnowflakeId(),
					FileID:      file.ID,
					OssConfigID: file.OssConfigID,
					Size:        file.Size,
					ObjName:     file.ObjName,
					CreatorID:   file.CreatorID,
				}
				err = database.DB.Create(fileVersion).Error
				if err != nil {
					logger.Errorln(err)
					return ctx.JSON(&server.CommonResponse{
						Code: server.ResponseCodeDatabase,
						Msg:  server.ResponseMsgDatabase,
					})
				}
			} else {
				logger.Errorln(err)
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeDatabase,
					Msg:  server.ResponseMsgDatabase,
				})
			}
		}
	} else {
		fileVersion, err = assetService.AssetFileService.GetFileVersion(fileVersion)
		if err != nil {
			logger.Errorln(err)
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase,
			})
		}
	}
	if fileVersion == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	//构造配置信息
	openOfficeConfig := new(domain.EditConfig)
	ft := translateFileType(file.Suffix)
	if ft == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ": " + file.Suffix,
		})
	}
	openOfficeConfig.Document = domain.Document{
		FileType: file.Suffix,
		Key:      strconv.FormatUint(fileVersion.ID, 10),
		Title:    file.Name,
		Url:      baseUrl + "/api/v1/office/download?id=" + strconv.FormatUint(fileVersion.ID, 10),
		Permissions: domain.Permissions{
			Comment: filePermission.Permission >= 1,
			CommentGroups: domain.CommentGroups{
				Edit:   []string{},
				Remove: []string{},
				View:   "",
			},
			Download:              filePermission.Permission >= 3,
			Edit:                  filePermission.Permission >= 2,
			EditCommentAuthorOnly: true,
		},
	}
	openOfficeConfig.DocumentType = ft
	openOfficeConfig.Type = domain.EditConfigTypeDesktop // TODO: 2021/8/17 临时写死，后面让页面自己去判断？
	openOfficeConfig.EditorConfig = domain.Editor{
		Mode: domain.EditConfigModeEdit, // TODO: 2021/8/17 临时写死，后面根据权限进行判断
		Lang: domain.EditConfigLangZhCn,
		User: domain.User{
			ID:   uid,
			Name: user.RealName,
		},
		CallbackUrl: baseUrl + "/api/v1/office/callback",
	}

	// Token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = openOfficeConfig
	openOfficeConfig.Token, err = token.SignedString([]byte(config.GetString("onlyoffice.secret")))
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
	}

	return ctx.JSON(&server.CommonResponse{
		Data: openOfficeConfig,
	})
}

func (c *onlyOfficeController) Download(ctx *fiber.Ctx) error {
	tokenStr := ctx.Get(fiber.HeaderAuthorization)
	tokenStr = tokenStr[len("Bearer "):]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetString("onlyoffice.secret")), nil
	})
	if err != nil {
		logger.Errorln(err)
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeUnknownError,
			Msg:  server.ResponseMsgUnknownError,
		})
	}
	claims := token.Claims.(jwt.MapClaims)
	payload, ok := claims["payload"].(map[string]interface{})
	if !ok {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}
	var u string
	u, ok = payload["url"].(string)
	if !ok {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}
	var uu *url.URL
	uu, err = url.Parse(u)
	if err != nil {
		logger.Errorln(err)
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	q := uu.Query()
	fileVersionId := q.Get("id")

	var id uint64
	id, err = strconv.ParseUint(fileVersionId, 10, 64)
	if err != nil {
		logger.Errorln(err)
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	fileReader, err := assetService.AssetFileService.DownloadVersion(id)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.SendStream(fileReader)
}

func translateFileType(suffix string) string {
	switch suffix {
	case "doc", "docx":
		return domain.DocumentTypeWord
	case "xls", "xlsx":
		return domain.DocumentTypeCell
	case "ppt", "pptx":
		return domain.DocumentTypeSlide
	default:
		return ""
	}
}
