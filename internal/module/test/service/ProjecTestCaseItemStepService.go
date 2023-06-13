package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ProjectTestCaseItemStepService = new(projectTestCaseItemStepService)

type projectTestCaseItemStepService struct{}

// Add 添加资源
func (s *projectTestCaseItemStepService) Add(instance *model.ProjectTestCaseItemStep) (duplicated bool, success bool, err error) {
	if instance.CaseItemID == 0 || instance.Content == "" {
		err = errors.New("Content and case Id is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTestCaseItemStep{}).Where(&model.ProjectTestCaseItemStep{
		CaseItemID: instance.CaseItemID,
		OrderNum:   instance.OrderNum,
	}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectTestCaseItemStepService) Update(instance *model.ProjectTestCaseItemStep) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTestCaseItemStep{ID: instance.ID}).Updates(&model.ProjectTestCaseItemStep{
		CaseItemID: instance.CaseItemID,
		OrderNum:   instance.OrderNum,
		Content:    instance.Content,
		Expect:     instance.Expect,
		Status:     instance.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectTestCaseItemStepService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectTestCaseItemStep{ID: id}).Delete(&model.ProjectTestCaseItemStep{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTestCaseItemStepService) PaginateBetweenTimes(condition *model.ProjectTestCaseItemStep, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTestCaseItemStep, err error) {
	tx := database.DB.Model(&model.ProjectTestCaseItemStep{})
	if limit > -1 {
		tx = tx.Limit(limit)
	}
	if offset > -1 {
		tx = tx.Offset(offset)
	}
	if orderBy != "" {
		tx = tx.Order(orderBy)
	}

	// 处理时间字段，在某段时间之间
	for tc, tr := range tcList {
		if tc != "" {
			if !tr.Start.IsZero() && !tr.End.IsZero() {
				tx.Where(tc+" between ? and ?", time.Time(tr.Start).UnixMilli(), time.Time(tr.End).UnixMilli())
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx.Where(tc+" <= ?", time.Time(tr.End).UnixMilli())
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx.Where(tc+" > ?", time.Time(tr.Start).UnixMilli())
			}
		}
	}

	if condition != nil {
		//if condition.Name != "" {
		//	tx = tx.Where("name like ?", "%"+condition.Name+"%")
		//}
	}

	err = tx.Find(&list, &model.ProjectTestCaseItemStep{
		CaseItemID: condition.CaseItemID,
		Status:     condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectTestCaseItemStepService) Instance(id uint64) (instance *model.ProjectTestCaseItemStep, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectTestCaseItemStep{}
	if err = database.DB.Where(&model.ProjectTestCaseItemStep{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
