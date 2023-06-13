package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectTestCaseItem struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID  uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	TestCaseID uint64         `json:"testCaseId,omitempty,string" gorm:"index;comment:测试ID"`
	Name       string         `json:"name,omitempty" gorm:"size:50;comment:测试用例名称"`
	Type       int            `json:"type,omitempty" gorm:"comment:测试用例类型 1-功能测试 2-性能测试 3-安全测试 4-兼容性测试 5-接口测试 9-其他"`
	Content    string         `json:"content,omitempty" gorm:"comment:测试用例内容"`
	Status     int            `json:"status,omitempty" gorm:"comment:测试用例状态 1-未测试 2-已通过 -1-未通过"`
	CreateTime int64          `json:"createTime,omitempty" gorm:"autoCreateTime:milli"`
	UpdateTime int64          `json:"updateTime,omitempty" gorm:"autoUpdateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTestCaseItem) TableComment() string {
	return "项目测试用例表"
}

func (ptc *ProjectTestCaseItem) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptc.ID = j.Get("id").Uint()
	ptc.ProjectID = j.Get("projectId").Uint()
	ptc.TestCaseID = j.Get("testCaseId").Uint()
	ptc.Name = j.Get("name").String()
	ptc.Type = int(j.Get("type").Int())
	ptc.Content = j.Get("content").String()
	ptc.Status = int(j.Get("status").Int())
	ptc.CreateTime = j.Get("createTime").Int()
	ptc.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTestCaseItem{})
}
