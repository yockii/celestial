package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/task/domain"
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ProjectTaskService = new(projectTaskService)

type projectTaskService struct{}

// Add 添加资源
func (s *projectTaskService) Add(instance *model.ProjectTask) (duplicated bool, success bool, err error) {
	if instance.ProjectID == 0 || instance.TaskName == "" {
		err = errors.New("task Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{
		ProjectID: instance.ProjectID,
		TaskName:  instance.TaskName,
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
func (s *projectTaskService) Update(instance *model.ProjectTask) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectTask{ID: instance.ID}).Updates(&model.ProjectTask{
		ProjectID:        instance.ProjectID,
		StageID:          instance.StageID,
		ParentID:         instance.ParentID,
		TaskName:         instance.TaskName,
		TaskDesc:         instance.TaskDesc,
		StartTime:        instance.StartTime,
		EndTime:          instance.EndTime,
		Priority:         instance.Priority,
		OwnerID:          instance.OwnerID,
		ActualStartTime:  instance.ActualStartTime,
		ActualEndTime:    instance.ActualEndTime,
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
func (s *projectTaskService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectTask{ID: id}).Delete(&model.ProjectTask{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTaskService) PaginateBetweenTimes(condition *model.ProjectTask, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTask, err error) {
	tx := database.DB.Model(&model.ProjectTask{}).Limit(100)
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
		if condition.TaskName != "" {
			tx = tx.Where("task_name like ?", "%"+condition.TaskName+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectTask{
		ProjectID: condition.ProjectID,
		ParentID:  condition.ParentID,
		StageID:   condition.StageID,
		Status:    condition.Status,
		Priority:  condition.Priority,
		OwnerID:   condition.OwnerID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectTaskService) Instance(id uint64) (instance *model.ProjectTask, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectTask{}
	if err = database.DB.Where(&model.ProjectTask{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectTaskService) TaskDurationByProject(projectID uint64, tcList map[string]*server.TimeCondition) (result *domain.ProjectTaskWorkTimeStatistics, err error) {
	result = &domain.ProjectTaskWorkTimeStatistics{
		ProjectID: projectID,
	}
	tx := database.DB.Model(&model.ProjectTask{})
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

	err = tx.Select("count(1) as task_count,sum(actual_duration) as actual_duration, sum(estimate_duration) as estimate_duration").Where(&model.ProjectTask{
		ProjectID: projectID,
	}).Scan(result).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
