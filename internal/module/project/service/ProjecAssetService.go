package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ProjectAssetService = new(projectAssetService)

type projectAssetService struct{}

// Add 添加资源
func (s *projectAssetService) Add(instance *model.ProjectAsset) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.ProjectID == 0 {
		err = errors.New("assetName and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectAsset{}).Where(&model.ProjectAsset{
		ProjectID: instance.ProjectID,
		Name:      instance.Name,
		Version:   instance.Version,
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
func (s *projectAssetService) Update(instance *model.ProjectAsset) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectAsset{ID: instance.ID}).Updates(&model.ProjectAsset{
		Name:         instance.Name,
		Type:         instance.Type,
		Version:      instance.Version,
		FileID:       instance.FileID,
		Remark:       instance.Remark,
		Status:       instance.Status,
		VerifyUserID: instance.VerifyUserID,
		VerifyTime:   instance.VerifyTime,
		ReleaseTime:  instance.ReleaseTime,
		ArchiveTime:  instance.ArchiveTime,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectAssetService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectAsset{ID: id}).Delete(&model.ProjectAsset{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectAssetService) PaginateBetweenTimes(condition *model.ProjectAsset, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectAsset, err error) {
	tx := database.DB.Model(&model.ProjectAsset{}).Limit(100)
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

	err = tx.Find(&list, &model.ProjectAsset{
		Type:         condition.Type,
		ProjectID:    condition.ProjectID,
		CreatorID:    condition.CreatorID,
		VerifyUserID: condition.VerifyUserID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectAssetService) Instance(id uint64) (instance *model.ProjectAsset, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectAsset{}
	if err = database.DB.Where(&model.ProjectAsset{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
