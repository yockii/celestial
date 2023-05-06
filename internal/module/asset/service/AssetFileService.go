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

var AssetFileService = new(assetFileService)

type assetFileService struct{}

// Add 添加资源
func (s *assetFileService) Add(instance *model.AssetFile) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.CategoryID == 0 || instance.OssConfigID == 0 {
		err = errors.New("Name and CategoryID and oss config is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.AssetFile{}).Where(&model.AssetFile{
		CategoryID: instance.CategoryID,
		Name:       instance.Name,
		Suffix:     instance.Suffix,
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
func (s *assetFileService) Update(instance *model.AssetFile) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.AssetFile{ID: instance.ID}).Updates(&model.AssetFile{
		CategoryID: instance.CategoryID,
		Name:       instance.Name,
		Suffix:     instance.Suffix,
		CreatorID:  instance.CreatorID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *assetFileService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.AssetFile{ID: id}).Delete(&model.AssetFile{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *assetFileService) PaginateBetweenTimes(condition *model.AssetFile, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.AssetFile, err error) {
	tx := database.DB.Model(&model.AssetFile{}).Limit(100)
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
		if condition.Suffix != "" {
			tx = tx.Where("suffix like ?", "%"+condition.Suffix+"%")
		}
	}

	err = tx.Find(&list, &model.AssetFile{
		CategoryID: condition.CategoryID,
		CreatorID:  condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *assetFileService) Instance(id uint64) (instance *model.AssetFile, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.AssetFile{}
	if err = database.DB.Where(&model.AssetFile{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
