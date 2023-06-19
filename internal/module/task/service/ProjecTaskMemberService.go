package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/task/domain"
	"github.com/yockii/celestial/internal/module/task/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
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

// ListWithRealName 查询列表，并附带用户真实姓名
func (s *projectTaskMemberService) ListWithRealName(condition *model.ProjectTaskMember) (list []*domain.ProjectTaskMemberWithRealName, err error) {
	tx := database.DB.Model(&model.ProjectTaskMember{})

	if condition != nil {
		//if condition.Name != "" {
		//	tx = tx.Where("name like ?", "%"+condition.Name+"%")
		//}
	}

	sm := gorm.Statement{DB: database.DB}
	_ = sm.Parse(&ucModel.User{})
	userTableName := sm.Schema.Table
	_ = sm.Parse(&model.ProjectTaskMember{})
	ptmTableName := sm.Schema.Table

	tx.Select(ptmTableName+".*", "real_name")

	err = tx.Joins("left join "+userTableName+" on "+ptmTableName+".user_id = "+userTableName+".id").Find(&list, &model.ProjectTaskMember{
		ProjectID: condition.ProjectID,
		TaskID:    condition.TaskID,
		UserID:    condition.UserID,
		RoleID:    condition.RoleID,
		Status:    condition.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// List 查询列表
func (s *projectTaskMemberService) List(condition *model.ProjectTaskMember) (list []*model.ProjectTaskMember, err error) {
	tx := database.DB.Model(&model.ProjectTaskMember{})

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
	}).Error
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

// UpdateStatus 更新状态
func (s *projectTaskMemberService) UpdateStatus(taskID, userID uint64, status int) (success bool, err error) {
	if taskID == 0 || userID == 0 {
		err = errors.New("id is required")
		return
	}

	// 检查任务是否已取消状态，该状态不允许做任何状态变更
	var task = new(model.ProjectTask)
	err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{
		ID: taskID,
	}).First(&task).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln(err)
			return
		}
		return false, nil
	}
	if task.Status == model.ProjectTaskStatusCancel {
		return false, nil
	}

	// 检查旧状态
	var oldTM = new(model.ProjectTaskMember)
	err = database.DB.Model(&model.ProjectTaskMember{}).Where(&model.ProjectTaskMember{
		TaskID: taskID,
		UserID: userID,
	}).First(&oldTM).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln(err)
			return
		}
		return false, nil
	}
	if oldTM.Status == status {
		return true, nil
	}

	// 判断当前状态是否可变更为目标状态
	var canChange bool
	switch oldTM.Status {
	case model.ProjectTaskStatusCancel: // 已取消，可以变更为未开始
		canChange = status == model.ProjectTaskStatusNotStart
	case 0: // 未知状态，可以变更为未开始
		fallthrough
	case model.ProjectTaskStatusNotStart: // 未开始，可以变更为已取消、已确认
		canChange = status == model.ProjectTaskStatusCancel || status == model.ProjectTaskStatusConfirmed
	case model.ProjectTaskStatusConfirmed: // 已确认，可以变更为进行中、已取消
		canChange = status == model.ProjectTaskStatusCancel || status == model.ProjectTaskStatusDoing
	case model.ProjectTaskStatusDoing: // 进行中，可以变更为已完成、已取消
		canChange = status == model.ProjectTaskStatusCancel || status == model.ProjectTaskStatusDone
	case model.ProjectTaskStatusDone: // 已完成，可以变更为已取消
		canChange = status == model.ProjectTaskStatusCancel
	}

	if !canChange {
		return false, nil
	}

	// 开启事务
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.ProjectTaskMember{}).Where(&model.ProjectTaskMember{
			ID: oldTM.ID,
		}).Update("status", status).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 任务成员状态变更后，检查是否需要更新任务状态
		// 情况1：任务成员全部确认，任务状态更新为已确认；情况2：任务成员开始工作，任务状态更新为进行中；情况3：任务成员全部完成，任务状态更新为已完成；
		if status == model.ProjectTaskStatusConfirmed || status == model.ProjectTaskStatusDone {
			var count int64
			err = tx.Model(&model.ProjectTaskMember{}).Where(&model.ProjectTaskMember{
				TaskID: taskID,
			}).Where("status <> ?", status).Count(&count).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
			if count == 0 && task.Status != status {
				// 都已确认，任务更新为已确认
				err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{
					ID: taskID,
				}).Update("status", status).Error
				if err != nil {
					logger.Errorln(err)
					return err
				}
			}
		} else if status == model.ProjectTaskStatusDoing && task.Status != status {
			// 任务成员开始工作，任务更新为进行中
			err = database.DB.Model(&model.ProjectTask{}).Where(&model.ProjectTask{
				ID: taskID,
			}).Update("status", status).Error
			if err != nil {
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
