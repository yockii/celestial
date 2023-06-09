package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	ProjectRequirementStatusPendingDesign = 1  // 待评审
	ProjectRequirementStatusPendingReview = 2  // 待评审
	ProjectRequirementStatusReviewed      = 3  // 评审通过
	ProjectRequirementStatusCompleted     = 9  // 已完成
	ProjectRequirementStatusRejected      = -1 // 评审未通过
)

type ProjectRequirement struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID   uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	ModuleID    uint64         `json:"moduleId,omitempty,string" gorm:"index;comment:模块ID"`
	Name        string         `json:"name,omitempty" gorm:"size:50;comment:需求名称"`
	Detail      string         `json:"detail,omitempty" gorm:"comment:需求详情"`
	Type        int            `json:"type,omitempty" gorm:"comment:需求类型 1-功能 2-接口 3-性能 4-安全 5-体验 6-改进 7-其他"`
	Priority    int            `json:"priority,omitempty" gorm:"comment:优先级 1-低 2-中 3-高"`
	StageID     uint64         `json:"stageId,omitempty,string" gorm:"index;comment:阶段ID"`
	Source      int            `json:"source,omitempty" gorm:"comment:来源 1-客户 2-内部"`
	OwnerID     uint64         `json:"ownerId,omitempty,string" gorm:"index;comment:需求负责人ID"`
	Feasibility int            `json:"feasibility,omitempty" gorm:"comment:可行性 -1-不可行 1-低 2-中 3-高"`
	Status      int            `json:"status,omitempty" gorm:"comment:状态 1-待设计 2-待评审 3-评审通过 9-已完成 -1-评审未通过"`
	FullPath    string         `json:"fullPath,omitempty" gorm:"size:1000;comment:全路径"`
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
	pr.ModuleID = j.Get("moduleId").Uint()
	pr.Name = j.Get("name").String()
	pr.Detail = j.Get("detail").String()
	pr.Type = int(j.Get("type").Int())
	pr.Priority = int(j.Get("priority").Int())
	pr.StageID = j.Get("stageId").Uint()
	pr.Source = int(j.Get("source").Int())
	pr.OwnerID = j.Get("ownerId").Uint()
	pr.Feasibility = int(j.Get("feasibility").Int())
	pr.Status = int(j.Get("status").Int())
	pr.FullPath = j.Get("fullPath").String()
	pr.CreateTime = j.Get("createTime").Int()
	pr.UpdateTime = j.Get("updateTime").Int()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectRequirement{})
}
