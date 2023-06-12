package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectRequirementService = new(projectRequirementService)

type projectRequirementService struct{}

// Add 添加资源
func (s *projectRequirementService) Add(instance *model.ProjectRequirement) (duplicated bool, success bool, err error) {
	if instance.ProjectID == 0 || instance.Name == "" {
		err = errors.New("requirement Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectRequirement{}).Where(&model.ProjectRequirement{
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

	// 如果有对应的模块ID，则获取模块的fullPath
	if instance.ModuleID != 0 {
		module := &model.ProjectModule{ID: instance.ModuleID}
		if err = database.DB.Model(module).First(&module).Error; err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = module.FullPath
	}

	instance.ID = util.SnowflakeId()
	instance.Status = model.ProjectRequirementStatusPendingReview

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectRequirementService) Update(instance *model.ProjectRequirement) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 对比原来的模块ID和新的模块ID是否一致，不一致则需要更新fullPath
	var oldModuleID uint64
	err = database.DB.Model(&model.ProjectRequirement{}).Where(&model.ProjectRequirement{ID: instance.ID}).Select("module_id").First(&oldModuleID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln(err)
			return
		}
	}
	if oldModuleID != instance.ModuleID {
		module := &model.ProjectModule{ID: instance.ModuleID}
		if err = database.DB.Model(module).First(&module).Error; err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = module.FullPath
	}

	err = database.DB.Where(&model.ProjectRequirement{ID: instance.ID}).Updates(&model.ProjectRequirement{
		ProjectID:   instance.ProjectID,
		ModuleID:    instance.ModuleID,
		StageID:     instance.StageID,
		Type:        instance.Type,
		Name:        instance.Name,
		Detail:      instance.Detail,
		Priority:    instance.Priority,
		Source:      instance.Source,
		OwnerID:     instance.OwnerID,
		Feasibility: instance.Feasibility,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}

	success = true
	return
}

// Delete 删除资源
func (s *projectRequirementService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectRequirement{ID: id}).Delete(&model.ProjectRequirement{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectRequirementService) PaginateBetweenTimes(condition *model.ProjectRequirement, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectRequirement, err error) {
	tx := database.DB.Model(&model.ProjectRequirement{})
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

	// 大字段不查询
	tx.Omit("detail", "full_path")

	err = tx.Find(&list, &model.ProjectRequirement{
		ProjectID:   condition.ProjectID,
		StageID:     condition.StageID,
		Type:        condition.Type,
		Priority:    condition.Priority,
		Source:      condition.Source,
		OwnerID:     condition.OwnerID,
		Feasibility: condition.Feasibility,
		Status:      condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectRequirementService) Instance(id uint64) (instance *model.ProjectRequirement, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectRequirement{}
	if err = database.DB.Where(&model.ProjectRequirement{ID: id}).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectRequirementService) UpdateStatus(id uint64, status int) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}

	// 获取原始状态
	var oldStatus int
	err = database.DB.Model(&model.ProjectRequirement{}).Where(&model.ProjectRequirement{ID: id}).Select("status").First(&oldStatus).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln(err)
			return
		}
		return false, nil
	}

	// 如果状态相同，则不更新
	if oldStatus == status {
		return true, nil
	}

	// 判断当前状态是否可改变为目标状态
	canChangeStatus := false
	switch oldStatus {
	case model.ProjectRequirementStatusPendingDesign: // 待设计
		fallthrough
	case model.ProjectRequirementStatusRejected: // 评审不通过
		if status == model.ProjectRequirementStatusPendingReview {
			canChangeStatus = true
			break
		}
	case model.ProjectRequirementStatusPendingReview: // 待评审
		//ProjectRequirementStatusReviewed || ProjectRequirementStatusRejected
		if status == model.ProjectRequirementStatusReviewed || status == model.ProjectRequirementStatusRejected {
			canChangeStatus = true
			break
		}
	case model.ProjectRequirementStatusReviewed: // 评审通过
		//ProjectRequirementStatusCompleted
		if status == model.ProjectRequirementStatusCompleted {
			canChangeStatus = true
			break
		}

	}
	if canChangeStatus {
		err = database.DB.Model(&model.ProjectRequirement{ID: id}).Update("status", status).Error
		if err != nil {
			logger.Errorln(err)
			return
		}
		success = true
	}
	return
}
