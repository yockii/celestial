package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/asset/domain"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/celestial/internal/module/asset/service"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	ucService "github.com/yockii/celestial/internal/module/uc/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
)

var AssetFileController = new(assetFileController)

type assetFileController struct{}

func (c *assetFileController) Add(ctx *fiber.Ctx) error {
	if form, err := ctx.MultipartForm(); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	} else {
		instance := new(model.File)
		instance.CreatorID, err = helper.GetCurrentUserID(ctx)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " user Id",
			})
		}
		// 资产目录
		if categoryIdList, has := form.Value["categoryId"]; !has {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " categoryId",
			})
		} else {
			categoryId := categoryIdList[0]
			instance.CategoryID, err = strconv.ParseUint(categoryId, 10, 64)
			if err != nil {
				logger.Errorln(err)
				return ctx.JSON(&server.CommonResponse{
					Code: server.ResponseCodeParamParseError,
					Msg:  server.ResponseMsgParamParseError + " categoryId",
				})
			}
		}
		// 资产名称
		if nameList, has := form.Value["name"]; !has {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " name",
			})
		} else {
			instance.Name = nameList[0]
		}
		if instance.CategoryID == 0 || instance.Name == "" {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " categoryId or name",
			})
		}

		if fhList, ok := form.File["file"]; !ok {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " file",
			})
		} else {
			if len(fhList) > 0 {
				fileHeader := fhList[0]
				duplicated := false
				success := false
				instance.Size = fileHeader.Size
				fileName := fileHeader.Filename
				instance.OriginName = fileName
				instance.Suffix = fileName[strings.LastIndex(fileName, ".")+1:]

				instance.CreatorID, err = helper.GetCurrentUserID(ctx)
				if err != nil {
					return ctx.JSON(&server.CommonResponse{
						Code: server.ResponseCodeParamNotEnough,
						Msg:  server.ResponseMsgParamNotEnough + " user Id",
					})
				}

				var file multipart.File
				file, err = fileHeader.Open()
				if err != nil {
					logger.Errorln(err)
					return ctx.JSON(&server.CommonResponse{
						Code: server.ResponseCodeParamParseError,
						Msg:  server.ResponseMsgParamParseError,
					})
				}

				defer func(file multipart.File) {
					_ = file.Close()
				}(file)

				duplicated, success, err = service.AssetFileService.Add(instance, file)
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
			}
		}

		_ = ants.Submit(data.AddDocumentAntsWrapper(&search.Document{
			ID:         instance.ID,
			Title:      instance.Name,
			Content:    instance.OriginName,
			Route:      fmt.Sprintf("/asset/file?id=%d", instance.ID),
			CreateTime: instance.CreateTime,
			UpdateTime: instance.UpdateTime,
		}, instance.CreatorID))

		return ctx.JSON(&server.CommonResponse{
			Data: instance,
		})
	}
}

func (c *assetFileController) Update(ctx *fiber.Ctx) error {
	instance := new(model.File)
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
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}

	success, err := service.AssetFileService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			return func() {
				d, e := service.AssetFileService.Instance(id)
				if e != nil {
					logger.Errorln(err)
					return
				}
				_ = data.AddDocument(&search.Document{
					ID:         d.ID,
					Title:      d.Name,
					Content:    d.OriginName,
					Route:      fmt.Sprintf("/asset/file?id=%d", d.ID),
					CreateTime: d.CreateTime,
					UpdateTime: d.UpdateTime,
				}, d.CreatorID)
			}
		}(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *assetFileController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.File)
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
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}

	success, err := service.AssetFileService.Delete(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(data.DeleteDocumentsAntsWrapper(instance.ID))
	}
	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *assetFileController) List(ctx *fiber.Ctx) error {
	instance := new(domain.AssetFileListRequest)
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

	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	total, list, err := service.AssetFileService.PaginateDomainListBetweenTimes(&instance.File, uid, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	//total, list, err := service.AssetFileService.PaginateBetweenTimes(&instance.File, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.AssetFileWithCreator
	{
		var wg sync.WaitGroup
		for _, item := range list {
			//fwu := &domain.AssetFileWithCreator{
			//	File: *item,
			//}
			fwu := item
			resultList = append(resultList, fwu)
			wg.Add(1)
			go func(result *domain.AssetFileWithCreator) {
				defer wg.Done()
				user, err := ucService.UserService.Instance(&ucModel.User{ID: result.CreatorID})
				if err != nil {
					logger.Errorln(err)
					return
				}
				result.Creator = &ucModel.User{
					ID:       user.ID,
					Username: user.Username,
					RealName: user.RealName,
				}
			}(fwu)
		}
		wg.Wait()
	}

	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  resultList,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

func (c *assetFileController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.File)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if condition.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	dept, err := service.AssetFileService.Instance(condition.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: dept,
	})
}

func (c *assetFileController) Download(ctx *fiber.Ctx) error {
	instance := new(model.File)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if instance.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		})
	}
	fileReader, err := service.AssetFileService.Download(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.SendStream(fileReader)
}

func (c *assetFileController) FilePermissionUsers(ctx *fiber.Ctx) error {
	filePermission := new(domain.FilePermissionUser)
	if err := ctx.QueryParser(filePermission); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	list, err := service.AssetFileService.GetPermissionUsers(filePermission)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: list,
	})
}

func (c *assetFileController) UpdateFileUserPermission(ctx *fiber.Ctx) error {
	filePermission := new(model.FilePermission)
	if err := ctx.BodyParser(filePermission); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 校验uid是否有fileId的管理权限
	has, err := service.AssetFileService.CheckPermission(filePermission.FileID, uid, model.FilePermissionManage)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if !has {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch,
		})
	}
	// 更新权限
	err = service.AssetFileService.UpdateFileUserPermission(filePermission)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: true,
	})
}
