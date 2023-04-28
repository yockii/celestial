package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ProjectRiskService = new(projectRiskService)

type projectRiskService struct{}

// Add 添加资源
func (s *projectRiskService) Add(instance *model.ProjectRisk) (duplicated bool, success bool, err error) {
	if instance.ProjectID == 0 || instance.RiskName == "" {
		err = errors.New("plan Name and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectRisk{}).Where(&model.ProjectRisk{
		ProjectID: instance.ProjectID,
		RiskName:  instance.RiskName,
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
func (s *projectRiskService) Update(instance *model.ProjectRisk) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectRisk{ID: instance.ID}).Updates(&model.ProjectRisk{
		ProjectID:       instance.ProjectID,
		StageID:         instance.StageID,
		RiskName:        instance.RiskName,
		RiskProbability: instance.RiskProbability,
		RiskImpact:      instance.RiskImpact,
		RiskLevel:       instance.RiskLevel,
		Status:          instance.Status,
		StartTime:       instance.StartTime,
		EndTime:         instance.EndTime,
		Result:          instance.Result,
		CreateUserID:    instance.CreateUserID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectRiskService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectRisk{ID: id}).Delete(&model.ProjectRisk{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectRiskService) PaginateBetweenTimes(condition *model.ProjectRisk, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectRisk, err error) {
	tx := database.DB.Model(&model.ProjectRisk{}).Limit(100)
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
		if condition.RiskName != "" {
			tx = tx.Where("risk_name like ?", "%"+condition.RiskName+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectRisk{
		ProjectID:       condition.ProjectID,
		StageID:         condition.StageID,
		RiskProbability: condition.RiskProbability,
		RiskImpact:      condition.RiskImpact,
		RiskLevel:       condition.RiskLevel,
		CreateUserID:    condition.CreateUserID,
		Status:          condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectRiskService) Instance(id uint64) (instance *model.ProjectRisk, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectRisk{}
	if err = database.DB.Where(&model.ProjectRisk{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
