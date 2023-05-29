package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var CommonTestCaseService = new(commonTestCaseService)

type commonTestCaseService struct{}

// Add 添加测试用例
func (s *commonTestCaseService) Add(instance *model.CommonTestCase) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.CategoryID == 0 {
		err = errors.New("Name and categoryID is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.CommonTestCase{}).Where(&model.CommonTestCase{
		CategoryID: instance.CategoryID,
		Name:       instance.Name,
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

// Update 更新测试用例基本信息
func (s *commonTestCaseService) Update(instance *model.CommonTestCase) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.CommonTestCase{ID: instance.ID}).Updates(&model.CommonTestCase{
		CategoryID: instance.CategoryID,
		Name:       instance.Name,
		CreatorID:  instance.CreatorID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除测试用例
func (s *commonTestCaseService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.CommonTestCase{ID: id}).Delete(&model.CommonTestCase{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *commonTestCaseService) PaginateBetweenTimes(condition *model.CommonTestCase, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.CommonTestCase, err error) {
	tx := database.DB.Model(&model.CommonTestCase{}).Limit(100)
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

	err = tx.Find(&list, &model.CommonTestCase{
		CategoryID: condition.CategoryID,
		CreatorID:  condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取测试用例实例
func (s *commonTestCaseService) Instance(id uint64) (instance *model.CommonTestCase, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.CommonTestCase{}
	if err = database.DB.Where(&model.CommonTestCase{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
