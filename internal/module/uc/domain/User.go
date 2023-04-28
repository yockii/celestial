package domain

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type UserListRequest struct {
	model.User
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type UserDispatchRolesRequest struct {
	UserID     uint64   `json:"userId,string"`
	RoleIDList []uint64 `json:"roleIdList"`
}

func (r *UserDispatchRolesRequest) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.UserID = j.Get("userId").Uint()
	for _, v := range j.Get("roleIdList").Array() {
		r.RoleIDList = append(r.RoleIDList, v.Uint())
	}
	return nil
}
