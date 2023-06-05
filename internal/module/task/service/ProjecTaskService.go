package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	projectModel "github.com/yockii/celestial/internal/module/project/model"
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
	if instance.ProjectID == 0 || instance.Name == "" {
		err = errors.New("task Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{
		ProjectID: instance.ProjectID,
		Name:      instance.Name,
	}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	// 获取对应的需求，赋值fullPath
	requirement := &projectModel.ProjectRequirement{ID: instance.RequirementID}
	if err = database.DB.Model(requirement).First(&requirement).Error; err != nil {
		logger.Errorln(err)
		return
	}
	instance.FullPath = requirement.FullPath + "/" + instance.Name

	instance.ID = util.SnowflakeId()
	instance.Status = model.ProjectTaskStatusNotStart

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

	if instance.RequirementID != 0 {
		// 判断需求ID是否发生变更，如果发生变更，则更新fullPath
		var oldInstance model.ProjectTask
		if err = database.DB.Model(&model.ProjectTask{}).First(&oldInstance, instance.ID).Error; err != nil {
			logger.Errorln(err)
			return
		}
		if oldInstance.RequirementID != instance.RequirementID {
			requirement := &projectModel.ProjectRequirement{ID: instance.RequirementID}
			if err = database.DB.Model(requirement).First(&requirement).Error; err != nil {
				logger.Errorln(err)
				return
			}
			instance.FullPath = requirement.FullPath + "/" + instance.Name
		}
	}

	err = database.DB.Where(&model.ProjectTask{ID: instance.ID}).Updates(&model.ProjectTask{
		ProjectID:        instance.ProjectID,
		StageID:          instance.StageID,
		ParentID:         instance.ParentID,
		RequirementID:    instance.RequirementID,
		Name:             instance.Name,
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
func (s *projectTaskService) PaginateBetweenTimes(condition *model.ProjectTask, onlyParent bool, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTask, err error) {
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
		if condition.Name != "" {
			tx = tx.Where("name like ?", "%"+condition.Name+"%")
		}
		if condition.FullPath != "" {
			tx = tx.Where("full_path like ?", condition.FullPath+"%")
		}
	}
	if onlyParent {
		tx = tx.Where("parent_id = 0")
	}

	// 大字段不查询
	tx.Omit("task_desc", "full_path")

	err = tx.Find(&list, &model.ProjectTask{
		ID:        condition.ID,
		ProjectID: condition.ProjectID,
		ParentID:  condition.ParentID,
		ModuleID:  condition.ModuleID,
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
