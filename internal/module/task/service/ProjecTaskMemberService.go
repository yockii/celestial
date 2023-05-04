package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ProjectTaskMemberService = new(projectTaskMemberService)

type projectTaskMemberService struct{}

// Add 添加资源
func (s *projectTaskMemberService) Add(instance *model.ProjectTaskMember) (duplicated bool, success bool, err error) {
	if instance.TaskID == 0 || instance.ProjectID == 0 || instance.UserID == 0 {
		err = errors.New("taskId / projectId / userId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTaskMember{}).Where(&model.ProjectTaskMember{
		ProjectID: instance.ProjectID,
		TaskID:    instance.TaskID,
		UserID:    instance.UserID,
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
func (s *projectTaskMemberService) Update(instance *model.ProjectTaskMember) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTaskMember{ID: instance.ID}).Updates(&model.ProjectTaskMember{
		ProjectID:        instance.ProjectID,
		TaskID:           instance.TaskID,
		UserID:           instance.UserID,
		RoleID:           instance.RoleID,
		EstimateDuration: instance.EstimateDuration,
		ActualDuration:   instance.ActualDuration,
		Status:           instance.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectTaskMemberService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectTaskMember{ID: id}).Delete(&model.ProjectTaskMember{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTaskMemberService) PaginateBetweenTimes(condition *model.ProjectTaskMember, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTaskMember, err error) {
	tx := database.DB.Model(&model.ProjectTaskMember{}).Limit(100)
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
		//if condition.Name != "" {
		//	tx = tx.Where("name like ?", "%"+condition.Name+"%")
		//}
	}

	err = tx.Find(&list, &model.ProjectTaskMember{
		ProjectID: condition.ProjectID,
		TaskID:    condition.TaskID,
		UserID:    condition.UserID,
		RoleID:    condition.RoleID,
		Status:    condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectTaskMemberService) Instance(id uint64) (instance *model.ProjectTaskMember, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectTaskMember{}
	if err = database.DB.Where(&model.ProjectTaskMember{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
