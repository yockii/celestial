package service

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/celestial/internal/module/asset/provider"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

var AssetFileService = &assetFileService{
	locker: new(sync.Mutex),
}

type assetFileService struct {
	osManager provider.OsManager
	locker    sync.Locker
}

// Add 添加资源
func (s *assetFileService) Add(instance *model.File, reader io.Reader) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.CategoryID == 0 {
		err = errors.New("Name and CategoryID is required ")
		return
	}
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}

	// 检查是否有重复的资源文件
	var c int64
	err = database.DB.Model(&model.File{}).Where(&model.File{
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
	// 上传文件
	now := time.Now()
	objName := fmt.Sprintf("%d%d%d/%d.%s", now.Year(), now.Month(), now.Day(), instance.ID, instance.Suffix)
	if err = s.osManager.PutObject(objName, reader); err != nil {
		return false, false, err
	}

	instance.OssConfigID = s.osManager.GetOssConfigID()
	instance.ObjName = objName
	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *assetFileService) Update(instance *model.File) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.File{ID: instance.ID}).Updates(&model.File{
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
	err = database.DB.Where(&model.File{ID: id}).Delete(&model.File{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *assetFileService) PaginateBetweenTimes(condition *model.File, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.File, err error) {
	tx := database.DB.Model(&model.File{}).Limit(100)
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

	err = tx.Find(&list, &model.File{
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
func (s *assetFileService) Instance(id uint64) (instance *model.File, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.File{}
	if err = database.DB.Where(&model.File{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *assetFileService) initOsManager() (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.osManager == nil {
		// 从数据库中取出可用的云存储配置
		c := &model.OssConfig{}
		if err = database.DB.Where(&model.OssConfig{Status: 1}).First(c).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("no available oss config")
			}
			return
		}
		// 初始化云存储
		osm := provider.GetProvider(c)
		if err = osm.Auth(); err != nil {
			return
		}
		s.osManager = osm
	}
	return nil
}
