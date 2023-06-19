package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectTaskMember struct {
	ID               uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID        uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	TaskID           uint64         `json:"taskId,omitempty,string" gorm:"index;comment:任务ID"`
	UserID           uint64         `json:"userId,omitempty,string" gorm:"index;comment:用户ID"`
	RoleID           uint64         `json:"roleId,omitempty,string" gorm:"index;comment:角色ID"`
	EstimateDuration int64          `json:"estimateDuration,omitempty" gorm:"comment:预计工期,单位:秒"`
	ActualDuration   int64          `json:"actualDuration,omitempty" gorm:"comment:实际工期,单位:秒"`
	Status           int            `json:"status,omitempty" gorm:"comment:任务状态 -1-已取消 1-未开始 2-已确认 3-进行中 9-已完成"`
	CreateTime       int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime       int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime       gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectTaskMember) TableComment() string {
	return "项目任务成员表"
}

func (ptm *ProjectTaskMember) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptm.ID = j.Get("id").Uint()
	ptm.ProjectID = j.Get("projectId").Uint()
	ptm.TaskID = j.Get("taskId").Uint()
	ptm.UserID = j.Get("userId").Uint()
	ptm.RoleID = j.Get("roleId").Uint()
	ptm.EstimateDuration = j.Get("estimateDuration").Int()
	ptm.ActualDuration = j.Get("actualDuration").Int()
	ptm.Status = int(j.Get("status").Int())
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTaskMember{})
}
