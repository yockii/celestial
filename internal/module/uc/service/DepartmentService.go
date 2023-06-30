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

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(instance).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 更新父部门的子部门数量
		if instance.ParentID != 0 {
			err = tx.Model(&model.Department{ID: instance.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error
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
	return
}

// Update 更新部门基本信息
func (s *departmentService) Update(instance, oldInstance *model.Department) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 旧数据的parentID与新数据不同，需要更新fullPath及childCount
	var parent *model.Department
	if (instance.ParentID != 0 && instance.ParentID != oldInstance.ParentID) || instance.Name != oldInstance.Name {
		parent = &model.Department{ID: instance.ParentID}
		err = database.DB.First(parent).Error
		if err != nil {
			logger.Errorln(err)
			return
		}
		instance.FullPath = parent.FullPath + "/" + instance.Name
	} else {
		instance.FullPath = ""
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&model.Department{ID: instance.ID}).Updates(&model.Department{
			Name:         instance.Name,
			ExternalID:   instance.ExternalID,
			ExternalJson: instance.ExternalJson,
			OrderNum:     instance.OrderNum,
			FullPath:     instance.FullPath,
		}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 如果fullPath变更，则所有用户部门的DepartmentPath也要变更
		if instance.FullPath != "" && instance.FullPath != oldInstance.FullPath {
			err = tx.Model(&model.UserDepartment{}).Where(&model.UserDepartment{DepartmentID: instance.ID}).Updates(&model.UserDepartment{
				DepartmentPath: instance.FullPath,
			}).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
		}

		// 更新父部门的子部门数量
		if instance.ParentID != 0 && instance.ParentID != oldInstance.ParentID {
			// 原有的父级要-1
			err = tx.Model(&model.Department{ID: oldInstance.ParentID}).Update("child_count", gorm.Expr("child_count - ?", 1)).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
			// 新的父级要+1
			err = tx.Model(&model.Department{ID: instance.ParentID}).Update("child_count", gorm.Expr("child_count + ?", 1)).Error
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
func (s *departmentService) PaginateBetweenTimes(condition *model.Department, onlyParent bool, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Department, err error) {
	tx := database.DB.Model(&model.Department{})
	if limit > -1 {
		tx = tx.Limit(limit)
	}
	if offset > -1 {
		tx = tx.Offset(offset)
	}
	if orderBy != "" {
		tx = tx.Order(orderBy)
	}
	tx = tx.Order("order_num")

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
		if condition.FullPath != "" {
			tx = tx.Where("full_path like ?", condition.FullPath+"%")
		}
	}

	if onlyParent {
		tx = tx.Where("parent_id = 0")
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

	//获取部门信息
	department, err := s.Instance(instance.DepartmentID)
	if err != nil {
		return
	}
	if department == nil {
		err = errors.New("department not found")
		return
	}

	// 检查是否已经存在
	var c int64
	err = database.DB.Model(&model.UserDepartment{}).Where(&model.UserDepartment{
		DepartmentID:   instance.DepartmentID,
		UserID:         instance.UserID,
		DepartmentPath: department.FullPath,
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

// AddUsers 部门添加用户
func (s *departmentService) AddUsers(departmentID uint64, userIDList []uint64) (success bool, err error) {
	if departmentID == 0 {
		err = errors.New("department id is required")
		return
	}
	if len(userIDList) == 0 {
		err = errors.New("user id is required")
		return
	}

	// 取出已有的部门用户ID
	var oldUserIDList []uint64
	err = database.DB.Model(&model.UserDepartment{}).Where(&model.UserDepartment{DepartmentID: departmentID}).Pluck("user_id", &oldUserIDList).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 选出要添加的用户ID
	var newUserIDList []uint64
	for n := len(userIDList) - 1; n >= 0; n-- {
		var exist bool
		for i := len(oldUserIDList) - 1; i >= 0; i-- {
			if userIDList[n] == oldUserIDList[i] {
				exist = true
				// 移除, 注意最后一个
				if i == len(oldUserIDList)-1 {
					oldUserIDList = oldUserIDList[:i]
				} else {
					oldUserIDList = append(oldUserIDList[:i], oldUserIDList[i+1:]...)
				}
				break
			}
		}
		if !exist {
			newUserIDList = append(newUserIDList, userIDList[n])
			// 移除
			if n == len(userIDList)-1 {
				userIDList = userIDList[:n]
			} else {
				userIDList = append(userIDList[:n], userIDList[n+1:]...)
			}
		}
	}
	// 选出要删除的用户ID
	var delUserIDList []uint64
	for _, oldUserID := range oldUserIDList {
		var exist bool
		for n := len(userIDList) - 1; n >= 0; n-- {
			if oldUserID == userIDList[n] {
				exist = true
				// 移除
				if n == len(userIDList)-1 {
					userIDList = userIDList[:n]
				} else {
					userIDList = append(userIDList[:n], userIDList[n+1:]...)
				}
				break
			}
		}
		if !exist {
			delUserIDList = append(delUserIDList, oldUserID)
		}
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 批量添加用户部门
		if len(newUserIDList) > 0 {
			var userDepartmentList []model.UserDepartment
			for _, userID := range newUserIDList {
				userDepartmentList = append(userDepartmentList, model.UserDepartment{
					ID:           util.SnowflakeId(),
					DepartmentID: departmentID,
					UserID:       userID,
				})
			}
			err = tx.Create(&userDepartmentList).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
		}
		// 批量删除用户部门
		if len(delUserIDList) > 0 {
			err = tx.Where(&model.UserDepartment{
				DepartmentID: departmentID,
			}).Where("user_id in (?)", delUserIDList).
				Delete(&model.UserDepartment{}).Error
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
