package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectTest struct {
	ID          uint64         `json:"ID,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID   uint64         `json:"projectID,omitempty,string" gorm:"index;comment:项目ID"`
	RelatedID   uint64         `json:"relatedID,omitempty,string" gorm:"index;comment:关联ID"`
	RelatedType int            `json:"relatedType,omitempty" gorm:"index;comment:关联类型 1-需求 2-任务 3-缺陷问题"`
	Name        string         `json:"name,omitempty" gorm:"size:50;comment:测试名称"`
	Remark      string         `json:"remark,omitempty" gorm:"size:500;comment:备注"`
	CreatorID   uint64         `json:"creatorID,omitempty,string" gorm:"comment:创建人ID"`
	CreateTime  int64          `json:"createTime,omitempty" gorm:"comment:创建时间"`
	UpdateTime  int64          `json:"updateTime,omitempty" gorm:"comment:更新时间"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTest) TableComment() string {
	return "项目测试表"
}

func (pt *ProjectTest) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pt.ID = j.Get("ID").Uint()
	pt.ProjectID = j.Get("projectID").Uint()
	pt.RelatedID = j.Get("relatedID").Uint()
	pt.RelatedType = int(j.Get("relatedType").Int())
	pt.Name = j.Get("name").String()
	pt.Remark = j.Get("remark").String()
	pt.CreatorID = j.Get("creatorID").Uint()
	pt.CreateTime = j.Get("createTime").Int()
	pt.UpdateTime = j.Get("updateTime").Int()
	return nil
}

type ProjectTestCase struct {
	ID         uint64         `json:"ID,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID  uint64         `json:"projectID,omitempty,string" gorm:"index;comment:项目ID"`
	TestID     uint64         `json:"testID,omitempty,string" gorm:"index;comment:测试ID"`
	Name       string         `json:"name,omitempty" gorm:"size:50;comment:测试用例名称"`
	Type       int            `json:"type,omitempty" gorm:"comment:测试用例类型 1-功能测试 2-性能测试 3-安全测试 4-兼容性测试 5-接口测试 9-其他"`
	Content    string         `json:"content,omitempty" gorm:"comment:测试用例内容"`
	Status     int            `json:"status,omitempty" gorm:"comment:测试用例状态 1-未测试 2-已通过 -1-未通过"`
	CreateTime int64          `json:"createTime,omitempty" gorm:"comment:创建时间"`
	UpdateTime int64          `json:"updateTime,omitempty" gorm:"comment:更新时间"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTestCase) TableComment() string {
	return "项目测试用例表"
}

func (ptc *ProjectTestCase) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ptc.ID = j.Get("ID").Uint()
	ptc.ProjectID = j.Get("projectID").Uint()
	ptc.TestID = j.Get("testID").Uint()
	ptc.Name = j.Get("name").String()
	ptc.Type = int(j.Get("type").Int())
	ptc.Content = j.Get("content").String()
	ptc.Status = int(j.Get("status").Int())
	ptc.CreateTime = j.Get("createTime").Int()
	ptc.UpdateTime = j.Get("updateTime").Int()
	return nil
}

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
	Models = append(Models, &ProjectTest{}, &ProjectTestCase{}, &ProjectTestCaseStep{})
}
