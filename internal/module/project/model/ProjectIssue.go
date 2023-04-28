package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectIssue struct {
	ID           uint64         `json:"ID,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID    uint64         `json:"projectID,omitempty,string" gorm:"index;comment:项目ID"`
	Title        string         `json:"title,omitempty" gorm:"type:varchar(255);comment:问题标题"`
	Content      string         `json:"content,omitempty" gorm:"comment:问题内容"`
	Type         uint8          `json:"type,omitempty" gorm:"comment:问题类型 1-代码错误 2-功能异常 3-界面优化 4-配置相关 5-安全相关 6-性能问题 9-其他问题"`
	Status       uint8          `json:"status,omitempty" gorm:"comment:问题状态 1-新建 2-已指定 3-处理中 4-待验证 5-已解决 9-已关闭"`
	AssignUserID uint64         `json:"assignUserID,omitempty,string" gorm:"index;comment:指派人ID"`
	StartTime    uint64         `json:"startTime,omitempty,string" gorm:"comment:开始解决时间"`
	EndTime      uint64         `json:"endTime,omitempty,string" gorm:"comment:解决完成时间"`
	SolveTime    int64          `json:"solveTime,omitempty" gorm:"comment:解决耗时"`
	IssueCause   string         `json:"issueCause,omitempty" gorm:"comment:问题原因"`
	SolveMethod  string         `json:"solveMethod,omitempty" gorm:"comment:解决方法"`
	CreateTime   uint64         `json:"createTime,omitempty,string" gorm:"comment:创建时间"`
	UpdateTime   uint64         `json:"updateTime,omitempty,string" gorm:"comment:更新时间"`
	DeleteTime   gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectIssue) TableComment() string {
	return "项目问题表"
}

func (pi *ProjectIssue) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pi.ID = j.Get("id").Uint()
	pi.ProjectID = j.Get("project_id").Uint()
	pi.Title = j.Get("title").String()
	pi.Content = j.Get("content").String()
	pi.Type = uint8(j.Get("type").Uint())
	pi.Status = uint8(j.Get("status").Uint())
	pi.AssignUserID = j.Get("assign_user_id").Uint()
	pi.StartTime = j.Get("start_time").Uint()
	pi.EndTime = j.Get("end_time").Uint()
	pi.SolveTime = j.Get("solve_time").Int()
	pi.IssueCause = j.Get("issue_cause").String()
	pi.SolveMethod = j.Get("solve_method").String()

	return nil
}

func init() {
	Models = append(Models, &ProjectIssue{})
}
