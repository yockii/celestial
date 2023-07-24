package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/meeting/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
)

var MeetingRoomReservationService = new(meetingRoomReservationService)

type meetingRoomReservationService struct{}

func (s *meetingRoomReservationService) List(condition *model.MeetingRoomReservation) (list []*model.MeetingRoomReservation, err error) {
	tx := database.DB.Model(&model.MeetingRoomReservation{})
	tx = tx.Order("start_time asc")
	if condition.MeetingRoomID != 0 {
		tx = tx.Where("meeting_room_id = ?", condition.MeetingRoomID)
	}
	if condition.StartTime != 0 {
		tx = tx.Where("end_time >= ?", condition.StartTime)
	}
	err = tx.Find(&list).Error
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *meetingRoomReservationService) Reserve(reservation *model.MeetingRoomReservation) (inUse bool, instance *model.MeetingRoomReservation, err error) {
	// 检查会议室这个时间段是否被占用
	var c int64
	err = database.DB.Model(&model.MeetingRoomReservation{}).Where(&model.MeetingRoomReservation{MeetingRoomID: reservation.MeetingRoomID}).
		Where(
			database.DB.Where("start_time > ? and start_time < ?", reservation.StartTime, reservation.EndTime).
				Or("end_time > ? and end_time < ?", reservation.StartTime, reservation.EndTime).
				Or("start_time <= ? and end_time >= ?", reservation.StartTime, reservation.EndTime).
				Or("start_time >= ? and end_time <= ?", reservation.StartTime, reservation.EndTime),
		).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		inUse = true
		return
	}
	reservation.ID = util.SnowflakeId()
	err = database.DB.Create(reservation).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	instance = reservation
	return
}

func (s *meetingRoomReservationService) UpdateReservation(old, reservation *model.MeetingRoomReservation) (inUse bool, err error) {
	// 检查会议室这个时间段是否被占用
	var c int64
	err = database.DB.Model(&model.MeetingRoomReservation{}).Where(&model.MeetingRoomReservation{MeetingRoomID: reservation.MeetingRoomID}).
		Where("id != ?", old.ID).
		Where(
			database.DB.Where("start_time > ? and start_time < ?", reservation.StartTime, reservation.EndTime).
				Or("end_time > ? and end_time < ?", reservation.StartTime, reservation.EndTime).
				Or("start_time <= ? and end_time >= ?", reservation.StartTime, reservation.EndTime).
				Or("start_time >= ? and end_time <= ?", reservation.StartTime, reservation.EndTime),
		).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		inUse = true
		return
	}

	err = database.DB.Model(&model.MeetingRoomReservation{}).Where(&model.MeetingRoomReservation{ID: old.ID}).Updates(reservation).Error
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *meetingRoomReservationService) Instance(id uint64) (*model.MeetingRoomReservation, error) {
	instance := new(model.MeetingRoomReservation)
	err := database.DB.Where(&model.MeetingRoomReservation{ID: id}).First(instance).Error
	if err != nil {
		logger.Errorln(err)
	}
	return instance, err
}

func (s *meetingRoomReservationService) DeleteReservation(reservation *model.MeetingRoomReservation) error {
	err := database.DB.Delete(reservation).Error
	if err != nil {
		logger.Errorln(err)
	}
	return err
}
