package domain

import (
	"github.com/yockii/celestial/internal/module/log/model"
	"github.com/yockii/ruomu-core/server"
)

type WorkTimeListRequest struct {
	model.WorkTime
	StartDateCondition  *server.TimeCondition `json:"startDateCondition"`
	EndDateCondition    *server.TimeCondition `json:"endDateCondition"`
	ReviewTimeCondition *server.TimeCondition `json:"reviewTimeCondition"`
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
