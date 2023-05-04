package domain

import (
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTestCaseStepListRequest struct {
	model.ProjectTestCaseStep
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
