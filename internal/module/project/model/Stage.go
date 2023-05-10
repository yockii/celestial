package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type Stage struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name       string         `json:"name,omitempty" gorm:"size:50;comment:阶段名称"` // 立项阶段、计划阶段、执行阶段、验收阶段、结项阶段...
	OrderNum   int            `json:"orderNum,omitempty" gorm:"comment:排序号"`
	Status     int            `json:"status,omitempty" gorm:"comment:状态 1-正常 2-禁用"`
	CreateTime int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *Stage) TableComment() string {
	return "阶段表"
}

func (s *Stage) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	s.ID = j.Get("id").Uint()
	s.Name = j.Get("name").String()
	s.OrderNum = int(j.Get("orderNum").Int())

	return nil
}

func init() {
	Models = append(Models, &Stage{})
}
