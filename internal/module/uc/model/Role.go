package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

const (
	RoleTypeSuperAdmin = -1
	RoleTypeNormal     = 1
	RoleTypeProject    = 2
)

const (
	RoleDataPermissionAll        = 1
	RoleDataPermissionDept       = 2
	RoleDataPermissionDeptAndSub = 3
	RoleDataPermissionSelf       = 4
)

const (
	RoleStatusNormal  = 1
	RoleStatusDisable = 2
)

type Role struct {
	ID             uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name           string `json:"name,omitempty" gorm:"size:50;comment:角色名称"`
	Desc           string `json:"desc,omitempty" gorm:"size:200;comment:角色描述"`
	Type           int    `json:"type,omitempty" gorm:"comment:角色类型 1-普通角色 2-项目角色 -1-超级管理员角色"`
	DataPermission int    `json:"dataPermission,omitempty" gorm:"comment:数据权限 1-全部数据权限 2-本部门及以下数据权限 3-仅本人数据权限"`
	Style          string `json:"style,omitempty" gorm:"size:500;comment:角色样式"`
	DefaultRole    int    `json:"defaultRole" gorm:"comment:默认角色 1-是 其他否"`
	Status         int    `json:"status,omitempty" gorm:"comment:状态 1-启用 2-禁用"`
	CreateTime     int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime     int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *Role) TableComment() string {
	return "角色表"
}
func (r *Role) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ID = j.Get("id").Uint()
	r.Name = j.Get("name").String()
	r.Desc = j.Get("desc").String()
	r.Type = int(j.Get("type").Int())
	r.DataPermission = int(j.Get("dataPermission").Int())
	r.DefaultRole = int(j.Get("defaultRole").Int())
	r.Status = int(j.Get("status").Int())
	r.Style = j.Get("style").String()
	return nil
}

type UserRole struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	UserID     uint64 `json:"userId,omitempty,string"`
	RoleID     uint64 `json:"roleId,omitempty,string"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *UserRole) TableComment() string {
	return "用户角色表"
}
func (ur *UserRole) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ur.ID = j.Get("id").Uint()
	ur.UserID = j.Get("userId").Uint()
	ur.RoleID = j.Get("roleId").Uint()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &Role{}, &UserRole{})
}
