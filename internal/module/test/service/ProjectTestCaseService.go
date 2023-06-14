package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/test/domain"
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectTestCaseService = new(projectTestCaseService)

type projectTestCaseService struct{}

// Add 添加资源
func (s *projectTestCaseService) Add(instance *model.ProjectTestCase) (duplicated bool, success bool, err error) {
	if instance.ProjectID == 0 || instance.Name == "" {
		err = errors.New("test Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTestCase{}).Where(&model.ProjectTestCase{
		ProjectID: instance.ProjectID,
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
		ProjectID:   instance.ProjectID,
		RelatedID:   instance.RelatedID,
		RelatedType: instance.RelatedType,
		Name:        instance.Name,
		Remark:      instance.Remark,
		CreatorID:   instance.CreatorID,
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

	// 事务处理，删除测试用例的同时删除测试用例下的测试项及测试项下的测试步骤
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除测试步骤
		err = tx.Where(&model.ProjectTestCaseItemStep{TestCaseID: id}).Delete(&model.ProjectTestCaseItemStep{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除测试用例下的测试项
		err = tx.Where(&model.ProjectTestCaseItem{TestCaseID: id}).Delete(&model.ProjectTestCaseItem{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除测试用例
		err = tx.Where(&model.ProjectTestCase{ID: id}).Delete(&model.ProjectTestCase{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		return nil
	})
	if err != nil {
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
		ProjectID:   condition.ProjectID,
		RelatedID:   condition.RelatedID,
		RelatedType: condition.RelatedType,
		CreatorID:   condition.CreatorID,
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

func (s *projectTestCaseService) BatchAdd(list []*domain.ProjectTestCaseWithItems, creatorID uint64) (bool, error) {
	if len(list) == 0 {
		return false, nil
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, v := range list {
			if v.ProjectID == 0 || v.Name == "" {
				return errors.New("test Name and projectId is required")
			}
			instance := new(model.ProjectTestCase)
			err := tx.Model(&model.ProjectTestCase{}).Where(&model.ProjectTestCase{
				ProjectID: v.ProjectID,
				Name:      v.Name,
			}).First(instance).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					v.ID = util.SnowflakeId()
					v.CreatorID = creatorID
					instance = &v.ProjectTestCase
					if err = tx.Create(instance).Error; err != nil {
						logger.Errorln(err)
						return err
					}
				} else {
					logger.Errorln(err)
					continue
				}
			}

			for _, item := range v.Items {
				// 去重
				if item.Name == "" {
					continue
				}
				var c int64
				err = tx.Model(&model.ProjectTestCaseItem{}).Where(&model.ProjectTestCaseItem{
					TestCaseID: instance.ID,
					Name:       item.Name,
				}).Count(&c).Error
				if err != nil {
					logger.Errorln(err)
					continue
				}

				item.ID = util.SnowflakeId()
				item.TestCaseID = instance.ID
				item.ProjectID = instance.ProjectID
				item.Type = 1
				item.Status = 1
				if err = tx.Create(item).Error; err != nil {
					logger.Errorln(err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
