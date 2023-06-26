package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/task/domain"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/celestial/internal/module/task/service"
	"github.com/yockii/celestial/pkg/search"
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

	// 判断权限
	if uid, success, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectTaskAdd); err != nil {
		return err
	} else if !success {
		return nil
	} else {
		instance.OwnerID = uid
		instance.CreatorID = uid
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

	_ = ants.Submit(func(d *domain.ProjectTaskWithMembers) func() {
		var relatedUids []uint64
		ruMap := make(map[uint64]struct{})
		ruMap[d.OwnerID] = struct{}{}
		ruMap[d.CreatorID] = struct{}{}
		for _, member := range d.Members {
			ruMap[member.UserID] = struct{}{}
		}
		for uid := range ruMap {
			relatedUids = append(relatedUids, uid)
		}

		p := "低"
		switch d.Priority {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		status := "未开始"
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.Name,
			Content: fmt.Sprintf("[%s]: %s, 状态：%s",
				p,

				d.TaskDesc,
				status,
			),
			Route:      fmt.Sprintf("/project/detail/%d/task?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, relatedUids...)
	}(instance))

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

	// 先取出原始数据
	oldInstance, err := service.ProjectTaskService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if oldInstance == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTaskUpdate); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectTaskService.Update(&instance.ProjectTask, oldInstance, instance.Members)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		c.addSearchDocument(instance.ID)
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

	// 先取出原始数据
	oldInstance, err := service.ProjectTaskService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if oldInstance == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	var success bool
	if _, success, err = helper.CheckResourceCodeInProject(ctx, oldInstance.ProjectID, constant.ResourceProjectTaskDelete); err != nil {
		return err
	} else if !success {
		return nil
	}

	success, err = service.ProjectTaskService.Delete(oldInstance)
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

func (c *projectTaskController) addSearchDocument(id uint64) {
	_ = ants.Submit(func(id uint64) func() {
		d, e := service.ProjectTaskService.Instance(id)
		if e != nil {
			logger.Errorln(e)
			return func() {}
		}
		members, err := service.ProjectTaskMemberService.List(&model.ProjectTaskMember{TaskID: id})
		if err != nil {
			logger.Errorln(err)
			return func() {}
		}

		var relatedUids []uint64
		ruMap := make(map[uint64]struct{})
		ruMap[d.OwnerID] = struct{}{}
		ruMap[d.CreatorID] = struct{}{}
		for _, member := range members {
			ruMap[member.UserID] = struct{}{}
		}
		for uid := range ruMap {
			relatedUids = append(relatedUids, uid)
		}

		p := "低"
		switch d.Priority {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		status := "未开始"
		//任务状态 -1-已取消 1-未开始 2-已确认 3-进行中 9-已完成
		switch d.Status {
		case -1:
			status = "已取消"
		case 1:
			status = "未开始"
		case 2:
			status = "已确认"
		case 3:
			status = "进行中"
		case 9:
			status = "已完成"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.Name,
			Content: fmt.Sprintf("[%s]: %s, 状态：%s",
				p,
				d.TaskDesc,
				status,
			),
			Route:      fmt.Sprintf("/project/detail/%d/task?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, relatedUids...)
	}(id))
}

func (c *projectTaskController) MemberUpdateStatus(status int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		condition := new(model.ProjectTask)
		if err := ctx.BodyParser(condition); err != nil {
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

		// 先取出原始数据
		oldTask, err := service.ProjectTaskService.Instance(condition.ID)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		if oldTask == nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeModuleNotExists,
				Msg:  server.ResponseMsgDataNotExists + " task",
			})
		}
		var userID uint64
		userID, err = helper.GetCurrentUserID(ctx)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeParamNotEnough,
				Msg:  server.ResponseMsgParamNotEnough + " userID In Session",
			})
		}
		oldTaskMember, err := service.ProjectTaskMemberService.Instance(&model.ProjectTaskMember{
			TaskID: condition.ID,
			UserID: userID,
		})
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		if oldTaskMember == nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeModuleNotExists,
				Msg:  server.ResponseMsgDataNotExists + " task member",
			})
		}
		// 判断权限
		code := constant.ResourceProjectTaskConfirm
		switch status {
		case model.ProjectTaskStatusConfirmed:
			code = constant.ResourceProjectTaskConfirm
		case model.ProjectTaskStatusDoing:
			code = constant.ResourceProjectTaskStart
		case model.ProjectTaskStatusDone:
			code = constant.ResourceProjectTaskDone
		}
		var success bool
		if _, success, err = helper.CheckResourceCodeInProject(ctx, condition.ProjectID, code); err != nil {
			return err
		} else if !success {
			return nil
		}

		if success, err = service.ProjectTaskMemberService.UpdateStatus(oldTask, oldTaskMember, status, condition.EstimateDuration, condition.ActualDuration); err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		if success {
			c.addSearchDocument(condition.ID)
		}
		return ctx.JSON(&server.CommonResponse{Data: success})
	}
}

func (c *projectTaskController) UpdateStatus(status int) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

		// 先取出原始数据
		oldTask, err := service.ProjectTaskService.Instance(condition.ID)
		if err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		}
		if oldTask == nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeModuleNotExists,
				Msg:  server.ResponseMsgDataNotExists + " task",
			})
		}
		// 判断权限
		code := constant.ResourceProjectTaskCancel
		switch status {
		case model.ProjectTaskStatusCancel:
			code = constant.ResourceProjectTaskCancel
		case model.ProjectTaskStatusNotStart:
			code = constant.ResourceProjectTaskRestart
		}
		var success bool
		if _, success, err = helper.CheckResourceCodeInProject(ctx, oldTask.ProjectID, code); err != nil {
			return err
		} else if !success {
			return nil
		}

		if success, err = service.ProjectTaskService.UpdateStatus(oldTask, status); err != nil {
			return ctx.JSON(&server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			})
		} else {
			if success {
				c.addSearchDocument(condition.ID)
			}
			return ctx.JSON(&server.CommonResponse{Data: success})
		}
	}
}
