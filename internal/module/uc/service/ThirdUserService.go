package service

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/database"
)

var ThirdUserService = new(thirdUserService)

type thirdUserService struct{}

func (s *thirdUserService) ListByUserIDListAndSourceCode(userIdList []uint64, sourceCode string) (list []*model.ThirdUser, err error) {
	if len(userIdList) == 0 {
		return
	}
	if err = database.DB.Where(&model.ThirdUser{
		SourceCode: sourceCode,
	}).Where("user_id in (?)", userIdList).Find(&list).Error; err != nil {
		logger.Error(err)
		return
	}
	return
}
