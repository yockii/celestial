package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/model"
	taskModel "github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectIssueService = new(projectIssueService)

type projectIssueService struct{}

// Add 添加资源
func (s *projectIssueService) Add(instance *model.ProjectIssue) (duplicated bool, success bool, err error) {
	if instance.Title == "" || instance.Type == 0 {
		err = errors.New("assetName and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectIssue{}).Where(&model.ProjectIssue{
		ProjectID: instance.ProjectID,
		Title:     instance.Title,
		Type:      instance.Type,
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
	instance.Status = 1

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 如果存在关联任务，则项目任务及任务人员设置为打回
		if instance.TaskID > 0 {
			err = tx.Model(&taskModel.ProjectTask{}).Where(&taskModel.ProjectTask{ID: instance.TaskID}).Updates(&taskModel.ProjectTask{
				Status: taskModel.ProjectTaskStatusTestReject,
			}).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
			err = tx.Model(&taskModel.ProjectTaskMember{}).Where(&taskModel.ProjectTaskMember{
				TaskID: instance.TaskID,
				UserID: instance.CreatorID,
			}).Updates(&taskModel.ProjectTaskMember{
				Status: taskModel.ProjectTaskStatusTestReject,
			}).Error
			if err != nil {
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
func (s *projectIssueService) Update(instance *model.ProjectIssue) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectIssue{ID: instance.ID}).Updates(&model.ProjectIssue{
		ProjectID:     instance.ProjectID,
		Title:         instance.Title,
		Content:       instance.Content,
		Type:          instance.Type,
		Status:        instance.Status,
		AssigneeID:    instance.AssigneeID,
		StartTime:     instance.StartTime,
		EndTime:       instance.EndTime,
		ResolvedTime:  instance.ResolvedTime,
		SolveDuration: instance.SolveDuration,
		IssueCause:    instance.IssueCause,
		SolveMethod:   instance.SolveMethod,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectIssueService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectIssue{ID: id}).Delete(&model.ProjectIssue{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectIssueService) PaginateBetweenTimes(condition *model.ProjectIssue, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectIssue, err error) {
	tx := database.DB.Model(&model.ProjectIssue{})
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
		if condition.Title != "" {
			tx = tx.Where("title like ?", "%"+condition.Title+"%")
		}
	}

	// 排除大字段
	tx = tx.Omit("content", "issue_cause", "solve_method")

	err = tx.Find(&list, &model.ProjectIssue{
		ID:         condition.ID,
		Type:       condition.Type,
		ProjectID:  condition.ProjectID,
		AssigneeID: condition.AssigneeID,
		Status:     condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectIssueService) Instance(id uint64) (instance *model.ProjectIssue, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectIssue{}
	if err = database.DB.Where(&model.ProjectIssue{ID: id}).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectIssueService) Assign(instance *model.ProjectIssue, assigneeID uint64) (success bool, err error) {
	if instance == nil || instance.ID == 0 || assigneeID == 0 {
		err = errors.New("id is required")
		return
	}
	// 原有数据的状态
	status := instance.Status

	var changedStatus int8 = 0
	if status == 0 || status == 1 {
		changedStatus = 2
	}

	err = database.DB.Where(&model.ProjectIssue{ID: instance.ID}).Updates(&model.ProjectIssue{
		AssigneeID: assigneeID,
		Status:     changedStatus,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// UpdateStatus 更新状态
func (s *projectIssueService) UpdateStatus(instance *model.ProjectIssue, status int8) (success bool, err error) {
	if instance == nil || instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	// 原有数据的状态
	oldStatus := instance.Status
	// 验证状态是否允许更改
	var canChange bool
	switch oldStatus {
	case 0:
		fallthrough
	case model.ProjectIssueStatusNew: // 新建的只能关闭或者指派
		canChange = status == model.ProjectIssueStatusClosed || status == model.ProjectIssueStatusAssigned
	case model.ProjectIssueStatusAssigned: // 已指派的只能关闭或者处理中或者拒绝
		canChange = status == model.ProjectIssueStatusClosed || status == model.ProjectIssueStatusProcessing || status == model.ProjectIssueStatusReject
	case model.ProjectIssueStatusProcessing: // 处理中的只能关闭或者验证中
		canChange = status == model.ProjectIssueStatusClosed || status == model.ProjectIssueStatusVerifying
	case model.ProjectIssueStatusVerifying: // 验证中的只能关闭或者已解决或已指派
		canChange = status == model.ProjectIssueStatusClosed || status == model.ProjectIssueStatusResolved || status == model.ProjectIssueStatusAssigned
	case model.ProjectIssueStatusResolved: // 已解决的只能关闭
		canChange = status == model.ProjectIssueStatusClosed
	case model.ProjectIssueStatusReject: // 被拒绝的状态只能关闭或者重新打开
		canChange = status == model.ProjectIssueStatusClosed || status == model.ProjectIssueStatusNew
	}

	if canChange {
		tx := database.DB.Model(&model.ProjectIssue{}).Where(&model.ProjectIssue{ID: instance.ID})
		if status == model.ProjectIssueStatusReject {
			err = tx.Updates(&model.ProjectIssue{
				Status:         status,
				RejectedReason: instance.RejectedReason,
			}).Error
		} else if status == model.ProjectIssueStatusVerifying {
			err = tx.Updates(map[string]interface{}{
				"status":         status,
				"end_time":       time.Now().UnixMilli(),
				"solve_duration": gorm.Expr("solve_duration + ?", instance.SolveDuration),
			}).Error
		} else if status == model.ProjectIssueStatusProcessing {
			err = tx.Updates(&model.ProjectIssue{
				Status: status,
				// 开始解决时间
				StartTime: time.Now().UnixMilli(),
			}).Error
		} else if status == model.ProjectIssueStatusResolved {
			err = tx.Updates(&model.ProjectIssue{
				Status: status,
				// 解决确认时间
				ResolvedTime: time.Now().UnixMilli(),
			}).Error
		} else {
			err = tx.Updates(&model.ProjectIssue{
				Status: status,
			}).Error
		}
		if err != nil {
			logger.Errorln(err)
			return
		}
		success = true
	}
	return
}
