package domain

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/server"
)

type AssetFileListRequest struct {
	model.File
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	DeleteTimeCondition *server.TimeCondition `json:"deleteTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
