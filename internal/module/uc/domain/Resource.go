package domain

import (
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type ResourceListRequest struct {
	model.Resource
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
