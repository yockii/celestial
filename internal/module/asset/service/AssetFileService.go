package service

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/domain"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/celestial/internal/module/asset/provider"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

var AssetFileService = &assetFileService{
	locker: new(sync.Mutex),
}

type assetFileService struct {
	osManager provider.OsManager
	locker    sync.Locker
}

func (s *assetFileService) Upload(suffix string, reader io.Reader) (objName string, err error) {
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}
	// 上传文件
	now := time.Now().Format("20060102")
	objName = fmt.Sprintf("%s/%d.%s", now, util.SnowflakeId(), suffix)
	if err = s.osManager.PutObject(objName, reader); err != nil {
		return
	}
	return
}

// Add 添加资源
func (s *assetFileService) Add(instance *model.File, reader io.Reader) (duplicated bool, success bool, err error) {
	if instance.Name == "" || instance.CategoryID == 0 {
		err = errors.New("Name and CategoryID is required ")
		return
	}
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}

	// 检查是否有重复的资源文件
	//var c int64
	//err = database.DB.Model(&model.File{}).Where(&model.File{
	//	CategoryID: instance.CategoryID,
	//	Name:       instance.Name,
	//	Suffix:     instance.Suffix,
	//}).Count(&c).Error
	//if err != nil {
	//	logger.Errorln(err)
	//	return
	//}
	//if c > 0 {
	//	duplicated = true
	//	return
	//}

	instance.ID = util.SnowflakeId()
	// 上传文件
	now := time.Now().Format("20060102")
	objName := fmt.Sprintf("%s/%d.%s", now, instance.ID, instance.Suffix)
	if err = s.osManager.PutObject(objName, reader); err != nil {
		return false, false, err
	}

	instance.OssConfigID = s.osManager.GetOssConfigID()
	instance.ObjName = objName

	// 设置CategoryPath
	category := new(model.AssetCategory)
	if err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{ID: instance.CategoryID}).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.New("category not found")
		}
		logger.Errorln(err)
		return
	}
	instance.CategoryPath = category.FullPath

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 添加文件版本信息
		if err = tx.Create(&model.FileVersion{
			ID:          util.SnowflakeId(),
			FileID:      instance.ID,
			OssConfigID: s.osManager.GetOssConfigID(),
			Size:        instance.Size,
			ObjName:     objName,
			CreatorID:   instance.CreatorID,
		}).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 设置创建者管理权限
		assetFilePermission := &model.FilePermission{
			ID:         util.SnowflakeId(),
			FileID:     instance.ID,
			UserID:     instance.CreatorID,
			Permission: 4,
		}
		if err = tx.Create(assetFilePermission).Error; err != nil {
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

// AddNewVersion 添加资源新版本, 其中instance必须有：id、size、creatorID、suffix
func (s *assetFileService) AddNewVersion(instance *model.File, reader io.Reader) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("ID is required ")
		return
	}
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}

	// 上传文件
	now := time.Now().Format("20060102")
	objName := fmt.Sprintf("%s/%d.%s", now, instance.ID, instance.Suffix)
	if err = s.osManager.PutObject(objName, reader); err != nil {
		return false, err
	}

	instance.OssConfigID = s.osManager.GetOssConfigID()
	instance.ObjName = objName

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where(&model.File{ID: instance.ID}).Updates(instance).Error; err != nil {
			logger.Errorln(err)
			return err
		}

		// 添加文件版本信息
		if err = tx.Create(&model.FileVersion{
			ID:          util.SnowflakeId(),
			FileID:      instance.ID,
			OssConfigID: s.osManager.GetOssConfigID(),
			Size:        instance.Size,
			ObjName:     objName,
			CreatorID:   instance.CreatorID,
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

// Update 更新资源基本信息
func (s *assetFileService) Update(instance *model.File) (success bool, err error) {
	if instance.ID == 0 {
		err = errors.New("id is required")
		return
	}

	// 检查是否变动了categoryID
	var old model.File
	if err = database.DB.Model(&model.File{}).Where(&model.File{ID: instance.ID}).First(&old).Error; err != nil {
		logger.Errorln(err)
		return
	}
	if old.CategoryID != instance.CategoryID {
		// 变动了categoryID，需要更新CategoryPath
		category := new(model.AssetCategory)
		if err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{ID: instance.CategoryID}).First(&category).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = errors.New("category not found")
			}
			logger.Errorln(err)
			return
		}
		instance.CategoryPath = category.FullPath
	}

	err = database.DB.Where(&model.File{ID: instance.ID}).Updates(&model.File{
		CategoryID:   instance.CategoryID,
		Name:         instance.Name,
		Suffix:       instance.Suffix,
		CreatorID:    instance.CreatorID,
		CategoryPath: instance.CategoryPath,
	}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// Delete 删除资源
func (s *assetFileService) Delete(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.File{ID: id}).Delete(&model.File{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

// PaginateBetweenTimes 带时间范围的分页查询
func (s *assetFileService) PaginateBetweenTimes(condition *model.File, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*model.File, err error) {
	tx := database.DB.Model(&model.File{})
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
		if condition.Suffix != "" {
			tx = tx.Where("suffix like ?", "%"+condition.Suffix+"%")
		}
		if condition.CategoryID != 0 {
			// 使用CategoryPath来查询
			category := new(model.AssetCategory)
			if err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{ID: condition.CategoryID}).First(&category).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = errors.New("category not found")
				}
				logger.Errorln(err)
				return
			}
			tx = tx.Where("category_path like ?", category.FullPath+"%")
		}
	}

	err = tx.Find(&list, &model.File{
		ID:        condition.ID,
		CreatorID: condition.CreatorID,
	}).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// PaginateDomainListBetweenTimes 带时间范围的分页查询 domain
