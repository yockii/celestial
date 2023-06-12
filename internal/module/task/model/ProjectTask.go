package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	ProjectTaskStatusCancel    = -1
	ProjectTaskStatusNotStart  = 1
	ProjectTaskStatusConfirmed = 2
	ProjectTaskStatusDoing     = 3
	ProjectTaskStatusDone      = 9
)

type ProjectTask struct {
	ID               uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID        uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	RequirementID    uint64         `json:"requirementId,omitempty,string" gorm:"index;comment:关联的需求ID"`
	StageID          uint64         `json:"stageId,omitempty,string" gorm:"index;comment:阶段ID"`
	ModuleID         uint64         `json:"moduleId,omitempty,string" gorm:"index;comment:模块ID"`
	ParentID         uint64         `json:"parentId,omitempty,string" gorm:"index;comment:父任务ID"`
	Name             string         `json:"name,omitempty" gorm:"size:50;comment:任务名称"`
	StartTime        int64          `json:"startTime,omitempty" gorm:"comment:开始时间"`
	EndTime          int64          `json:"endTime,omitempty" gorm:"comment:结束时间"`
	TaskDesc         string         `json:"taskDesc,omitempty" gorm:"comment:任务描述"`
	Priority         int            `json:"priority,omitempty" gorm:"comment:优先级 1-低 2-中 3-高"`
	OwnerID          uint64         `json:"ownerId,omitempty,string" gorm:"index;comment:任务负责人ID"`
	ActualStartTime  int64          `json:"actualStartTime,omitempty" gorm:"comment:实际开始时间"`
	ActualEndTime    int64          `json:"actualEndTime,omitempty" gorm:"comment:实际结束时间"`
	EstimateDuration int64          `json:"estimateDuration,omitempty" gorm:"comment:预计工期,单位:秒"`
	ActualDuration   int64          `json:"actualDuration,omitempty" gorm:"comment:实际工期,单位:秒"`
	ChildrenCount    int            `json:"childrenCount,omitempty" gorm:"comment:子任务数量"`
	Status           int            `json:"status,omitempty" gorm:"comment:任务状态 -1-已取消 1-未开始 2-已确认 3-进行中 9-已完成"`
	CreatorID        uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建人ID"`
	FullPath         string         `json:"fullPath,omitempty" gorm:"size:1000;comment:全路径, 需求全路径 + / + 需求名 + / + 任务名"`
	CreateTime       int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime       int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime       gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectTask) TableComment() string {
	return "项目任务表"
}

func (pt *ProjectTask) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pt.ID = j.Get("id").Uint()
	pt.ProjectID = j.Get("projectId").Uint()
	pt.RequirementID = j.Get("requirementId").Uint()
	pt.Priority = int(j.Get("priority").Int())
	pt.StageID = j.Get("stageId").Uint()
	pt.ModuleID = j.Get("moduleId").Uint()
	pt.OwnerID = j.Get("ownerId").Uint()
	pt.ParentID = j.Get("parentId").Uint()
	pt.Name = j.Get("name").String()
	pt.StartTime = j.Get("startTime").Int()
	pt.EndTime = j.Get("endTime").Int()
	pt.TaskDesc = j.Get("taskDesc").String()
	pt.ActualStartTime = j.Get("actualStartTime").Int()
	pt.ActualEndTime = j.Get("actualEndTime").Int()
	pt.EstimateDuration = j.Get("estimateDuration").Int()
	pt.ActualDuration = j.Get("actualDuration").Int()
	pt.ChildrenCount = int(j.Get("childrenCount").Int())
	pt.Status = int(j.Get("status").Int())
	pt.CreatorID = j.Get("creatorId").Uint()
	pt.FullPath = j.Get("fullPath").String()
	pt.CreateTime = j.Get("createTime").Int()
	pt.UpdateTime = j.Get("updateTime").Int()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTask{})
}
