package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"strings"
	"time"
)

var OssConfigService = new(ossConfigService)

type ossConfigService struct{}

// Add 添加资源
func (s *ossConfigService) Add(instance *model.OssConfig) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.Type == "" || instance.Endpoint == "" || instance.AccessKeyID == "" || instance.SecretAccessKey == "" || instance.Bucket == "" {
		err = errors.New("Name/Type/Endpoint/AccessKeyID/SecretAccessKey/BucketName is required ")
		return
	}
	instance.Type = strings.ToLower(instance.Type)
	var c int64
	err = database.DB.Model(&model.OssConfig{}).Where(&model.OssConfig{
		Type:     instance.Type,
		Name:     instance.Name,
		Endpoint: instance.Endpoint,
		Bucket:   instance.Bucket,
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
func (s *ossConfigService) Update(instance *model.OssConfig) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.OssConfig{ID: instance.ID}).Updates(&model.OssConfig{
		Type:            strings.ToLower(instance.Type),
		Name:            instance.Name,
		Endpoint:        instance.Endpoint,
		AccessKeyID:     instance.AccessKeyID,
		SecretAccessKey: instance.SecretAccessKey,
		Bucket:          instance.Bucket,
		Region:          instance.Region,
		Secure:          instance.Secure,
		SelfDomain:      instance.SelfDomain,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *ossConfigService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.OssConfig{ID: id}).Delete(&model.OssConfig{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *ossConfigService) PaginateBetweenTimes(condition *model.OssConfig, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.OssConfig, err error) {
	tx := database.DB.Model(&model.OssConfig{})
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
		if condition.Type != "" {
			tx = tx.Where("type = ?", strings.ToLower(condition.Type))
		}
		if condition.Bucket != "" {
			tx = tx.Where("bucket like ?", "%"+condition.Bucket+"%")
		}
	}

	err = tx.Find(&list, &model.OssConfig{
		Secure:     condition.Secure,
		SelfDomain: condition.SelfDomain,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *ossConfigService) Instance(id uint64) (instance *model.OssConfig, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.OssConfig{}
	if err = database.DB.Where(&model.OssConfig{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// UpdateStatus 更新状态
func (s *ossConfigService) UpdateStatus(id uint64, status int) (success bool, err error) {
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 如果status=1则先把其他的都设置为-1
		if status == 1 {
			if err = tx.Model(&model.OssConfig{}).Where("status = ?", 1).Update("status", -1).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}
		if err = tx.Model(&model.OssConfig{ID: id}).Update("status", status).Error; err != nil {
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
