package domain

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectMemberListRequest struct {
	model.ProjectMember
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectMemberLite struct {
	UserID   uint64 `json:"userId,string"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	RoleID   uint64 `json:"roleId,string"`
}

type ProjectMemberBatchAddRequest struct {
	ProjectID  uint64   `json:"projectId,string"`
	RoleIdList []uint64 `json:"roleIdList"`
	UserIdList []uint64 `json:"userIdList"`
}

func (r *ProjectMemberBatchAddRequest) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ProjectID = j.Get("projectId").Uint()
	if j.Get("roleIdList").Exists() {
		for _, v := range j.Get("roleIdList").Array() {
			r.RoleIdList = append(r.RoleIdList, v.Uint())
		}
	}
	if j.Get("userIdList").Exists() {
		for _, v := range j.Get("userIdList").Array() {
			r.UserIdList = append(r.UserIdList, v.Uint())
		}
	}
	return nil
}
