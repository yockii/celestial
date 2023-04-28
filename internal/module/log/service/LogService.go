package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/log/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var LogService = new(stageService)

type stageService struct{}

// Add 添加资源
func (s *stageService) Add(instance *model.Log) (duplicated bool, success bool, err error) {
	if instance.TargetID == 0 || instance.TargetType == 0 {
		err = errors.New("Name and order is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.Log{}).Where("target_id = ? or target_type = ?", instance.TargetID, instance.TargetType).Count(&c).Error
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

// PaginateBetweenTimes 带时间范围的分页查询
func (s *stageService) PaginateBetweenTimes(condition *model.Log, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Log, err error) {
	tx := database.DB.Model(&model.Log{}).Limit(100)
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
		if condition.Action != "" {
			tx = tx.Where("action like ?", "%"+condition.Action+"%")
		}
	}

	err = tx.Find(&list, &model.Log{
		TargetID:   condition.TargetID,
		TargetType: condition.TargetType,
		Status:     condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *stageService) Instance(id uint64) (instance *model.Log, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.Log{}
	if err = database.DB.Where(&model.Log{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
