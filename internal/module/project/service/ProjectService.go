package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectService = new(projectService)

type projectService struct{}

// Add 添加资源
func (s *projectService) Add(instance *model.Project) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.Code == "" {
		err = errors.New("projectName and projectCode is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.Project{}).Where("name = ? or code = ?", instance.Name, instance.Code).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()

	// 事务处理
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		// 将创建者添加到项目成员中
		if err = tx.Create(&model.ProjectMember{
			ID:        util.SnowflakeId(),
			ProjectID: instance.ID,
			UserID:    instance.OwnerID,
		}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		return nil
	})
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectService) Update(instance *model.Project) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.Project{ID: instance.ID}).Updates(&model.Project{
		Name:        instance.Name,
		Code:        instance.Code,
		Description: instance.Description,
		OwnerID:     instance.OwnerID,
		StageID:     instance.StageID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.Project{ID: id}).Delete(&model.Project{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectService) PaginateBetweenTimes(condition *model.Project, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Project, err error) {
	tx := database.DB.Model(&model.Project{}).Limit(100)
	if limit > -1 {
		tx = tx.Limit(limit)
	}
	if offset > -1 {
		tx = tx.Offset(offset)
	}
	if orderBy != "" {
		tx = tx.Order(orderBy)
	} else {
		tx = tx.Order("update_time desc")
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
		if condition.Code != "" {
			tx = tx.Where("code like ?", "%"+condition.Code+"%")
		}
		if condition.ParentID == 0 {
			tx = tx.Where("parent_id = ?", condition.ParentID)
		}
	}

	err = tx.Find(&list, &model.Project{
		ParentID: condition.ParentID,
		OwnerID:  condition.OwnerID,
		StageID:  condition.StageID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectService) Instance(id uint64) (instance *model.Project, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.Project{}
	if err = database.DB.Where(&model.Project{ID: id}).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}

// StatisticsByStage 统计各阶段的项目数量
func (s *projectService) StatisticsByStage() (list []*domain.ProjectCountByStage, err error) {
	err = database.DB.Model(&model.Project{}).Select("stage_id, count(*) as count").Group("stage_id").Scan(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
