package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectModuleListRequest struct {
	model.ProjectModule
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
