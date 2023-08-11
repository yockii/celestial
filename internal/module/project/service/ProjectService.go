package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	taskModel "github.com/yockii/celestial/internal/module/task/model"
	testModel "github.com/yockii/celestial/internal/module/test/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
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

		// 如果有父级项目，则父级项目的子项目数+1
		if instance.ParentID != 0 {
			if err = tx.Model(&model.Project{}).Where(&model.Project{ID: instance.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		return nil
	})
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectService) Update(instance, oldInstance *model.Project) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Where(&model.Project{ID: instance.ID}).Updates(&model.Project{
			Name:        instance.Name,
			Code:        instance.Code,
			Description: instance.Description,
			OwnerID:     instance.OwnerID,
			StageID:     instance.StageID,
		}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 如果父级发生变化，则更新父级项目的子项目数
		if instance.ParentID != 0 && instance.ParentID != oldInstance.ParentID {
			if err = tx.Model(&model.Project{}).Where(&model.Project{ID: instance.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			if err = tx.Model(&model.Project{}).Where(&model.Project{ID: oldInstance.ParentID}).Update("child_count", gorm.Expr("child_count - ?", 1)).Error; err != nil {
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

// Delete 删除资源
func (s *projectService) Delete(instance *model.Project) (success bool, err error) {
	if instance == nil || instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Where(&model.Project{ID: instance.ID}).Delete(&model.Project{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目成员
		err = tx.Where(&model.ProjectMember{ProjectID: instance.ID}).Delete(&model.ProjectMember{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的任务
		err = tx.Where(&taskModel.ProjectTask{ProjectID: instance.ID}).Delete(&taskModel.ProjectTask{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的任务成员
		err = tx.Where(&taskModel.ProjectTaskMember{ProjectID: instance.ID}).Delete(&taskModel.ProjectTaskMember{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的计划
		err = tx.Where(&model.ProjectPlan{ProjectID: instance.ID}).Delete(&model.ProjectPlan{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的功能模块
		err = tx.Where(&model.ProjectModule{ProjectID: instance.ID}).Delete(&model.ProjectModule{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的需求
		err = tx.Where(&model.ProjectRequirement{ProjectID: instance.ID}).Delete(&model.ProjectRequirement{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的测试
		err = tx.Where(&testModel.ProjectTest{ProjectID: instance.ID}).Delete(&testModel.ProjectTest{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的测试用例
		err = tx.Where(&testModel.ProjectTestCase{ProjectID: instance.ID}).Delete(&testModel.ProjectTestCase{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的测试用例项
		err = tx.Where(&testModel.ProjectTestCaseItem{ProjectID: instance.ID}).Delete(&testModel.ProjectTestCaseItem{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的测试用例项步骤
		err = tx.Where(&testModel.ProjectTestCaseItemStep{ProjectID: instance.ID}).Delete(&testModel.ProjectTestCaseItemStep{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的缺陷
		err = tx.Where(&model.ProjectIssue{ProjectID: instance.ID}).Delete(&model.ProjectIssue{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的变更
		err = tx.Where(&model.ProjectChange{ProjectID: instance.ID}).Delete(&model.ProjectChange{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的风险
		err = tx.Where(&model.ProjectRisk{ProjectID: instance.ID}).Delete(&model.ProjectRisk{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 删除项目关联的资产
		err = tx.Where(&model.ProjectAsset{ProjectID: instance.ID}).Delete(&model.ProjectAsset{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 如果有父级项目，更新父级项目的子项目数量
		if instance.ParentID != 0 {
			err = tx.Model(&model.Project{}).Where(&model.Project{ID: instance.ParentID}).Update("child_count", gorm.Expr("child_count - ?", 1)).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
		}

		return nil
	})

	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectService) PaginateBetweenTimes(condition *model.Project, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition, currentUserID uint64, dataPermit int) (total int64, list []*model.Project, err error) {
	tx := database.DB.Model(&model.Project{})
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
				tx = tx.Where(tc+" between ? and ?", time.Time(tr.Start).UnixMilli(), time.Time(tr.End).UnixMilli())
			} else if tr.Start.IsZero() && !tr.End.IsZero() {
				tx = tx.Where(tc+" <= ?", time.Time(tr.End).UnixMilli())
			} else if !tr.Start.IsZero() && tr.End.IsZero() {
				tx = tx.Where(tc+" > ?", time.Time(tr.Start).UnixMilli())
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

	if dataPermit == ucModel.RoleDataPermissionSelf {
		tx = tx.Where("id in (?)", database.DB.Model(&model.ProjectMember{}).Select("project_id").Where("user_id = ?", currentUserID))
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

func (s *projectService) MyProjects(uid uint64, condition *model.Project) (list []*model.Project, err error) {
	tx := database.DB.Model(&model.Project{})
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
		tx = tx.Where(&model.Project{
			ParentID: condition.ParentID,
			OwnerID:  condition.OwnerID,
			StageID:  condition.StageID,
		})
	}

	err = tx.Where("id in (?)", database.DB.Model(&model.ProjectMember{}).Distinct("project_id").Where("user_id = ?", uid)).Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectService) ListAllForWorkTimeStatistics() (list []*model.Project, err error) {
	err = database.DB.Model(&model.Project{}).
		Omit("description").
		Where("parent_id = ?", 0).
		Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
