package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

// ProjectTest 项目测试记录
type ProjectTest struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID  uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	Round      int            `json:"round,omitempty" gorm:"comment:测试轮次"`
	TestRecord string         `json:"testRecord,omitempty" gorm:"comment:测试记录"`
	Remark     string         `json:"remark,omitempty" gorm:"size:500;comment:备注"`
	StartTime  int64          `json:"startTime,omitempty" gorm:"comment:开始时间"`
	EndTime    int64          `json:"endTime,omitempty" gorm:"comment:结束时间"`
	CreatorID  uint64         `json:"creatorID,omitempty,string" gorm:"comment:创建人ID"`
	CloserID   uint64         `json:"closerID,omitempty,string" gorm:"comment:封版人ID"`
	CreateTime int64          `json:"createTime,omitempty" gorm:"autoCreateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*ProjectTest) TableComment() string {
	return "项目测试记录表"
}

func (pt *ProjectTest) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pt.ID = j.Get("id").Uint()
	pt.ProjectID = j.Get("projectId").Uint()
	pt.Round = int(j.Get("round").Int())
	pt.TestRecord = j.Get("testRecord").String()
	pt.Remark = j.Get("remark").String()
	pt.StartTime = j.Get("startTime").Int()
	pt.EndTime = j.Get("endTime").Int()
	pt.CreatorID = j.Get("creatorId").Uint()
	pt.CloserID = j.Get("closerId").Uint()
	pt.CreateTime = j.Get("createTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectTest{})
}
