package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"time"
)

var ThirdSourceService = new(thirdSourceService)

type thirdSourceService struct{}

// Add 添加第三方登录源
func (s *thirdSourceService) Add(instance *model.ThirdSource) (duplicated bool, success bool, err error) {
	// 判断必填
	if instance.SourceName == "" || instance.SourceCode == "" {
		err = errors.New("必填项不能为空")
		return
	}
	var c int64
	err = database.DB.Model(&model.ThirdSource{}).Where(&model.ThirdSource{SourceName: instance.SourceName, SourceCode: instance.SourceCode}).Count(&c).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
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

// Update 更新第三方登录源
func (s *thirdSourceService) Update(instance *model.ThirdSource) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.ThirdSource{ID: instance.ID}).Updates(&model.ThirdSource{
		SourceName:    instance.SourceName,
		SourceCode:    instance.SourceCode,
		Configuration: instance.Configuration,
		MatchConfig:   instance.MatchConfig,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除第三方登录源
func (s *thirdSourceService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Delete(&model.ThirdSource{ID: id}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *thirdSourceService) PaginateBetweenTimes(condition *model.ThirdSource, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.ThirdSource, err error) {
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
		if condition.SourceName != "" {
			tx = tx.Where("source_name like ?", "%"+condition.SourceName+"%")
		}
		if condition.SourceCode != "" {
			tx = tx.Where("source_code like ?", "%"+condition.SourceCode+"%")
		}
	}
	err = tx.Omit("configuration", "matchConfig").Find(&list, &model.ThirdSource{
		ID: condition.ID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取单个实例
func (s *thirdSourceService) Instance(condition *model.ThirdSource) (instance *model.ThirdSource, err error) {
	if condition.ID == 0 {
		err = errors.New("id is required")
		logger.Errorln(err)
		return
	}
	instance = &model.ThirdSource{}
	err = database.DB.Where(condition).First(instance).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
