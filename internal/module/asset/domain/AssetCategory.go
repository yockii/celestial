package domain

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/server"
)

type AssetCategoryListRequest struct {
	model.AssetCategory
	OnlyParent          bool                  `json:"onlyParent"`
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