func (s *assetFileService) PaginateDomainListBetweenTimes(condition *model.File, uid uint64, limit int, offset int, orderBy string, tcList map[string]*server.TimeCondition) (total int64, list []*domain.AssetFileWithCreator, err error) {
	tx := database.DB.Model(&model.File{})
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
		if condition.Suffix != "" {
			tx = tx.Where("suffix like ?", "%"+condition.Suffix+"%")
		}
		if condition.CategoryID != 0 {
			// 使用CategoryPath来查询
			category := new(model.AssetCategory)
			if err = database.DB.Model(&model.AssetCategory{}).Where(&model.AssetCategory{ID: condition.CategoryID}).First(&category).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = errors.New("category not found")
				}
				logger.Errorln(err)
				return
			}
			tx = tx.Where("category_path like ?", category.FullPath+"%")
		}
	}

	tx = tx.Where(&model.File{
		ID:        condition.ID,
		CreatorID: condition.CreatorID,
	})

	//tx = tx.Joins("t_file_permission", map[string]interface{}{})
	tx = tx.Joins("join t_file_permission on file_id=t_file.id and user_id=? and permission > 0", uid)
	tx = tx.Select("t_file.*, t_file_permission.permission as permission")

	err = tx.Find(&list).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

// Instance 获取资源实例
func (s *assetFileService) Instance(id uint64) (instance *model.File, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	instance = &model.File{}
	if err = database.DB.Where(&model.File{ID: id}).First(instance).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (s *assetFileService) initOsManager() (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.osManager == nil {
		// 从数据库中取出可用的云存储配置
		c := &model.OssConfig{}
		if err = database.DB.Where(&model.OssConfig{Status: 1}).First(c).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = errors.New("no available oss config")
			}
			return
		}
		// 初始化云存储
		osm := provider.GetProvider(c)
		if err = osm.Auth(); err != nil {
			return
		}
		s.osManager = osm
	}
	return nil
}

func (s *assetFileService) Download(id uint64) (reader io.ReadCloser, err error) {
	instance, err := s.Instance(id)
	if err != nil {
		return nil, err
	}
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}
	return s.osManager.GetObject(instance.ObjName)
}

func (s *assetFileService) DownloadByObjName(objName string) (reader io.ReadCloser, err error) {
	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}
	return s.osManager.GetObject(objName)
}

