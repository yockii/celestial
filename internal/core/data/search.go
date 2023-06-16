package data

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/database"
)

func AddDocument(document *search.Document, uidList ...uint64) error {
	// 如果有uidList，则查询所有uid对应的用户，将用户的realName加入到document的对应字段中
	if len(uidList) > 0 {
		var realNameList []string
		if err := database.DB.Model(&model.User{}).Where(uidList).Pluck("real_name", &realNameList).Error; err != nil {
			logger.Error(err)
			return err
		}
		document.RelatedUsers = realNameList
	}

	_, err := search.AddDocument([]*search.Document{
		document,
	})
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func DeleteDocuments(idList ...uint64) error {
	_, err := search.DeleteDocuments(idList...)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func AddDocumentAntsWrapper(document *search.Document, uidList ...uint64) func() {
	return func() {
		if err := AddDocument(document, uidList...); err != nil {
			logger.Error(err)
		}
	}
}

func DeleteDocumentsAntsWrapper(idList ...uint64) func() {
	return func() {
		if err := DeleteDocuments(idList...); err != nil {
			logger.Error(err)
		}
	}
}
