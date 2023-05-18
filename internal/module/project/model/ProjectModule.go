package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectModule struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID  uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	ParentID   uint64         `json:"parentId,omitempty,string" gorm:"index;comment:父模块ID"`
	Name       string         `json:"name,omitempty" gorm:"size:50;comment:模块名称"`
	Alias      string         `json:"alias,omitempty" gorm:"size:50;comment:模块别名"`
	Remark     string         `json:"remark,omitempty" gorm:"size:200;comment:备注"`
	CreatorID  uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建人ID"`
	CreateTime int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectModule) TableComment() string {
	return "项目模块表"
}

func (pm *ProjectModule) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pm.ID = j.Get("id").Uint()
	pm.ProjectID = j.Get("projectId").Uint()
	pm.ParentID = j.Get("parentId").Uint()
	pm.Name = j.Get("name").String()
	pm.Alias = j.Get("alias").String()
	pm.Remark = j.Get("remark").String()
	pm.CreatorID = j.Get("creatorId").Uint()
	pm.CreateTime = j.Get("createTime").Int()

	return nil
}
