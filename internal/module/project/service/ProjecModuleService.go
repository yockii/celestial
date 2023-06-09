package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/model"
	taskModel "github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"strings"
	"time"
)

var ProjectModuleService = new(projectModuleService)

type projectModuleService struct{}

// Add 添加资源
func (s *projectModuleService) Add(instance *model.ProjectModule) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.ProjectID == 0 {
		err = errors.New("name and projectId is required")
		return
	}

	var parent *model.ProjectModule
	if instance.ParentID != 0 {
		parent = &model.ProjectModule{ID: instance.ParentID}
		if err = database.DB.Where(parent).First(parent).Error; err != nil {
			logger.Errorln(err)
			return
		}
	}

	var c int64
	err = database.DB.Model(&model.ProjectModule{}).Where(&model.ProjectModule{
		ProjectID: instance.ProjectID,
		Name:      instance.Name,
		ParentID:  instance.ParentID,
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

	// 事务处理，更新父级的子数量、获取父级的路径
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if parent != nil {
			if err = tx.Model(parent).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			instance.FullPath = parent.FullPath + "/" + instance.Name
		} else {
			instance.FullPath = "/" + instance.Name
		}
		if err = tx.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		return nil
	}); err != nil {
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *projectModuleService) Update(instance, oldInstance *model.ProjectModule) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 检查parentId与原来的是否一致
	if oldInstance.ParentID != instance.ParentID || (instance.Name != "" && oldInstance.Name != instance.Name) {
		// 不一致则需要更新父级的子数量以及更新自己的路径
		var parent *model.ProjectModule
		if instance.ParentID != 0 {
			parent = &model.ProjectModule{ID: instance.ParentID}
			if err = database.DB.Where(parent).First(parent).Error; err != nil {
				logger.Errorln(err)
				return
			}
			instance.FullPath = parent.FullPath + "/" + instance.Name
		} else {
			instance.FullPath = "/" + instance.Name
		}

		if err = database.DB.Transaction(func(tx *gorm.DB) error {
			if oldInstance.ParentID != instance.ParentID {
				if parent != nil {
					if err = tx.Model(parent).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
						logger.Errorln(err)
						return err
					}
					// 旧的父级的子数量减1
					if oldInstance.ParentID != 0 {
						if err = tx.Model(&model.ProjectModule{}).Where(&model.ProjectModule{ID: oldInstance.ParentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
							logger.Errorln(err)
							return err
						}
					}
				}
			}
			if err = tx.Model(&model.ProjectModule{ID: instance.ID}).Updates(&model.ProjectModule{
				ParentID:  instance.ParentID,
				FullPath:  instance.FullPath,
				Name:      instance.Name,
				Alias:     instance.Alias,
				Remark:    instance.Remark,
				CreatorID: instance.CreatorID,
			}).Error; err != nil {
				logger.Errorln(err)
				return err
			}

			// 更新所有子模块的路径
			if err = tx.Model(&model.ProjectModule{}).Where("full_path like ?", oldInstance.FullPath+"%").
				Update("full_path", gorm.Expr("concat(?, substring(full_path, ?))", instance.FullPath, len(oldInstance.FullPath)+1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}

			// 更新所有需求的路径
			if err = tx.Model(&model.ProjectRequirement{}).Where("full_path like ?", oldInstance.FullPath+"%").
				Update("full_path", gorm.Expr("concat(?, substring(full_path, ?))", instance.FullPath, len(oldInstance.FullPath)+1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}

			// 更新所有任务的需求路径
			if err = tx.Model(&taskModel.ProjectTask{}).Where("full_path like ?", oldInstance.FullPath+"%").
				Update("full_path", gorm.Expr("concat(?, substring(full_path, ?))", instance.FullPath, len(oldInstance.FullPath)+1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}

			return nil
		}); err != nil {
			return
		}
	} else {
		// 一致则只更新基本信息
		if err = database.DB.Model(&model.ProjectModule{ID: instance.ID}).Updates(&model.ProjectModule{
			Name:      instance.Name,
			Alias:     instance.Alias,
			Remark:    instance.Remark,
			CreatorID: instance.CreatorID,
		}).Error; err != nil {
			logger.Errorln(err)
			return
		}
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectModuleService) Delete(instance *model.ProjectModule) (success bool, err error) {
	if instance == nil || instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 事务处理，如果有更新父级的子数量
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除
		if err = tx.Delete(&model.ProjectModule{ID: instance.ID}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if instance.ParentID != 0 {
			if err = tx.Model(&model.ProjectModule{ID: instance.ParentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}
		return nil
	}); err != nil {
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectModuleService) PaginateBetweenTimes(condition *model.ProjectModule, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectModule, err error) {
	tx := database.DB.Model(&model.ProjectModule{})
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
		if condition.Name != "" {
			tx = tx.Where("name like ?", "%"+condition.Name+"%")
		}
	}

	err = tx.Find(&list, &model.ProjectModule{
		ID:        condition.ID,
		ProjectID: condition.ProjectID,
		ParentID:  condition.ParentID,
		Status:    condition.Status,
		CreatorID: condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectModuleService) Instance(id uint64) (instance *model.ProjectModule, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectModule{}
	if err = database.DB.Where(&model.ProjectModule{ID: id}).First(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorln(err)
		return
	}
	return
}

func (s *projectModuleService) UpdateStatus(instance *model.ProjectModule, status int) (success bool, err error) {
	if instance == nil || instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	// 获取原始状态
	oldStatus := instance.Status
	if oldStatus == status {
		return true, nil
	}

	// 判断当前状态是否可改变为目标状态
	canChangeStatus := false
	switch oldStatus {
	case model.ProjectModuleStatusPendingReview: // 待评审
		fallthrough
	case model.ProjectModuleStatusRejected: // 评审不通过
		if status == model.ProjectModuleStatusPendingDev {
			canChangeStatus = true
		}
	case model.ProjectModuleStatusPendingDev: // 评审通过待开发
		if status == model.ProjectModuleStatusCompleted {
			canChangeStatus = true
		}
	}

	if canChangeStatus {
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&model.ProjectModule{}).Where(&model.ProjectModule{ID: instance.ID}).Update("status", status).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
			if status == model.ProjectModuleStatusPendingDev {
				// 模块评审通过后，所有上级模块都通过评审
				parentPathArray := strings.Split(instance.FullPath, "/")
				var parentFullPathList []string
				path := ""
				for _, parentPath := range parentPathArray {
					if parentPath == "" {
						continue
					}
					path += "/" + parentPath
					parentFullPathList = append(parentFullPathList, path)
				}
				err = tx.Model(&model.ProjectModule{}).Where(&model.ProjectModule{
					ProjectID: instance.ProjectID,
				}).Where("full_path in (?)", parentFullPathList).Update("status", model.ProjectModuleStatusPendingDev).Error
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
	}
	return
}
