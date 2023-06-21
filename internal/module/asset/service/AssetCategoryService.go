package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var AssetCategoryService = new(assetCategoryService)

type assetCategoryService struct{}

// Add 添加资源
func (s *assetCategoryService) Add(instance *model.AssetCategory) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.Type == 0 {
		err = errors.New("Name and type is required ")
		return
	}
	var c int64
	err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{
		ParentID: instance.ParentID,
		Name:     instance.Name,
		Type:     instance.Type,
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
	// 全路径设置
	if instance.ParentID != 0 {
		var parent model.AssetCategory
		if err = database.DB.Select("full_path").Model(&model.AssetCategory{ID: instance.ParentID}).First(&parent).Error; err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = parent.FullPath + "/" + instance.Name
	} else {
		instance.FullPath = instance.Name
	}

	// 事务处理，创建资源的同时，更新父级的子级数量
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(instance).Error; err != nil {
			return err
		}
		if instance.ParentID != 0 {
			if err = tx.Model(&model.AssetCategory{ID: instance.ParentID}).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新资源基本信息
func (s *assetCategoryService) Update(instance *model.AssetCategory) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 检查parentId是否与旧数据相同，不同则更新各自的parent的子级数量
	var oldParentID uint64
	if err = database.DB.Model(&model.AssetCategory{ID: instance.ID}).Select("parent_id").Scan(&oldParentID).Error; err != nil {
		logger.Errorln(err)
		return
	}

	if oldParentID != instance.ParentID {
		// 更新fullPath
		if instance.ParentID != 0 {
			var parent model.AssetCategory
			if err = database.DB.Select("full_path").Model(&model.AssetCategory{ID: instance.ParentID}).First(&parent).Error; err != nil {
				logger.Errorln(err)
				return
			}
			instance.FullPath = parent.FullPath + "/" + instance.Name
		} else {
			instance.FullPath = instance.Name
		}
	}

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if oldParentID != instance.ParentID {
			if err = tx.Model(&model.AssetCategory{ID: oldParentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			if err = tx.Model(&model.AssetCategory{ID: instance.ParentID}).Update("children_count", gorm.Expr("children_count + ?", 1)).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		}

		updateColumns := []string{
			"parent_id", "name", "type",
		}
		if instance.FullPath != "" {
			updateColumns = append(updateColumns, "full_path")
		}

		if err = tx.Where(&model.AssetCategory{ID: instance.ID}).Select(updateColumns).Updates(&model.AssetCategory{
			ParentID: instance.ParentID,
			Name:     instance.Name,
			Type:     instance.Type,
			FullPath: instance.FullPath,
		}).Error; err != nil {
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

// Delete 删除资源
func (s *assetCategoryService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}

	// 事务处理，删除资源的同时，更新父级的子级数量
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		var parentID uint64
		if err = tx.Model(&model.AssetCategory{ID: id}).Select("parent_id").Scan(&parentID).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if err = tx.Where(&model.AssetCategory{ID: id}).Delete(&model.AssetCategory{}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if parentID != 0 {
			if err = tx.Model(&model.AssetCategory{ID: parentID}).Update("children_count", gorm.Expr("children_count - ?", 1)).Error; err != nil {
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
func (s *assetCategoryService) PaginateBetweenTimes(condition *model.AssetCategory, onlyParent bool, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.AssetCategory, err error) {
	tx := database.DB.Model(&model.AssetCategory{})
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

	if condition.ParentID == 0 && onlyParent {
		tx = tx.Where("parent_id = 0")
	}

	err = tx.Find(&list, &model.AssetCategory{
		Type:      condition.Type,
		ParentID:  condition.ParentID,
		CreatorID: condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *assetCategoryService) Instance(id uint64) (instance *model.AssetCategory, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.AssetCategory{}
	if err = database.DB.Where(&model.AssetCategory{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}
