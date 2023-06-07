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
