package domain

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"github.com/yockii/ruomu-core/server"
)

type CommonTestCaseListRequest struct {
	model.CommonTestCase
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type CommonTestCaseItemListRequest struct {
	model.CommonTestCaseItem
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type CommonTestCaseWithItem struct {
	model.CommonTestCase
	Items []*model.CommonTestCaseItem `json:"items"`
}
