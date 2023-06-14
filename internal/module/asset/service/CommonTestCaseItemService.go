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

var CommonTestCaseItemService = new(commonTestCaseItemService)

type commonTestCaseItemService struct{}

// Add 添加测试用例
func (s *commonTestCaseItemService) Add(instance *model.CommonTestCaseItem) (duplicated bool, success bool, err error) {
	if instance.Content == "" || instance.TestCaseID == 0 {
		err = errors.New("Name and testCaseId is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.CommonTestCaseItem{}).Where(&model.CommonTestCaseItem{
		TestCaseID: instance.TestCaseID,
		Content:    instance.Content,
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

// Delete 删除测试用例
func (s *commonTestCaseItemService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.CommonTestCaseItem{ID: id}).Delete(&model.CommonTestCaseItem{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新测试用例项
func (s *commonTestCaseItemService) Update(instance *model.CommonTestCaseItem) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.CommonTestCaseItem{ID: instance.ID}).Updates(&model.CommonTestCaseItem{
		Content:   instance.Content,
		Remark:    instance.Remark,
		CreatorID: instance.CreatorID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *commonTestCaseItemService) PaginateBetweenTimes(condition *model.CommonTestCaseItem, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.CommonTestCaseItem, err error) {
	tx := database.DB.Model(&model.CommonTestCaseItem{})
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
		if condition.Content != "" {
			tx = tx.Where("content like ?", "%"+condition.Content+"%")
		}
	}

	err = tx.Find(&list, &model.CommonTestCaseItem{
		TestCaseID: condition.TestCaseID,
		CreatorID:  condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取测试用例实例
func (s *commonTestCaseItemService) Instance(id uint64) (instance *model.CommonTestCaseItem, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.CommonTestCaseItem{}
	if err = database.DB.Where(&model.CommonTestCaseItem{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *commonTestCaseItemService) ListAllOnlyShow(m *model.CommonTestCaseItem) (list []*model.CommonTestCaseItem, err error) {
	err = database.DB.Select("id,content").Where(&model.CommonTestCaseItem{
		TestCaseID: m.TestCaseID,
	}).Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
