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

var ProjectChangeService = new(projectChangeService)

type projectChangeService struct{}

// Add 添加资源
func (s *projectChangeService) Add(instance *model.ProjectChange) (duplicated bool, success bool, err error) {
	if instance.Title == "" || instance.Type == 0 || instance.ProjectID == 0 {
		err = errors.New("title / projectId / type is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectChange{}).Where(&model.ProjectChange{
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
func (s *projectChangeService) Update(instance *model.ProjectChange) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectChange{ID: instance.ID}).Updates(&model.ProjectChange{
		Title:          instance.Title,
		Type:           instance.Type,
		ProjectID:      instance.ProjectID,
		Level:          instance.Level,
		Reason:         instance.Reason,
		Plan:           instance.Plan,
		Review:         instance.Review,
		Risk:           instance.Risk,
		Status:         instance.Status,
		ApplyUserID:    instance.ApplyUserID,
		ReviewerIDList: instance.ReviewerIDList,
		Result:         instance.Result,
		ReviewTime:     instance.ReviewTime,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectChangeService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectChange{ID: id}).Delete(&model.ProjectChange{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectChangeService) PaginateBetweenTimes(condition *model.ProjectChange, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectChange, err error) {
	tx := database.DB.Model(&model.ProjectChange{})
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

	err = tx.Find(&list, &model.ProjectChange{
		Type:      condition.Type,
		ProjectID: condition.ProjectID,
		Level:     condition.Level,
		Status:    condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectChangeService) Instance(id uint64) (instance *model.ProjectChange, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectChange{}
	if err = database.DB.Where(&model.ProjectChange{ID: id}).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}
