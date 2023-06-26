package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/panjf2000/ants/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/data"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/celestial/internal/module/project/service"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
)

var ProjectRiskController = new(projectRiskController)

type projectRiskController struct{}

func (c *projectRiskController) Add(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRisk)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}

	// 处理必填
	if instance.RiskName == "" || instance.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name & projectId",
		})
	}

	// 判断权限
	if uid, err := helper.CheckResourceCodeInProject(ctx, instance.ProjectID, constant.ResourceProjectRiskAdd); err != nil {
		return err
	} else {
		instance.CreatorID = uid
	}

	duplicated, success, err := service.ProjectRiskService.Add(instance)
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

	_ = ants.Submit(func(d *model.ProjectRisk) func() {
		p := "低"
		switch d.RiskProbability {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		e := "低"
		switch d.RiskImpact {
		case 1:
			e = "低"
		case 2:
			e = "中"
		case 3:
			e = "高"
		}
		l := "低"
		switch d.RiskLevel {
		case 1:
			l = "低"
		case 2:
			l = "中"
		case 3:
			l = "高"
		}
		status := "已识别"
		//风险状态 1-已识别 2-已应对 3-已发生 4-已解决
		switch d.Status {
		case 1:
			status = "已识别"
		case 2:
			status = "已应对"
		case 3:
			status = "已发生"
		case 4:
			status = "已解决"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.RiskName,
			Content: fmt.Sprintf("%s\n概率：%s, 影响：%s, 等级：%s, 状态：%s\n应对措施：%s\n总结：%s",
				d.RiskDesc, p, e, l, status, d.Response, d.Result,
			),
			Route:      fmt.Sprintf("/project/detail/%d/risk?id=%d", d.ProjectID, d.ID),
			CreateTime: d.CreateTime,
			UpdateTime: d.UpdateTime,
		}, d.CreatorID)
	}(instance))

	return ctx.JSON(&server.CommonResponse{
		Data: instance,
	})
}

func (c *projectRiskController) Update(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRisk)
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
	old, err := service.ProjectRiskService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	if _, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectRiskUpdate); err != nil {
		return err
	}

	success, err := service.ProjectRiskService.Update(instance)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}

	if success {
		c.updateSearchDocument(instance.ID)
	}

	return ctx.JSON(&server.CommonResponse{
		Data: success,
	})
}

func (c *projectRiskController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.ProjectRisk)
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
	old, err := service.ProjectRiskService.Instance(instance.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	// 判断权限
	if _, err = helper.CheckResourceCodeInProject(ctx, old.ProjectID, constant.ResourceProjectRiskDelete); err != nil {
		return err
	}

	success, err := service.ProjectRiskService.Delete(instance.ID)
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

func (c *projectRiskController) List(ctx *fiber.Ctx) error {
	instance := new(domain.ProjectRiskListRisk)
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

	total, list, err := service.ProjectRiskService.PaginateBetweenTimes(&instance.ProjectRisk, paginate.Limit, paginate.Offset, instance.OrderBy, tcList)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Items:  list,
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		},
	})
}

func (c *projectRiskController) Instance(ctx *fiber.Ctx) error {
	condition := new(model.ProjectRisk)
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
	dept, err := service.ProjectRiskService.Instance(condition.ID)
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

// CalculateRiskByProject 计算项目风险
func (c *projectRiskController) CalculateRiskByProject(ctx *fiber.Ctx) error {
	condition := new(model.ProjectRisk)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if condition.ProjectID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " projectId",
		})
	}
	riskScore, maxRiskInfo, err := service.ProjectRiskService.CalculateRiskByProject(condition.ProjectID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: &domain.ProjectRiskCoefficient{
			RiskCoefficient: riskScore,
			MaxRisk:         maxRiskInfo,
		},
	})
}

func (c *projectRiskController) updateSearchDocument(id uint64) {
	_ = ants.Submit(func(id uint64) func() {
		d, er := service.ProjectRiskService.Instance(id)
		if er != nil {
			logger.Errorln(er)
			return func() {}
		}
		p := "低"
		switch d.RiskProbability {
		case 1:
			p = "低"
		case 2:
			p = "中"
		case 3:
			p = "高"
		}
		e := "低"
		switch d.RiskImpact {
		case 1:
			e = "低"
		case 2:
			e = "中"
		case 3:
			e = "高"
		}
		l := "低"
		switch d.RiskLevel {
		case 1:
			l = "低"
		case 2:
			l = "中"
		case 3:
			l = "高"
		}
		status := "已识别"
		//风险状态 1-已识别 2-已应对 3-已发生 4-已解决
		switch d.Status {
		case 1:
			status = "已识别"
		case 2:
			status = "已应对"
		case 3:
			status = "已发生"
		case 4:
			status = "已解决"
		}
		return data.AddDocumentAntsWrapper(&search.Document{
			ID:    d.ID,
			Title: d.RiskName,
			Content: fmt.Sprintf("%s\n概率：%s, 影响：%s, 等级：%s, 状态：%s\n应对措施：%s\n总结：%s",
				d.RiskDesc, p, e, l, status, d.Response, d.Result,
			),
			UpdateTime: d.UpdateTime,
		})
	}(id))
}
