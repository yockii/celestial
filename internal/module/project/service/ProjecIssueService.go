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

	if err = database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
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
		ProjectID:    instance.ProjectID,
		Title:        instance.Title,
		Content:      instance.Content,
		Type:         instance.Type,
		Status:       instance.Status,
		AssignUserID: instance.AssignUserID,
		StartTime:    instance.StartTime,
		EndTime:      instance.EndTime,
		SolveTime:    instance.SolveTime,
		IssueCause:   instance.IssueCause,
		SolveMethod:  instance.SolveMethod,
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
	tx := database.DB.Model(&model.ProjectIssue{}).Limit(100)
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
		if condition.Title != "" {
			tx = tx.Where("title like ?", "%"+condition.Title+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectIssue{
		Type:         condition.Type,
		ProjectID:    condition.ProjectID,
		AssignUserID: condition.AssignUserID,
		Status:       condition.Status,
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
		logger.Errorln(err)
		return
	}
	return
}
