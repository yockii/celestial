package model

import (
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectPlan struct {
	ID           uint64          `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID    uint64          `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	StageID      uint64          `json:"stageId,omitempty,string" gorm:"index;comment:阶段ID"`
	PlanName     string          `json:"planName,omitempty" gorm:"size:50;comment:计划名称"`
	PlanDesc     string          `json:"planDesc,omitempty" gorm:"size:500;comment:计划描述"`
	StartTime    int64           `json:"startTime" gorm:"comment:计划开始时间"`
	EndTime      int64           `json:"endTime" gorm:"comment:计划结束时间"`
	Target       string          `json:"target,omitempty" gorm:"size:500;comment:计划目标"`
	Scope        string          `json:"scope,omitempty" gorm:"size:500;comment:计划范围"`
	Schedule     string          `json:"schedule,omitempty" gorm:"size:500;comment:计划进度"`
	Resource     string          `json:"resource,omitempty" gorm:"size:500;comment:计划资源"`
	Budget       decimal.Decimal `json:"budget,omitempty" gorm:"type:decimal(10,2);comment:计划预算"`
	CreateUserID uint64          `json:"createUserId,omitempty,string" gorm:"index;comment:创建人ID"`
	Status       int             `json:"status,omitempty" gorm:"comment:状态 -1-已废弃 1-未开始 2-进行中 3-已完成"`
	CreateTime   int64           `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64           `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime   gorm.DeletedAt  `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectPlan) TableComment() string {
	return "项目计划表"
}

func (pp *ProjectPlan) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pp.ID = j.Get("id").Uint()
	pp.ProjectID = j.Get("projectId").Uint()
	pp.StageID = j.Get("stageId").Uint()
	pp.PlanName = j.Get("planName").String()
	pp.PlanDesc = j.Get("planDesc").String()
	pp.StartTime = j.Get("startTime").Int()
	pp.EndTime = j.Get("endTime").Int()
	pp.Target = j.Get("target").String()
	pp.Scope = j.Get("scope").String()
	pp.Schedule = j.Get("schedule").String()
	pp.Resource = j.Get("resource").String()
	pp.Budget = decimal.NewFromFloat(j.Get("budget").Float())
	pp.CreateUserID = j.Get("createUserId").Uint()
	pp.Status = int(j.Get("status").Int())

	return nil
}

func init() {
	Models = append(Models, &ProjectPlan{})
}
