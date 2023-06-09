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

var ProjectPlanService = new(projectPlanService)

type projectPlanService struct{}

// Add 添加资源
func (s *projectPlanService) Add(instance *model.ProjectPlan) (duplicated bool, success bool, err error) {
	if instance.ProjectID == 0 || instance.PlanName == "" {
		err = errors.New("plan Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectPlan{}).Where(&model.ProjectPlan{
		ProjectID: instance.ProjectID,
		PlanName:  instance.PlanName,
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

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectPlanService) Update(instance *model.ProjectPlan) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectPlan{ID: instance.ID}).Updates(&model.ProjectPlan{
		ProjectID:    instance.ProjectID,
		StageID:      instance.StageID,
		PlanName:     instance.PlanName,
		PlanDesc:     instance.PlanDesc,
		StartTime:    instance.StartTime,
		EndTime:      instance.EndTime,
		Target:       instance.Target,
		Scope:        instance.Scope,
		Schedule:     instance.Schedule,
		Resource:     instance.Resource,
		Budget:       instance.Budget,
		CreateUserID: instance.CreateUserID,
		Status:       instance.Status,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectPlanService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectPlan{ID: id}).Delete(&model.ProjectPlan{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectPlanService) PaginateBetweenTimes(condition *model.ProjectPlan, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectPlan, err error) {
	tx := database.DB.Model(&model.ProjectPlan{})
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
		if condition.PlanName != "" {
			tx = tx.Where("plan_name like ?", "%"+condition.PlanName+"%")
		}
	}

	// 大字段不加载
	tx.Omit("plan_desc", "target", "scope", "schedule", "resource")

	err = tx.Find(&list, &model.ProjectPlan{
		ID:        condition.ID,
		ProjectID: condition.ProjectID,
		StageID:   condition.StageID,
		Status:    condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectPlanService) Instance(condition *model.ProjectPlan) (instance *model.ProjectPlan, err error) {
	if condition.ID == 0 && (condition.ProjectID == 0 || condition.Status != model.ProjectPlanStatusStarted) {
		err = errors.New("id is required or projectId with status started missing")
		return
	}
	instance = &model.ProjectPlan{}
	if err = database.DB.Where(condition).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}
