package service

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var RoleService = new(roleService)

type roleService struct{}

// Add 添加角色
func (s *roleService) Add(instance *model.Role) (duplicated bool, success bool, err error) {
	if instance.Name == "" {
		err = errors.New("roleName is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.Role{}).Where(&model.Role{Name: instance.Name}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	// 设置默认值
	if instance.Type == 0 {
		instance.Type = 1
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

// Update 更新角色基本信息
func (s *roleService) Update(instance *model.Role) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.Role{ID: instance.ID}).Updates(&model.Role{
		Name:           instance.Name,
		Desc:           instance.Desc,
		Type:           instance.Type,
		DataPermission: instance.DataPermission,
		Status:         instance.Status,
		Style:          instance.Style,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	s.removeCache(instance.ID)
	return
}

// Delete 删除角色
func (s *roleService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Delete(&model.Role{ID: id}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *roleService) PaginateBetweenTimes(condition *model.Role, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.Role, err error) {
	tx := database.DB.Limit(100)
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

	// 模糊查询
	if condition != nil {
		if condition.Name != "" {
			tx = tx.Where("role_name like ?", "%"+condition.Name+"%")
		}
	}
	err = tx.Find(&list, &model.Role{
		ID:             condition.ID,
		Type:           condition.Type,
		DataPermission: condition.DataPermission,
		Status:         condition.Status,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取单个角色
func (s *roleService) Instance(condition *model.Role) (instance *model.Role, err error) {
	if condition.ID == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.Role{}
	err = database.DB.Where(condition).First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// ResourceCodes 获取角色的资源编码列表
func (s *roleService) ResourceCodes(roleId uint64) (list []string, err error) {
	list = make([]string, 0)
	err = database.DB.Model(&model.RoleResource{}).Where(&model.RoleResource{RoleID: roleId}).Pluck("resource_code", &list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *roleService) DispatchResources(roleID uint64, ResourceCodeList []string) (success bool, err error) {
	if roleID == 0 {
		err = errors.New("roleID is required")
		return
	}
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除旧的
		err = tx.Where(&model.RoleResource{RoleID: roleID}).Delete(&model.RoleResource{}).Error
		if err != nil {
			logger.Errorln(err)
			return err
		}
		// 添加新的
		for _, resourceCode := range ResourceCodeList {
			roleResource := &model.RoleResource{
				ID:           util.SnowflakeId(),
				RoleID:       roleID,
				ResourceCode: resourceCode,
			}
			err = tx.Create(roleResource).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// SetDefault 设置默认角色
func (*roleService) SetDefault(id uint64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Role{}).Where("default_role=?", 1).Updates(&model.Role{DefaultRole: -1}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if err := tx.Model(&model.Role{ID: id}).Updates(&model.Role{DefaultRole: 1}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		return nil
	})
}

func (s *roleService) removeCache(id uint64) {
	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	_, _ = conn.Do("HDEL", constant.RedisKeyRoleDataPerm, id)
}
