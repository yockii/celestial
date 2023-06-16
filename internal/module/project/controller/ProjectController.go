package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/project/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
	"strings"
	"sync"
)

var ProjectController = new(projectController)

type projectController struct{}

func (_ *projectController) Add(ctx *fiber.Ctx) error {
	instance := new(model.Project)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Name == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name",
		})
	}

	// 当前登录用户作为项目拥有者
	// 获取当前登录的用户ID
	uid, err := helper.GetCurrentUserID(ctx)
	if err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	instance.OwnerID = uid

	instance.Name = strings.TrimSpace(instance.Name)

	duplicated, success, err := service.ProjectService.Add(instance)
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

	_ = ants.Submit(data.AddDocumentAntsWrapper(&search.Document{
		ID:           instance.ID,
		Title:        instance.Name,
		Content:      instance.Description,
		Route:        fmt.Sprintf("/project/detail/%d", instance.ID),
		RelatedUsers: nil,
		CreateTime:   instance.CreateTime,
		UpdateTime:   instance.UpdateTime,
	}, instance.OwnerID))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (_ *projectController) Update(ctx *fiber.Ctx) error {
	instance := new(model.Project)
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
	instance.Name = strings.TrimSpace(instance.Name)
	success, err := service.ProjectService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		_ = ants.Submit(func(id uint64) func() {
			d, e := service.ProjectService.Instance(id)
			if e != nil {
				logger.Errorln(e)
				return func() {
				}
			}
			_, members, e := service.ProjectMemberService.PaginateBetweenTimes(&model.ProjectMember{ProjectID: id}, -1, -1, "", nil)
			if e != nil {
				logger.Errorln(e)
				return func() {
				}
			}
			relatedIdList := []uint64{d.OwnerID}
			for _, member := range members {
				relatedIdList = append(relatedIdList, member.UserID)
			}
			return data.AddDocumentAntsWrapper(&search.Document{
				ID:           d.ID,
				Title:        d.Name,
				Content:      d.Description,
				Route:        fmt.Sprintf("/project/detail/%d", d.ID),
				RelatedUsers: nil,
				CreateTime:   d.CreateTime,
				UpdateTime:   d.UpdateTime,
			}, relatedIdList...)
		}(instance.ID))
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (_ *projectController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.Project)
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

	success, err := service.ProjectService.Delete(instance.ID)
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

func (_ *projectController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectListRequest)
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

	total, list, err := service.ProjectService.PaginateBetweenTimes(&instance.Project, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.ProjectWithMembers
	// 获取关联数据
	var wg sync.WaitGroup
	for _, project := range list {
		p := &domain.ProjectWithMembers{
			Project: *project,
		}
		resultList = append(resultList, p)
		wg.Add(1)
		go func(proj *domain.ProjectWithMembers) {
			defer wg.Done()
			// 获取项目成员
			members, err := service.ProjectMemberService.ListLiteByProjectID(proj.ID)
			if err != nil {
				logger.Errorln(err)
				return
			}
			proj.Members = members
		}(p)
	}
	wg.Wait()

	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  resultList,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

// Instance 获取角色详情
func (_ *projectController) Instance(ctx *fiber.Ctx) error {
	instance := new(model.Project)
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

	instance, err := service.ProjectService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	result := &domain.ProjectWithMembers{}
	if instance != nil {
		result.Project = *instance
		result.Members, _ = service.ProjectMemberService.ListLiteByProjectID(instance.ID)
	}
	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}

// StatisticsByStage 统计项目阶段
func (_ *projectController) StatisticsByStage(ctx *fiber.Ctx) error {
	result, err := service.ProjectService.StatisticsByStage()
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}
