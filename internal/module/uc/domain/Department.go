package domain

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type DepartmentListRequest struct {
	model.Department
	OnlyParent          bool                  `json:"onlyParent"`
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type AddDepartmentUsersRequest struct {
	DepartmentID uint64   `json:"departmentId,string"`
	UserIDs      []uint64 `json:"userIds"`
}

func (r *AddDepartmentUsersRequest) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.DepartmentID = j.Get("departmentId").Uint()
	for _, v := range j.Get("userIds").Array() {
		r.UserIDs = append(r.UserIDs, v.Uint())
	}
	return nil
}
