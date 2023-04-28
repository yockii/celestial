package model

import "github.com/tidwall/gjson"

type AssetCategory struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ParentID   uint64 `json:"parentId,omitempty,string" gorm:"index;comment:父级ID"`
	Name       string `json:"name,omitempty" gorm:"index;size:50;comment:名称"`
	Type       int    `json:"type,omitempty" gorm:"comment:类型 1-公共资产 2-项目资产 9-个人资产"`
	CreatorID  uint64 `json:"creatorId,omitempty,string" gorm:"index;comment:创建人ID"`
	CreateTime int64  `json:"createTime,omitempty" gorm:"comment:创建时间"`
}

func (*AssetCategory) TableComment() string {
	return `资产分类表`
}

func (ac *AssetCategory) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ac.ID = j.Get("id").Uint()
	ac.ParentID = j.Get("parentId").Uint()
	ac.Name = j.Get("name").String()
	ac.Type = int(j.Get("type").Int())
	ac.CreatorID = j.Get("creatorId").Uint()
	ac.CreateTime = j.Get("createTime").Int()
	return nil
}

func init() {
	Models = append(Models, &AssetCategory{})
}
