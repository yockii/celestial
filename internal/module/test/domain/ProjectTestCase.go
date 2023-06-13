package domain

import (
	"github.com/yockii/celestial/internal/module/test/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTestCaseListRequest struct {
	model.ProjectTestCase
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectTestCaseWithItems struct {
	model.ProjectTestCase
	Items []*model.ProjectTestCaseItem `json:"items"`
}

type ProjectTestCaseWithItemsWithSteps struct {
	model.ProjectTestCase
	Items []*ProjectTestCaseItemWithSteps `json:"items"`
}
