package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectTestCaseItemService = new(projectTestCaseItemService)

type projectTestCaseItemService struct{}

// Add 添加资源
func (s *projectTestCaseItemService) Add(instance *model.ProjectTestCaseItem) (duplicated bool, success bool, err error) {
	if instance.TestCaseID == 0 || instance.Name == "" {
		err = errors.New("Name and test Id is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTestCaseItem{}).Where(&model.ProjectTestCaseItem{
		ProjectID:  instance.ProjectID,
		TestCaseID: instance.TestCaseID,
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
	instance.Status = 1

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectTestCaseItemService) Update(instance *model.ProjectTestCaseItem) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTestCaseItem{ID: instance.ID}).Updates(&model.ProjectTestCaseItem{
		ProjectID:  instance.ProjectID,
		TestCaseID: instance.TestCaseID,
		Name:       instance.Name,
		Type:       instance.Type,
		Content:    instance.Content,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// UpdateStatus 更新资源状态
func (s *projectTestCaseItemService) UpdateStatus(id uint64, status int) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Model(&model.ProjectTestCaseItem{}).Where(&model.ProjectTestCaseItem{ID: id}).Update("status", status).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectTestCaseItemService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectTestCaseItem{ID: id}).Delete(&model.ProjectTestCaseItem{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTestCaseItemService) PaginateBetweenTimes(condition *model.ProjectTestCaseItem, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTestCaseItem, err error) {
	tx := database.DB.Model(&model.ProjectTestCaseItem{})
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
				tx = tx.Where(tc+" between ? and ?", time.Time(tr.Start).UnixMilli(), time.Time(tr.End).UnixMilli())
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx = tx.Where(tc+" <= ?", time.Time(tr.End).UnixMilli())
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx = tx.Where(tc+" > ?", time.Time(tr.Start).UnixMilli())
			}
		}
	}

	if condition != nil {
		if condition.Name != "" {
			tx = tx.Where("name like ?", "%"+condition.Name+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectTestCaseItem{
		ProjectID:  condition.ProjectID,
		TestCaseID: condition.TestCaseID,
		Type:       condition.Type,
		Status:     condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectTestCaseItemService) Instance(id uint64) (instance *model.ProjectTestCaseItem, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectTestCaseItem{}
	if err = database.DB.Where(&model.ProjectTestCaseItem{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectTestCaseItemService) BatchAdd(list []*model.ProjectTestCaseItem) (success bool, err error) {
	if len(list) == 0 {
		err = errors.New("list is required")
		return
	}

	// 事务处理
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, v := range list {
			if v.ProjectID == 0 || v.Name == "" || v.TestCaseID == 0 {
				continue
			}
			// 去重
			var c int64
			if err = tx.Model(&model.ProjectTestCaseItem{}).Where(&model.ProjectTestCaseItem{
				ProjectID:  v.ProjectID,
				TestCaseID: v.TestCaseID,
				Name:       v.Name,
			}).Count(&c).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			if c > 0 {
				continue
			}

			v.ID = util.SnowflakeId()
			v.Status = 1
			if err = tx.Create(v).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	success = true
	return
}
