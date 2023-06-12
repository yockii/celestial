package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"strings"
	"time"
)

var ResourceService = new(resourceService)

type resourceService struct{}

// Add 添加资源
func (s *resourceService) Add(instance *model.Resource) (duplicated bool, success bool, err error) {
	if instance.ResourceName == "" || instance.ResourceCode == "" {
		err = errors.New("resourceName and resourceCode is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.Resource{}).Where(&model.Resource{ResourceName: instance.ResourceName}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()

	err = database.DB.Create(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *resourceService) Update(instance *model.Resource) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.Resource{ID: instance.ID}).Updates(&model.Resource{
		ResourceName: instance.ResourceName,
		ResourceCode: strings.ToLower(instance.ResourceCode),
		Type:         instance.Type,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *resourceService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.Resource{ID: id}).Delete(&model.Resource{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *resourceService) PaginateBetweenTimes(condition *model.Resource, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Resource, err error) {
	tx := database.DB.Model(&model.Resource{})
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
		if condition.ResourceName != "" {
			tx = tx.Where("resource_name like ?", "%"+condition.ResourceName+"%")
		}
		if condition.ResourceCode != "" {
			tx = tx.Where("resource_code like ?", "%"+condition.ResourceCode+"%")
		}
	}

	err = tx.Where(&model.Resource{Type: condition.Type}).Find(&list).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *resourceService) Instance(id uint64) (instance *model.Resource, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.Resource{}
	err = database.DB.Where(&model.Resource{ID: id}).First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
