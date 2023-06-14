package domain

import (
	"github.com/tidwall/gjson"
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

func (o *ProjectTestCaseWithItems) UnmarshalJSON(data []byte) error {
	j := gjson.ParseBytes(data)
	err := o.ProjectTestCase.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	for _, itemJson := range j.Get("items").Array() {
		item := new(model.ProjectTestCaseItem)
		err = item.UnmarshalJSON([]byte(itemJson.Raw))
		if err != nil {
			return err
		}
		o.Items = append(o.Items, item)
	}

	return nil
}

type ProjectTestCaseWithItemsWithSteps struct {
	model.ProjectTestCase
	Items []*ProjectTestCaseItemWithSteps `json:"items"`
}
