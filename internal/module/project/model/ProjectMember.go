package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type ProjectMember struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID  uint64 `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	UserID     uint64 `json:"userId,omitempty,string" gorm:"index;comment:用户ID"`
	RoleID     uint64 `json:"roleId,omitempty,string" gorm:"index;comment:在项目中承担的角色ID"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *ProjectMember) TableComment() string {
	return "项目成员表"
}

func (p *ProjectMember) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.ProjectID = j.Get("projectId").Uint()
	p.UserID = j.Get("userId").Uint()
	p.RoleID = j.Get("roleId").Uint()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectMember{})
}
