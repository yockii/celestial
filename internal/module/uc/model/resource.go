package model

import (
	"github.com/tidwall/gjson"
)

type Resource struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ResourceName string `json:"resourceName,omitempty" gorm:"size:50;comment:资源名称"`             // 资源名称
	ResourceCode string `json:"resourceCode,omitempty" gorm:"size:50;uniqueIndex;comment:资源代码"` // 资源认证代码
	HttpMethod   string `json:"httpMethod,omitempty" gorm:"comment:http方法"`                     // http方法
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *Resource) TableComment() string {
	return "资源表"
}
func (r *Resource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ID = j.Get("id").Uint()
	r.ResourceName = j.Get("resourceName").String()
	r.ResourceCode = j.Get("resourceCode").String()
	r.HttpMethod = j.Get("httpMethod").String()
	return nil
}

type RoleResource struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	RoleID       uint64 `json:"roleId,omitempty,string"`
	ResourceCode string `json:"resourceCode,omitempty" gorm:"size:50;comment:资源代码"` // 资源认证代码
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *RoleResource) TableComment() string {
	return "角色资源表"
}
func (rr *RoleResource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	rr.ID = j.Get("id").Uint()
	rr.RoleID = j.Get("roleId").Uint()
	rr.ResourceCode = j.Get("resourceCode").String()
	return nil
}

func init() {
	Models = append(Models, &Resource{}, &RoleResource{})
}
