package service

import (
	"errors"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/cache"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"

	logger "github.com/sirupsen/logrus"
)

var SysConfigService = new(sysConfigService)

type sysConfigService struct{}

const sysConfigCacheSeconds = 60

func sysConfigCacheKey(key string) string {
	return constant.RedisKeySysConfig + ":" + key
}

// GetValue 读取系统配置值，带Redis缓存；未配置时返回空字符串
func (s *sysConfigService) GetValue(key string) (value string, err error) {
	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	cacheKey := sysConfigCacheKey(key)
	value, err = redis.String(conn.Do("GET", cacheKey))
	if err == nil {
		return
	}
	if err != redis.ErrNil {
		logger.Errorln(err)
	}

	instance := &model.SysConfig{}
	err = database.DB.Where(&model.SysConfig{Key: key}).First(instance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		logger.Errorln(err)
		return
	}
	value = instance.Value
	_, _ = conn.Do("SETEX", cacheKey, sysConfigCacheSeconds, value)
	return
}

// GetBool 读取布尔型配置，未配置或读取失败时返回fallback
func (s *sysConfigService) GetBool(key string, fallback bool) bool {
	value, err := s.GetValue(key)
	if err != nil || value == "" {
		return fallback
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		logger.Errorln(err)
		return fallback
	}
	return b
}

// List 全部配置
func (s *sysConfigService) List() (list []*model.SysConfig, err error) {
	err = database.DB.Order("config_key").Find(&list).Error
	if err != nil {
		logger.Errorln(err)
	}
	return
}

// Set 新增或更新配置，并使缓存立即失效，改动实时生效
func (s *sysConfigService) Set(key, value, comment string) error {
	instance := &model.SysConfig{}
	err := database.DB.Where(&model.SysConfig{Key: key}).First(instance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			instance = &model.SysConfig{
				ID:      util.SnowflakeId(),
				Key:     key,
				Value:   value,
				Comment: comment,
			}
			if err = database.DB.Create(instance).Error; err != nil {
				logger.Errorln(err)
				return err
			}
		} else {
			logger.Errorln(err)
			return err
		}
	} else {
		if err = database.DB.Model(&model.SysConfig{ID: instance.ID}).Updates(map[string]interface{}{
			"value":   value,
			"comment": comment,
		}).Error; err != nil {
			logger.Errorln(err)
			return err
		}
	}

	conn := cache.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)
	_, _ = conn.Do("DEL", sysConfigCacheKey(key))
	return nil
}
