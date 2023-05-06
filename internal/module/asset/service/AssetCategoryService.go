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

var AssetCategoryService = new(assetCategoryService)

type assetCategoryService struct{}

// Add 添加资源
func (s *assetCategoryService) Add(instance *model.AssetCategory) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.Type == 0 {
		err = errors.New("Name and type is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{
		ParentID: instance.ParentID,
		Name:     instance.Name,
		Type:     instance.Type,
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
func (s *assetCategoryService) Update(instance *model.AssetCategory) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.AssetCategory{ID: instance.ID}).Updates(&model.AssetCategory{
		ParentID:  instance.ParentID,
		Name:      instance.Name,
		Type:      instance.Type,
		CreatorID: instance.CreatorID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *assetCategoryService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.AssetCategory{ID: id}).Delete(&model.AssetCategory{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *assetCategoryService) PaginateBetweenTimes(condition *model.AssetCategory, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.AssetCategory, err error) {
	tx := database.DB.Model(&model.AssetCategory{}).Limit(100)
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

	err = tx.Find(&list, &model.AssetCategory{
		Type:      condition.Type,
		ParentID:  condition.ParentID,
		CreatorID: condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *assetCategoryService) Instance(id uint64) (instance *model.AssetCategory, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.AssetCategory{}
	if err = database.DB.Where(&model.AssetCategory{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
