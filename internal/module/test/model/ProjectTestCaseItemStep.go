package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectTestCaseItemStep struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	CaseItemID uint64         `json:"caseItemId,omitempty,string" gorm:"index;comment:测试用例ID"`
	OrderNum   int            `json:"orderNum,omitempty" gorm:"comment:排序号"`
	Content    string         `json:"content,omitempty" gorm:"size:500;comment:测试步骤内容"`
	Expect     string         `json:"expect,omitempty" gorm:"size:500;comment:预期结果"`
	Status     int            `json:"status,omitempty" gorm:"comment:测试步骤状态 1-未测试 2-已通过 -1-未通过"`
	CreateTime int64          `json:"createTime,omitempty" gorm:"autoCreateTime:milli"`
	UpdateTime int64          `json:"updateTime,omitempty" gorm:"autoUpdateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTestCaseItemStep) TableComment() string {
	return "项目测试用例步骤表"
}

func (ptcs *ProjectTestCaseItemStep) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptcs.ID = j.Get("id").Uint()
	ptcs.CaseItemID = j.Get("caseItemId").Uint()
	ptcs.OrderNum = int(j.Get("orderNum").Int())
	ptcs.Content = j.Get("content").String()
	ptcs.Expect = j.Get("expect").String()
	ptcs.Status = int(j.Get("status").Int())
	ptcs.CreateTime = j.Get("createTime").Int()
	ptcs.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTestCaseItemStep{})
}
