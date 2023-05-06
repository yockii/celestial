package domain

import (
	"github.com/yockii/celestial/internal/module/log/model"
	"github.com/yockii/ruomu-core/server"
)

type LogListRequest struct {
	model.Log
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
