package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectTestCase struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID   uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	RelatedID   uint64         `json:"relatedId,omitempty,string" gorm:"index;comment:关联ID"`
	RelatedType int            `json:"relatedType,omitempty" gorm:"index;comment:关联类型 1-需求 2-任务 3-缺陷问题"`
	Name        string         `json:"name,omitempty" gorm:"size:50;comment:测试名称"`
	Remark      string         `json:"remark,omitempty" gorm:"size:500;comment:备注"`
	CreatorID   uint64         `json:"creatorID,omitempty,string" gorm:"comment:创建人ID"`
	CreateTime  int64          `json:"createTime,omitempty" gorm:"autoCreateTime:milli"`
	UpdateTime  int64          `json:"updateTime,omitempty" gorm:"autoUpdateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTestCase) TableComment() string {
	return "项目测试表"
}

func (pt *ProjectTestCase) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pt.ID = j.Get("id").Uint()
	pt.ProjectID = j.Get("projectId").Uint()
	pt.RelatedID = j.Get("relatedId").Uint()
	pt.RelatedType = int(j.Get("relatedType").Int())
	pt.Name = j.Get("name").String()
	pt.Remark = j.Get("remark").String()
	pt.CreatorID = j.Get("creatorId").Uint()
	pt.CreateTime = j.Get("createTime").Int()
	pt.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTestCase{})
}
