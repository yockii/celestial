package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var DepartmentService = new(departmentService)

type departmentService struct{}

// Add 添加部门
func (s *departmentService) Add(instance *model.Department) (duplicated bool, success bool, err error) {
	if instance.Name == "" {
		err = errors.New("name is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.Department{}).Where(&model.Department{Name: instance.Name, ParentID: instance.ParentID}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	if instance.ParentID != 0 {
		// 获取parent
		parent := &model.Department{ID: instance.ParentID}
		err = database.DB.First(parent).Error
		if err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = parent.FullPath + "/" + instance.Name
	} else {
		instance.FullPath = instance.Name
	}

	instance.ID = util.SnowflakeId()

	err = database.DB.Create(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Update 更新部门基本信息
func (s *departmentService) Update(instance *model.Department) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.Department{ID: instance.ID}).Updates(&model.Department{
		ExternalID:   instance.ExternalID,
		ExternalJson: instance.ExternalJson,
		OrderNum:     instance.OrderNum,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// UpdateName 更新部门名称
func (s *departmentService) UpdateName(instance *model.Department) (duplicated bool, success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	if instance.Name == "" {
		err = errors.New("name is required")
		return
	}

	// 获取原有数据
	old := &model.Department{ID: instance.ID}
	err = database.DB.First(old).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if old.Name == instance.Name {
		// 名称未变更
		success = true
		return
	}

	var c int64
	err = database.DB.Model(&model.Department{}).Where(&model.Department{Name: instance.Name, ParentID: old.ParentID}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	// 事务
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 更新当前部门名称
		err = tx.Model(&model.Department{ID: instance.ID}).Updates(map[string]interface{}{
			"name":      instance.Name,
			"full_path": old.FullPath[:len(old.FullPath)-len(old.Name)] + instance.Name,
		}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}

		// 更新子部门全路径
		err = tx.Model(&model.Department{}).Where("full_path LIKE ?", old.FullPath+"/%").Updates(map[string]interface{}{
			"full_path": gorm.Expr("REPLACE(full_path, ?, ?)", old.FullPath+"/", old.FullPath[:len(old.FullPath)-len(old.Name)]+instance.Name+"/"),
		}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}

		return nil
	})

	if err != nil {
		return
	}
	success = true
	return
}

// ChangeParent 修改父级部门
func (s *departmentService) ChangeParent(instance *model.Department) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	if instance.ParentID == 0 {
		// 变为根部门
		instance.FullPath = instance.Name
	} else {
		// 获取parent
		parent := &model.Department{ID: instance.ParentID}
		err = database.DB.First(parent).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("父级部门不存在")
			}
			logger.Errorln(err)
			return
		}
		instance.FullPath = parent.FullPath + "/" + instance.Name
	}

	err = database.DB.Where(&model.Department{ID: instance.ID}).Updates(map[string]interface{}{
		"parent_id": instance.ParentID,
		"Full_Path": instance.FullPath,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除部门
func (s *departmentService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}

	// 检查存在子部门不允许删除
	var c int64
	err = database.DB.Model(&model.Department{}).Where(&model.Department{ParentID: id}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		err = errors.New("存在子部门，不允许删除")
		return
	}

	err = database.DB.Where(&model.Department{ID: id}).Delete(&model.Department{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *departmentService) PaginateBetweenTimes(condition *model.Department, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Department, err error) {
	tx := database.DB.Model(&model.Department{}).Limit(100)
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

	err = tx.Omit("external_json").Find(&list, &model.Department{
		ExternalID: condition.ExternalID,
		ParentID:   condition.ParentID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取单个实例
func (s *departmentService) Instance(id uint64) (instance *model.Department, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.Department{ID: id}
	err = database.DB.First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// AddUser 部门添加用户
func (s *departmentService) AddUser(instance *model.UserDepartment) (success bool, err error) {
	if instance.DepartmentID == 0 {
		err = errors.New("department id is required")
		return
	}
	if instance.UserID == 0 {
		err = errors.New("user id is required")
		return
	}

	// 检查是否已经存在
	var c int64
	err = database.DB.Model(&model.UserDepartment{}).Where(&model.UserDepartment{
		DepartmentID: instance.DepartmentID,
		UserID:       instance.UserID,
	}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		success = true
		return
	}

	err = database.DB.Create(&model.UserDepartment{
		ID:           util.SnowflakeId(),
		DepartmentID: instance.DepartmentID,
		UserID:       instance.UserID,
		ExternalJson: instance.ExternalJson,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// RemoveUser 部门移除用户
func (s *departmentService) RemoveUser(departmentID uint64, userID uint64) (success bool, err error) {
	if departmentID == 0 {
		err = errors.New("department id is required")
		return
	}
	if userID == 0 {
		err = errors.New("user id is required")
		return
	}

	err = database.DB.Where(&model.UserDepartment{
		DepartmentID: departmentID,
		UserID:       userID,
	}).Delete(&model.UserDepartment{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

func (s *departmentService) DeptRemoved(deptId string) error {
	// 获取部门
	dept := &model.Department{ExternalID: deptId}
	if err := database.DB.First(dept).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("部门不存在")
			return nil
		}
		logger.Errorln(err)
		return err
	}

	// 事务删除部门和部门用户关联
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除部门
		if err := tx.Where(&model.Department{ID: dept.ID}).Delete(&model.Department{}).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 删除部门用户关联
		if err := tx.Where(&model.UserDepartment{DepartmentID: dept.ID}).Delete(&model.UserDepartment{}).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
