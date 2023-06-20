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
		RiskDesc:        instance.RiskDesc,
		RiskProbability: instance.RiskProbability,
		RiskImpact:      instance.RiskImpact,
		RiskLevel:       instance.RiskLevel,
		Status:          instance.Status,
		StartTime:       instance.StartTime,
		EndTime:         instance.EndTime,
		Result:          instance.Result,
		CreatorID:       instance.CreatorID,
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
	tx := database.DB.Model(&model.ProjectRisk{})
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
		if condition.RiskName != "" {
			tx = tx.Where("risk_name like ?", "%"+condition.RiskName+"%")
		}
	}

	// 忽略大字段
	tx.Omit("risk_desc", "response", "result")

	err = tx.Find(&list, &model.ProjectRisk{
		ID:              condition.ID,
		ProjectID:       condition.ProjectID,
		StageID:         condition.StageID,
		RiskProbability: condition.RiskProbability,
		RiskImpact:      condition.RiskImpact,
		RiskLevel:       condition.RiskLevel,
		CreatorID:       condition.CreatorID,
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectRiskService) CalculateRiskByProject(id uint64) (float64, *model.ProjectRisk, error) {
	// 计算项目风险系数，将该项目下所有除了已解决的风险的 概率、影响、等级 乘积求和
	var total float64
	var list []*model.ProjectRisk
	err := database.DB.Model(&model.ProjectRisk{}).Where("status <> ?", model.ProjectRiskStatusSolved).Where(&model.ProjectRisk{
		ProjectID: id,
	}).Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return total, nil, err
	}
	var maxRisk *model.ProjectRisk
	for _, v := range list {
		riskScore := float64(v.RiskProbability * v.RiskImpact * v.RiskLevel)
		// 根据风险状态，计算风险系数
		switch v.Status {
		// 已识别的风险，分数系数设置为0.5
		case model.ProjectRiskStatusIdentified:
			riskScore = riskScore * 0.5
		// 已应对的风险，分数系数设置为0.3
		case model.ProjectRiskStatusResponded:
			riskScore = riskScore * 0.3
		// 已发生的风险，分数系数设置为1
		case model.ProjectRiskStatusOccurred:
			riskScore = riskScore * 1
		default:
			riskScore = 0
		}
		total += riskScore
		if maxRisk == nil || (maxRisk.RiskProbability+maxRisk.RiskImpact+maxRisk.RiskLevel <= v.RiskProbability+v.RiskImpact+v.RiskLevel) {
			maxRisk = v
		}
	}
	return total, maxRisk, nil
}
