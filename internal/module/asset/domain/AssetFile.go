package domain

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type AssetFileListRequest struct {
	model.File
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	DeleteTimeCondition *server.TimeCondition `json:"deleteTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type AssetFileWithCreator struct {
	model.File
	Creator *ucModel.User `json:"creator"`
}