func (s *assetFileService) DownloadVersion(id uint64) (reader io.ReadCloser, err error) {
	versionedFile := new(model.FileVersion)
	if err = database.DB.Where(&model.FileVersion{ID: id}).First(versionedFile).Error; err != nil {
		logger.Errorln(err)
		return
	}

	if s.osManager == nil {
		// 初始化新的osManager
		if err = s.initOsManager(); err != nil {
			return
		}
	}
	return s.osManager.GetObject(versionedFile.ObjName)
}

func (s *assetFileService) GetPermissionUsers(condition *domain.FilePermissionUser) ([]*domain.FilePermissionUser, error) {
	tx := database.DB.Model(&model.FilePermission{})
	realName := condition.RealName
	condition.RealName = ""
	tx = tx.Where(&model.FilePermission{
		FileID:     condition.FileID,
		Permission: condition.Permission,
	})
	tx = tx.Select("t_file_permission.*, t_user.real_name as real_name")
	if realName != "" {
		tx = tx.Joins("join t_user on user_id=t_user.id and real_name like ?", "%"+realName+"%")
	} else {
		tx = tx.Joins("join t_user on user_id=t_user.id")
	}
	var list []*domain.FilePermissionUser
	err := tx.Find(&list).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return list, nil
}

func (s *assetFileService) UpdateFileUserPermission(permission *model.FilePermission) error {
	// 检查是否存在
	instance := &model.FilePermission{FileID: permission.FileID, UserID: permission.UserID}
	if err := database.DB.Where(instance).First(instance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			err = database.DB.Create(&model.FilePermission{
				ID:         util.SnowflakeId(),
				FileID:     permission.FileID,
				UserID:     permission.UserID,
				Permission: permission.Permission,
			}).Error
			if err != nil {
				logger.Errorln(err)
				return err
			}
			return nil
		} else {
			logger.Errorln(err)
			return err
		}
	}
	// 更新permission
	err := database.DB.Where(&model.FilePermission{ID: instance.ID}).Updates(&model.FilePermission{Permission: permission.Permission}).Error
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

func (s *assetFileService) CheckPermission(fileID uint64, uid uint64, permission uint8) (bool, error) {
	instance := &model.FilePermission{
		FileID: fileID,
		UserID: uid,
	}
	if err := database.DB.Where(instance).First(instance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		logger.Errorln(err)
		return false, err
	}
	if instance.Permission >= permission {
		return true, nil
	}
	return false, nil
}

func (s *assetFileService) DeleteFilePermission(id uint64) (success bool, err error) {
	if id == 0 {
		err = errors.New("id is required")
		return
	}
	err = database.DB.Where(&model.FilePermission{ID: id}).Delete(&model.FilePermission{}).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	success = true
	return
}

func (s *assetFileService) LatestVersion(fileID uint64) (*model.FileVersion, error) {
	version := &model.FileVersion{}
	err := database.DB.Where(&model.FileVersion{FileID: fileID}).Order("create_time desc").First(version).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return version, nil
}

func (s *assetFileService) GetFileVersion(condition *model.FileVersion) (*model.FileVersion, error) {
	if condition.ID == 0 && condition.FileID == 0 {
		return nil, errors.New("id or file_id is required")
	}
	version := &model.FileVersion{}
	err := database.DB.Where(condition).Order("create_time desc").First(version).Error
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return version, nil
}

func (s *assetFileService) VersionList(condition *model.FileVersion, limit int, offset int, orderBy string) (total int64, list []*model.FileVersion, err error) {
	if condition.FileID == 0 {
		err = errors.New("file_id is required")
		return
	}
	tx := database.DB.Model(&model.FileVersion{})
	if limit > -1 {
		tx = tx.Limit(limit)
	}
	if offset > -1 {
		tx = tx.Offset(offset)
	}
	if orderBy != "" {
		tx = tx.Order(orderBy)
	} else {
		tx = tx.Order("create_time desc")
	}
	tx = tx.Where(condition)

	err = tx.Limit(limit).Offset(offset).Find(&list).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}
