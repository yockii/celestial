package domain

import (
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTestCaseItemListRequest struct {
	model.ProjectTestCaseItem
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectTestCaseItemWithSteps struct {
	model.ProjectTestCaseItem
	Steps []*model.ProjectTestCaseItemStep `json:"steps"`
}
