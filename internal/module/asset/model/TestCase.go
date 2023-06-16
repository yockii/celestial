package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type CommonTestCase struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name       string `json:"name,omitempty" gorm:"size:50;comment:通用测试用例名称"`
	Remark     string `json:"remark,omitempty" gorm:"size:200;comment:通用测试用例备注"`
	CreatorID  uint64 `json:"creatorId,omitempty,string" gorm:"index;comment:创建人ID"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (*CommonTestCase) TableComment() string {
	return `通用测试用例表`
}

func (ctc *CommonTestCase) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ctc.ID = j.Get("id").Uint()
	ctc.Name = j.Get("name").String()
	ctc.Remark = j.Get("remark").String()
	ctc.CreatorID = j.Get("creatorId").Uint()
	ctc.CreateTime = j.Get("createTime").Int()
	ctc.UpdateTime = j.Get("updateTime").Int()
	return nil
}

type CommonTestCaseItem struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	TestCaseID uint64 `json:"testCaseId,omitempty,string" gorm:"index;comment:通用测试用例ID"`
	Content    string `json:"content,omitempty" gorm:"size:100;comment:测试用例内容"`
	Remark     string `json:"remark,omitempty" gorm:"size:200;comment:通用测试用例备注"`
	CreatorID  uint64 `json:"creatorId,omitempty,string" gorm:"index;comment:创建人ID"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (*CommonTestCaseItem) TableComment() string {
	return `通用测试用例项表`
}

func (ctci *CommonTestCaseItem) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ctci.ID = j.Get("id").Uint()
	ctci.TestCaseID = j.Get("testCaseId").Uint()
	ctci.Content = j.Get("content").String()
	ctci.Remark = j.Get("remark").String()
	ctci.CreatorID = j.Get("creatorId").Uint()
	ctci.CreateTime = j.Get("createTime").Int()
	ctci.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &CommonTestCase{}, &CommonTestCaseItem{})
}
