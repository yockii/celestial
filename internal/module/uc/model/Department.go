package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type Department struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ExternalID   string `json:"externalId,omitempty" gorm:"index;size:50;comment:外部ID"`
	Name         string `json:"name,omitempty" gorm:"size:50;comment:部门名称"`
	ParentID     uint64 `json:"parentId,omitempty,string" gorm:"index;comment:父级ID"`
	FullPath     string `json:"fullPath,omitempty" gorm:"size:255;comment:部门路径,/分割"`
	ExternalJson string `json:"externalJson,omitempty" gorm:"size:2000;comment:扩展json信息"`
	OrderNum     int    `json:"orderNum,omitempty" gorm:"comment:排序号"`
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *Department) TableComment() string {
	return "部门表"
}

func (d *Department) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	d.ID = j.Get("id").Uint()
	d.ExternalID = j.Get("externalId").String()
	d.Name = j.Get("name").String()
	d.ParentID = j.Get("parentId").Uint()
	d.FullPath = j.Get("fullPath").String()
	d.OrderNum = int(j.Get("orderNum").Int())
	return nil
}

type UserDepartment struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	UserID       uint64 `json:"userId,omitempty,string" gorm:"index;comment:用户ID"`
	DepartmentID uint64 `json:"departmentId,omitempty,string" gorm:"index;comment:部门ID"`
	ExternalJson string `json:"externalJson,omitempty" gorm:"size:2000;comment:扩展json信息"`
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *UserDepartment) TableComment() string {
	return "用户部门表"
}

func (ud *UserDepartment) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ud.ID = j.Get("id").Uint()
	ud.UserID = j.Get("userId").Uint()
	ud.DepartmentID = j.Get("departmentId").Uint()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &Department{}, &UserDepartment{})
}
