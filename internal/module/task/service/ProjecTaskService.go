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
	"gorm.io/gorm"
	"time"
)

var ProjectTaskService = new(projectTaskService)

type projectTaskService struct{}

// Add 添加资源
func (s *projectTaskService) Add(instance *model.ProjectTask, members []*domain.ProjectTaskMemberWithRealName) (duplicated bool, success bool, err error) {
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
	// 如果有父级任务，赋值fullPath
	var parent *model.ProjectTask
	if instance.ParentID != 0 {
		parent = &model.ProjectTask{ID: instance.ParentID}
		if err = database.DB.Model(parent).First(&parent).Error; err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = parent.FullPath + "/" + instance.Name
	} else {
		requirement := &projectModel.ProjectRequirement{ID: instance.RequirementID}
		if err = database.DB.Model(requirement).First(&requirement).Error; err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = requirement.FullPath + "/" + instance.Name
	}
	instance.ID = util.SnowflakeId()
	instance.Status = model.ProjectTaskStatusNotStart

	// 事务，同步新增任务和任务成员
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 更新父级任务的子任务数量
		if parent != nil {
			if err = tx.Model(parent).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		for _, m := range members {
			member := &m.ProjectTaskMember
			member.ID = util.SnowflakeId()
			member.TaskID = instance.ID
			member.ProjectID = instance.ProjectID
			member.Status = 1

			if err = tx.Create(member).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}
		return nil
	}); err != nil {
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectTaskService) Update(instance *model.ProjectTask, members []*domain.ProjectTaskMemberWithRealName) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 判断父级任务是否发生变更，如果变更，更新full path，后面事务中还要更新各自父级数量
	var oldInstance model.ProjectTask
	if err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: instance.ID}).First(&oldInstance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	fullPathChanged := false
	parentChanged := false
	if instance.ParentID != 0 {
		if oldInstance.ParentID != instance.ParentID {
			parent := &model.ProjectTask{ID: instance.ParentID}
			if err = database.DB.Model(&model.ProjectTask{}).Where(parent).First(&parent).Error; err != nil {
				logger.Errorln(err)
				return
			}
			instance.FullPath = parent.FullPath + "/" + instance.Name
			fullPathChanged = true
			parentChanged = true
		}
	} else {
		if oldInstance.RequirementID != instance.RequirementID {
			requirement := &projectModel.ProjectRequirement{ID: instance.RequirementID}
			if err = database.DB.Model(requirement).First(&requirement).Error; err != nil {
				logger.Errorln(err)
				return
			}
			instance.FullPath = requirement.FullPath + "/" + instance.Name
			fullPathChanged = true
		}
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 更新
		err = tx.Where(&model.ProjectTask{ID: instance.ID}).Updates(&model.ProjectTask{
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
			FullPath:         instance.FullPath,
		}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 如果parentChanged，更新原有父级任务的子任务数量和现父级任务子任务数量
		if parentChanged {
			if err = tx.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: oldInstance.ParentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			if err = tx.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: instance.ParentID}).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		// 如果fullPathChanged，更新子任务的fullPath
		if fullPathChanged {
			if err = tx.Model(&model.ProjectTask{}).Where("full_path like ?", instance.FullPath+"%").Update("full_path", gorm.Expr("concat(?, substring(full_path, ?))", instance.FullPath, len(oldInstance.FullPath)+1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		// 取出原有members，进行对比，并进行新增和删除
		var oldMembers []*model.ProjectTaskMember
		if err = tx.Model(&model.ProjectTaskMember{}).Where(&model.ProjectTaskMember{TaskID: instance.ID}).Find(&oldMembers).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		// 需要删除的members
		var deleteMembers []*model.ProjectTaskMember
		for _, oldMember := range oldMembers {
			var found bool
			for _, member := range members {
				if oldMember.ID == member.ID {
					found = true
					break
				}
			}
			if !found {
				deleteMembers = append(deleteMembers, oldMember)
			}
		}
		// 需要新增的members
		var addMembers []*model.ProjectTaskMember
		for _, member := range members {
			var found bool
			for _, oldMember := range oldMembers {
				if oldMember.ID == member.ID {
					found = true
					break
				}
			}
			if !found {
				addMembers = append(addMembers, &member.ProjectTaskMember)
			}
		}
		// 删除members
		for _, member := range deleteMembers {
			if err = tx.Delete(member).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}
		// 新增members
		for _, member := range addMembers {
			member.ID = util.SnowflakeId()
			member.TaskID = instance.ID
			member.ProjectID = instance.ProjectID
			if err = tx.Create(member).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		return nil
	})
	if err != nil {
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

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除，如果有父级，需要更新父级的子任务数量
		task := new(model.ProjectTask)
		if err = tx.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: id}).First(&task).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if task.ParentID != 0 {
			if err = tx.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: task.ParentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		err = database.DB.Where(&model.ProjectTask{ID: id}).Delete(&model.ProjectTask{}).Error
		if err != nil {
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

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectTaskService) PaginateBetweenTimes(condition *model.ProjectTask, onlyParent bool, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectTask, err error) {
	tx := database.DB.Model(&model.ProjectTask{})
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
		if condition.FullPath != "" {
			tx = tx.Where("full_path like ?", condition.FullPath+"%")
		}
	}
	if onlyParent {
		tx = tx.Where("parent_id = 0")
	}

	// 大字段不查询
	tx = tx.Omit("task_desc", "full_path")

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
				tx = tx.Where(tc+" between ? and ?", time.Time(tr.Start).UnixMilli(), time.Time(tr.End).UnixMilli())
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx = tx.Where(tc+" <= ?", time.Time(tr.End).UnixMilli())
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx = tx.Where(tc+" > ?", time.Time(tr.Start).UnixMilli())
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

// UpdateStatus 更新状态
func (s *projectTaskService) UpdateStatus(id uint64, status int) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	// 检查旧状态
	var oldStatus int
	err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: id}).Select("status").First(&oldStatus).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln(err)
			return
		}
		return false, nil
	}
	if oldStatus == status {
		return true, nil
	}

	// 判断当前状态是否可变更为目标状态
	var canChange bool
	switch status {
	case model.ProjectTaskStatusCancel: // 已取消，可以变更为未开始
		canChange = oldStatus == model.ProjectTaskStatusNotStart
	case model.ProjectTaskStatusNotStart: // 未开始，可以变更为已取消、已确认
		canChange = oldStatus == model.ProjectTaskStatusCancel || oldStatus == model.ProjectTaskStatusConfirmed
	case model.ProjectTaskStatusConfirmed: // 已确认，可以变更为进行中、已取消
		canChange = oldStatus == model.ProjectTaskStatusCancel || oldStatus == model.ProjectTaskStatusDoing
	case model.ProjectTaskStatusDoing: // 进行中，可以变更为已完成、已取消
		canChange = oldStatus == model.ProjectTaskStatusCancel || oldStatus == model.ProjectTaskStatusDone
	case model.ProjectTaskStatusDone: // 已完成，可以变更为已取消
		canChange = oldStatus == model.ProjectTaskStatusCancel
	}

	if !canChange {
		return false, nil
	}

	err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{ID: id}).Update("status", status).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}
