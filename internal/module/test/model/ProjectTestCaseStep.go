package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectTestCaseStep struct {
	ID         uint64         `json:"ID,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	CaseID     uint64         `json:"caseID,omitempty,string" gorm:"index;comment:测试用例ID"`
	OrderNum   int            `json:"orderNum,omitempty" gorm:"comment:排序号"`
	Content    string         `json:"content,omitempty" gorm:"comment:测试步骤内容"`
	Expect     string         `json:"expect,omitempty" gorm:"comment:预期结果"`
	Status     int            `json:"status,omitempty" gorm:"comment:测试步骤状态 1-未测试 2-已通过 -1-未通过"`
	CreateTime int64          `json:"createTime,omitempty" gorm:"comment:创建时间"`
	UpdateTime int64          `json:"updateTime,omitempty" gorm:"comment:更新时间"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTestCaseStep) TableComment() string {
	return "项目测试用例步骤表"
}

func (ptcs *ProjectTestCaseStep) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptcs.ID = j.Get("ID").Uint()
	ptcs.CaseID = j.Get("caseID").Uint()
	ptcs.OrderNum = int(j.Get("orderNum").Int())
	ptcs.Content = j.Get("content").String()
	ptcs.Expect = j.Get("expect").String()
	ptcs.Status = int(j.Get("status").Int())
	ptcs.CreateTime = j.Get("createTime").Int()
	ptcs.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	Models = append(Models, &ProjectTestCaseStep{})
}
