package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	ProjectRiskStatusIdentified = 1
	ProjectRiskStatusResponded  = 2
	ProjectRiskStatusOccurred   = 3
	ProjectRiskStatusSolved     = 4
)

type ProjectRisk struct {
	ID              uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID       uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	StageID         uint64         `json:"stageId,omitempty,string" gorm:"index;comment:阶段ID"`
	RiskName        string         `json:"riskName,omitempty" gorm:"size:50;comment:风险名称"`
	RiskProbability int            `json:"riskProbability,omitempty" gorm:"comment:风险概率 1-低 2-中 3-高"`
	RiskImpact      int            `json:"riskImpact,omitempty" gorm:"comment:风险影响 1-低 2-中 3-高"`
	RiskLevel       int            `json:"riskLevel,omitempty" gorm:"comment:风险等级 1-低 2-中 3-高"`
	Status          int            `json:"status,omitempty" gorm:"comment:风险状态 1-已识别 2-已应对 3-已发生 4-已解决"`
	Response        string         `json:"response,omitempty" gorm:"comment:应对措施"`
	StartTime       int64          `json:"startTime,omitempty" gorm:"comment:开始时间"`
	EndTime         int64          `json:"endTime,omitempty" gorm:"comment:结束时间"`
	Result          string         `json:"result,omitempty" gorm:"comment:应对结果总结"`
	CreateUserID    uint64         `json:"createUserId,omitempty,string" gorm:"index;comment:创建人ID"`
	CreateTime      int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime      int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime      gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectRisk) TableComment() string {
	return "项目风险表"
}

func (pr *ProjectRisk) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pr.ID = j.Get("id").Uint()
	pr.ProjectID = j.Get("projectId").Uint()
	pr.StageID = j.Get("stageId").Uint()
	pr.RiskName = j.Get("riskName").String()
	pr.RiskProbability = int(j.Get("riskProbability").Int())
	pr.RiskImpact = int(j.Get("riskImpact").Int())
	pr.RiskLevel = int(j.Get("riskLevel").Int())
	pr.Status = int(j.Get("status").Int())
	pr.Response = j.Get("response").String()
	pr.StartTime = j.Get("startTime").Int()
	pr.EndTime = j.Get("endTime").Int()
	pr.Result = j.Get("result").String()
	pr.CreateUserID = j.Get("createUserId").Uint()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectRisk{})
}
