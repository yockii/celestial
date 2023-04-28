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

var WorkTimeService = new(workTimeService)

type workTimeService struct{}

// Add 添加资源
func (s *workTimeService) Add(instance *model.WorkTime) (duplicated bool, success bool, err error) {
	if instance.TargetID == 0 || instance.TargetType == 0 {
		err = errors.New("assetName and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.WorkTime{}).Where(&model.WorkTime{
		TargetID:   instance.TargetID,
		TargetType: instance.TargetType,
		UserID:     instance.UserID,
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
func (s *workTimeService) Update(instance *model.WorkTime) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.WorkTime{ID: instance.ID}).Updates(&model.WorkTime{
		TargetID:     instance.TargetID,
		TargetType:   instance.TargetType,
		WorkTime:     instance.WorkTime,
		WorkContent:  instance.WorkContent,
		ReviewerID:   instance.ReviewerID,
		Status:       instance.Status,
		ReviewTime:   instance.ReviewTime,
		RejectReason: instance.RejectReason,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *workTimeService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.WorkTime{ID: id}).Delete(&model.WorkTime{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *workTimeService) PaginateBetweenTimes(condition *model.WorkTime, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.WorkTime, err error) {
	tx := database.DB.Model(&model.WorkTime{}).Limit(100)
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
		//if condition.Title != "" {
		//	tx = tx.Where("title like ?", "%"+condition.Title+"%")
		//}
	}

	err = tx.Find(&list, &model.WorkTime{
		TargetID:   condition.TargetID,
		TargetType: condition.TargetType,
		UserID:     condition.UserID,
		ReviewerID: condition.ReviewerID,
		Status:     condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *workTimeService) Instance(id uint64) (instance *model.WorkTime, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.WorkTime{}
	if err = database.DB.Where(&model.WorkTime{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
