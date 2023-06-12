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

var ProjectTestCaseService = new(projectTestCaseService)

type projectTestCaseService struct{}

// Add 添加资源
func (s *projectTestCaseService) Add(instance *model.ProjectTestCase) (duplicated bool, success bool, err error) {
	if instance.TestID == 0 || instance.Name == "" {
		err = errors.New("Name and test Id is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTestCase{}).Where(&model.ProjectTestCase{
		ProjectID: instance.ProjectID,
		TestID:    instance.TestID,
		Name:      instance.Name,
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
func (s *projectTestCaseService) Update(instance *model.ProjectTestCase) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTestCase{ID: instance.ID}).Updates(&model.ProjectTestCase{
		ProjectID: instance.ProjectID,
		TestID:    instance.TestID,
		Name:      instance.Name,
		Type:      instance.Type,
		Content:   instance.Content,
		Status:    instance.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectTestCaseService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectTestCase{ID: id}).Delete(&model.ProjectTestCase{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTestCaseService) PaginateBetweenTimes(condition *model.ProjectTestCase, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTestCase, err error) {
	tx := database.DB.Model(&model.ProjectTestCase{})
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
		if condition.Name != "" {
			tx = tx.Where("name like ?", "%"+condition.Name+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectTestCase{
		ProjectID: condition.ProjectID,
		TestID:    condition.TestID,
		Type:      condition.Type,
		Status:    condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectTestCaseService) Instance(id uint64) (instance *model.ProjectTestCase, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectTestCase{}
	if err = database.DB.Where(&model.ProjectTestCase{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
