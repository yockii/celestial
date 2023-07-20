package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/database"
)

var AssetFilePermissionService = new(assetFilePermissionService)

type assetFilePermissionService struct{}

func (s assetFilePermissionService) Instance(condition *model.FilePermission) (permission *model.FilePermission, err error) {
	if condition.ID == 0 && (condition.FileID == 0 || condition.UserID == 0) {
		err = errors.New("id or fileId & userId is required")
		return
	}
	permission = new(model.FilePermission)
	if err = database.DB.Where(condition).First(permission).Error; err != nil {
		logger.Errorln(err)
		return
	}
	return permission, nil
}
