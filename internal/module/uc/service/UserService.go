package service

import (
	"errors"
	"github.com/yockii/celestial/internal/module/uc/model"
	"gorm.io/gorm"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"golang.org/x/crypto/bcrypt"
)

var UserService = new(userService)

type userService struct{}

// LoginWithUsernameAndPassword 用户登录
func (s *userService) LoginWithUsernameAndPassword(username, password string) (instance *model.User, passwordNotMatch bool, err error) {
	instance = new(model.User)
	err = database.DB.Where(&model.User{Username: username}).First(instance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			passwordNotMatch = true
			return
		}
		logger.Errorln(err)
		return
	}
	if instance.Status != model.UserStatusNormal {
		err = errors.New("用户已被禁用")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(instance.Password), []byte(password))
	if err != nil {
		passwordNotMatch = true
		err = nil
		return
	}
	// 完成后密码置空
	instance.Password = ""
	return
}

// Add 添加用户
func (s *userService) Add(instance *model.User) (duplicated bool, success bool, err error) {
	if instance.Username == "" {
		err = errors.New("username is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.User{}).Where(&model.User{Username: instance.Username}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	instance.ID = util.SnowflakeId()
	if instance.Password != "" {
		pwd, _ := bcrypt.GenerateFromPassword([]byte(instance.Password), bcrypt.DefaultCost)
		instance.Password = string(pwd)
	}
	instance.Status = model.UserStatusNormal

	// 获取默认角色
	defaultRole := &model.Role{DefaultRole: 1}
	if err = database.DB.Where(defaultRole).First(defaultRole).Error; err != nil {
		logger.Errorln(err)
		return
	}
	if defaultRole != nil && defaultRole.ID > 0 {
		// 添加用户的同时要添加默认角色
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			if err = tx.Create(instance).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			userRole := &model.UserRole{
				ID:     util.SnowflakeId(),
				UserID: instance.ID,
				RoleID: defaultRole.ID,
			}
			if err = tx.Create(userRole).Error; err != nil {
				logger.Errorln(err)
				return err
			}
			return nil
		})
	} else {
		err = database.DB.Create(instance).Error
	}
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 完成后密码置空
	instance.Password = ""
	success = true
	return
}

// Update 更新用户基本信息
func (s *userService) Update(instance *model.User) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.User{ID: instance.ID}).Updates(&model.User{
		RealName: instance.RealName,
		Status:   instance.Status,
		Email:    instance.Email,
		Mobile:   instance.Mobile,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// UpdatePassword 更新用户密码
func (s *userService) UpdatePassword(instance *model.User) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	if instance.Password == "" {
		err = errors.New("password is required")
		return
	}
	pwd, _ := bcrypt.GenerateFromPassword([]byte(instance.Password), bcrypt.DefaultCost)
	err = database.DB.Where(&model.User{ID: instance.ID}).Updates(&model.User{
		Password: string(pwd),
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *userService) PaginateBetweenTimes(condition *model.User, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition, departmentPath string) (total int64, list []*model.User, err error) {
	// 处理不允许查询的字段
	if condition.Password != "" {
		condition.Password = ""
	}
	tx := database.DB
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
	if departmentPath != "" {
		tx = tx.Where("id in (?)", database.DB.Model(&model.UserDepartment{}).Where("department_path like ?", departmentPath+"%").Select("user_id"))
	}

	// 模糊查找
	if condition.Username != "" {
		tx = tx.Where("username like ?", "%"+condition.Username+"%")
		condition.Username = ""
	}
	if condition.RealName != "" {
		tx = tx.Where("real_name like ?", "%"+condition.RealName+"%")
		condition.RealName = ""
	}

	// 不查询离职的人
	if condition.Status == 0 {
		tx = tx.Where("status != ?", model.UserStatusLeaved)
	}

	err = tx.Omit("password").Find(&list, condition).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}
	return total, list, nil
}

// Instance 获取单个用户
func (s *userService) Instance(condition *model.User) (instance *model.User, err error) {
	if condition.ID == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.User{}
	err = database.DB.Where(condition).First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Delete 删除用户
func (s *userService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}

	if err = database.DB.Delete(&model.User{}, id).Error; err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Roles 获取用户的角色列表
func (s *userService) Roles(userId uint64, types ...int) (roles []*model.Role, err error) {
	// 获取用户ID对应的所有角色信息
	sm := gorm.Statement{DB: database.DB}
	_ = sm.Parse(&model.Role{})
	ruleTableName := sm.Schema.Table

	tx := database.DB.Model(&model.UserRole{})
	if len(types) > 0 {
		tx = tx.Where("type in (?)", types)
	}
	err = tx.
		Select(ruleTableName + ".*").
		Joins("left join " + ruleTableName + " on " + ruleTableName + ".id = role_id").
		Where(&model.UserRole{UserID: userId}).Scan(&roles).Error

	//var list []*model.UserRole
	//err = database.DB.Where(&model.UserRole{UserID: userId}).Find(&list).Error
	//if err != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//var roleIds []uint64
	//for _, v := range list {
	//	roleIds = append(roleIds, v.RoleID)
	//}

	return
}

func (s *userService) DispatchRoles(userID uint64, roleIDList []uint64) (success bool, err error) {
	if userID == 0 {
		err = errors.New("id is required")
		return
	}
	// 在事务中处理
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有的角色
		if err = tx.Where(&model.UserRole{UserID: userID}).Delete(&model.UserRole{}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		// 添加新的角色
		for _, v := range roleIDList {
			if err = tx.Create(&model.UserRole{
				ID:     util.SnowflakeId(),
				UserID: userID,
				RoleID: v,
			}).Error; err != nil {
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

// UserLeaved 用户离职
func (s *userService) UserLeaved(corpId string, uid string) error {
	source, err := ThirdSourceService.Instance(&model.ThirdSource{
		CorpId: corpId,
	})
	if err != nil {
		return err
	}
	if source == nil {
		logger.Warn("未找到对应的第三方源配置")
		return nil
	}

	thirdUser := &model.ThirdUser{OpenID: uid, SourceID: source.ID}
	// 获取第三方用户
	if err = database.DB.Where(thirdUser).First(thirdUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("未找到对应的第三方用户")
			return nil
		}
		logger.Errorln(err)
		return err
	}
	// 找到后，得到对应的userId, 在事务中删除用户组织，将用户和第三方用户信息标记为离职
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where(&model.UserDepartment{UserID: thirdUser.UserID}).Delete(&model.UserDepartment{}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		// 将用户和第三方用户信息标记为离职
		if err = tx.Model(thirdUser).Update("status", model.UserStatusLeaved).Error; err != nil {
			logger.Errorln(err)
			return err
		}
		if err = tx.Model(&model.User{ID: thirdUser.UserID}).Update("status", model.UserStatusLeaved).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// TODO 还有其他用户相关的信息需要删除，如关联的项目、任务等

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
