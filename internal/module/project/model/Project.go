package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type Project struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ParentID    uint64         `json:"parentId,omitempty,string" gorm:"index;comment:父项目ID"`
	Name        string         `json:"name,omitempty" gorm:"size:50;comment:项目名称"`
	Code        string         `json:"code,omitempty" gorm:"size:50;comment:项目编码"`
	Description string         `json:"description,omitempty" gorm:"size:500;comment:项目描述"`
	OwnerID     uint64         `json:"ownerId,omitempty,string" gorm:"index;comment:项目负责人ID"`
	StageID     uint64         `json:"stageId,omitempty,string" gorm:"index;comment:项目阶段ID"`
	CreateTime  int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime  int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *Project) TableComment() string {
	return "项目表"
}

func (p *Project) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.Name = j.Get("name").String()
	p.Code = j.Get("code").String()
	p.Description = j.Get("description").String()
	p.OwnerID = j.Get("ownerId").Uint()
	p.StageID = j.Get("stageId").Uint()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &Project{})
}
