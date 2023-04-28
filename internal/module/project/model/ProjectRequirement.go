package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectRequirement struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID   uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	Name        string         `json:"name,omitempty" gorm:"size:50;comment:需求名称"`
	Detail      string         `json:"detail,omitempty" gorm:"comment:需求详情"`
	Priority    int            `json:"priority,omitempty" gorm:"comment:优先级 1-低 2-中 3-高"`
	StageID     uint64         `json:"stageId,omitempty,string" gorm:"index;comment:阶段ID"`
	Source      int            `json:"source,omitempty" gorm:"comment:来源 1-客户 2-内部"`
	OwnerID     uint64         `json:"ownerId,omitempty,string" gorm:"index;comment:需求负责人ID"`
	Feasibility int            `json:"feasibility,omitempty" gorm:"comment:可行性 -1-不可行 0-未评审 1-低 2-中 3-高"`
	CreateTime  int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime  int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectRequirement) TableComment() string {
	return "项目需求表"
}

func (pr *ProjectRequirement) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pr.ID = j.Get("id").Uint()
	pr.ProjectID = j.Get("projectId").Uint()
	pr.Name = j.Get("requestName").String()
	pr.Detail = j.Get("requestDetail").String()
	pr.Priority = int(j.Get("priority").Int())
	pr.StageID = j.Get("stageId").Uint()
	pr.Source = int(j.Get("source").Int())
	pr.OwnerID = j.Get("ownerId").Uint()

	return nil
}

func init() {
	Models = append(Models, &ProjectRequirement{})
}
