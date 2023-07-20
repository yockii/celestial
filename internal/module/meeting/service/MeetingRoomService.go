package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/meeting/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var MeetingRoomService = new(meetingRoomService)

type meetingRoomService struct{}

// Add 添加资源
func (s *meetingRoomService) Add(instance *model.MeetingRoom) (duplicated bool, success bool, err error) {
	if instance.Name == "" {
		err = errors.New("Name is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.MeetingRoom{}).Where(&model.MeetingRoom{
		Name: instance.Name,
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
	instance.Status = model.MeetingRoomStatusEnabled

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *meetingRoomService) Update(instance *model.MeetingRoom) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.MeetingRoom{ID: instance.ID}).Updates(&model.MeetingRoom{
		Name:     instance.Name,
		Position: instance.Position,
		Capacity: instance.Capacity,
		Devices:  instance.Devices,
		Status:   instance.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *meetingRoomService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.MeetingRoom{ID: id}).Delete(&model.MeetingRoom{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *meetingRoomService) PaginateBetweenTimes(condition *model.MeetingRoom, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.MeetingRoom, err error) {
	tx := database.DB.Model(&model.MeetingRoom{})
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
		if condition.Position != "" {
			tx = tx.Where("position like ?", "%"+condition.Position+"%")
		}
	}

	err = tx.Find(&list, &model.MeetingRoom{
		Status: condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *meetingRoomService) Instance(id uint64) (instance *model.MeetingRoom, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.MeetingRoom{}
	if err = database.DB.Where(&model.MeetingRoom{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
