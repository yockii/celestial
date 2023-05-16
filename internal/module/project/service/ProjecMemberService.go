package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/project/domain"
	"github.com/yockii/celestial/internal/module/project/model"
	userModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"time"
)

var ProjectMemberService = new(projectMemberService)

type projectMemberService struct{}

// Add 添加资源
func (s *projectMemberService) Add(instance *model.ProjectMember) (duplicated bool, success bool, err error) {
	if instance.UserID == 0 || instance.ProjectID == 0 {
		err = errors.New("userId and projectId is required")
		return
	}
	var c int64
	err = database.DB.Model(&model.ProjectMember{}).Where(&model.ProjectMember{
		ProjectID: instance.ProjectID,
		UserID:    instance.UserID,
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
func (s *projectMemberService) Update(instance *model.ProjectMember) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	err = database.DB.Where(&model.ProjectMember{ID: instance.ID}).Updates(&model.ProjectMember{
		UserID: instance.UserID,
		RoleID: instance.RoleID,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *projectMemberService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ProjectMember{ID: id}).Delete(&model.ProjectMember{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *projectMemberService) PaginateBetweenTimes(condition *model.ProjectMember, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ProjectMember, err error) {
	tx := database.DB.Model(&model.ProjectMember{}).Limit(100)
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
		//if condition.Name != "" {
		//	tx = tx.Where("name like ?", "%"+condition.Name+"%")
		//}
	}

	err = tx.Find(&list, &model.ProjectMember{
		ProjectID: condition.ProjectID,
		UserID:    condition.UserID,
		RoleID:    condition.RoleID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *projectMemberService) Instance(id uint64) (instance *model.ProjectMember, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.ProjectMember{}
	if err = database.DB.Where(&model.ProjectMember{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// ListLiteByProjectID 获取某项目的所有成员，仅获取id/username/realName字段信息
func (s *projectMemberService) ListLiteByProjectID(projectID uint64) (list []*domain.ProjectMemberLite, err error) {
	if projectID == 0 {
		err = errors.New("projectID is required")
		return
	}
	stmt := &gorm.Statement{DB: database.DB}
	_ = stmt.Parse(&userModel.User{})
	userTableName := stmt.Schema.Table
	err = database.DB.Model(&model.ProjectMember{}).Select("user_id, username, real_name, role_id").
		Joins("left join " + userTableName + " on user_id = " + userTableName + ".id").Where(&model.ProjectMember{ProjectID: projectID}).Scan(&list).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
