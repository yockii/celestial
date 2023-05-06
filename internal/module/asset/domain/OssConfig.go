package domain

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/server"
)

type OssConfigListRequest struct {
	model.OssConfig
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
