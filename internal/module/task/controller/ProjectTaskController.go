package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/task/domain"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/celestial/internal/module/task/service"
	"github.com/yockii/ruomu-core/server"
	"sync"
)

var ProjectTaskController = new(projectTaskController)

type projectTaskController struct{}

func (c *projectTaskController) Add(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectTaskWithMembers)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.Name == "" || instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name & projectId",
		})
	}

	// 获取当前用户id作为ownerId和creatorId
	if currentUserId, err := helper.GetCurrentUserID(ctx); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " user info",
		})
	} else {
		instance.OwnerID = currentUserId
		instance.CreatorID = currentUserId
	}

	duplicated, success, err := service.ProjectTaskService.Add(&instance.ProjectTask, instance.Members)
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
	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectTaskController) Update(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectTaskWithMembers)
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

	success, err := service.ProjectTaskService.Update(&instance.ProjectTask, instance.Members)
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

func (c *projectTaskController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectTask)
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

	success, err := service.ProjectTaskService.Delete(instance.ID)
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

func (c *projectTaskController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectTaskListTask)
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

	total, list, err := service.ProjectTaskService.PaginateBetweenTimes(&instance.ProjectTask, instance.OnlyParent, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	var resultList []*domain.ProjectTaskWithMembers
	if len(list) > 0 {
		var wg sync.WaitGroup

		for _, item := range list {
			wg.Add(1)
			ptwm := new(domain.ProjectTaskWithMembers)
			ptwm.ProjectTask = *item
			resultList = append(resultList, ptwm)
			go func(task *domain.ProjectTaskWithMembers) {
				defer wg.Done()
				members, err := service.ProjectTaskMemberService.ListWithRealName(&model.ProjectTaskMember{TaskID: task.ID})
				if err != nil {
					logger.Errorln(err)
					return
				}
				task.Members = members
			}(ptwm)
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

func (c *projectTaskController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectTask)
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
	dept, err := service.ProjectTaskService.Instance(condition.ID)
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

// TaskDurationByProject 项目下任务耗时统计
func (c *projectTaskController) TaskDurationByProject(ctx *fiber.Ctx) error {
	condition := new(domain.ProjectTaskListTask)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 项目ID必须传入
	if condition.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId",
		})
	}

	tcList := make(map[string]*server.TimeCondition)
	if condition.CreateTimeCondition != nil {
		tcList["create_time"] = condition.CreateTimeCondition
	}
	if condition.UpdateTimeCondition != nil {
		tcList["update_time"] = condition.UpdateTimeCondition
	}

	// 获取预计工时和实际工时的统计
	result, err := service.ProjectTaskService.TaskDurationByProject(condition.ProjectID, tcList)
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
