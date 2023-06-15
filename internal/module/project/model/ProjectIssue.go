package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectIssue struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID   uint64         `json:"projectIDd,omitempty,string" gorm:"index;comment:项目ID"`
	Title       string         `json:"title,omitempty" gorm:"type:varchar(255);comment:问题标题"`
	Type        uint8          `json:"type,omitempty" gorm:"comment:问题类型 1-代码错误 2-功能异常 3-界面优化 4-配置相关 5-安全相关 6-性能问题 9-其他问题"`
	Level       uint8          `json:"level,omitempty" gorm:"comment:问题级别 1-一般 2-重要 3-紧急"`
	Content     string         `json:"content,omitempty" gorm:"comment:问题内容"`
	Status      uint8          `json:"status,omitempty" gorm:"comment:问题状态 1-新建 2-已指定 3-处理中 4-待验证 5-已解决 9-已关闭"`
	CreatorID   uint64         `json:"creatorId,omitempty,string" gorm:"index;comment:创建人ID"`
	AssigneeID  uint64         `json:"assigneeId,omitempty,string" gorm:"index;comment:指派人ID"`
	StartTime   uint64         `json:"startTime,omitempty,string" gorm:"comment:开始解决时间"`
	EndTime     uint64         `json:"endTime,omitempty,string" gorm:"comment:解决完成时间"`
	SolveTime   int64          `json:"solveTime,omitempty" gorm:"comment:解决耗时"`
	IssueCause  string         `json:"issueCause,omitempty" gorm:"comment:问题原因"`
	SolveMethod string         `json:"solveMethod,omitempty" gorm:"comment:解决方法"`
	CreateTime  uint64         `json:"createTime,omitempty,string" gorm:"autoCreateTime:milli"`
	UpdateTime  uint64         `json:"updateTime,omitempty,string" gorm:"autoUpdateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectIssue) TableComment() string {
	return "项目问题表"
}

func (pi *ProjectIssue) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pi.ID = j.Get("id").Uint()
	pi.ProjectID = j.Get("projectId").Uint()
	pi.Title = j.Get("title").String()
	pi.Type = uint8(j.Get("type").Uint())
	pi.Level = uint8(j.Get("level").Uint())
	pi.Content = j.Get("content").String()
	pi.Status = uint8(j.Get("status").Uint())
	pi.CreatorID = j.Get("creatorId").Uint()
	pi.AssigneeID = j.Get("assigneeId").Uint()
	pi.StartTime = j.Get("startTime").Uint()
	pi.EndTime = j.Get("endTime").Uint()
	pi.SolveTime = j.Get("solveTime").Int()
	pi.IssueCause = j.Get("issueCause").String()
	pi.SolveMethod = j.Get("solveMethod").String()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectIssue{})
}
