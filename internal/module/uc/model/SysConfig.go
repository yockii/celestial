package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type SysConfig struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Key        string `json:"key" gorm:"column:config_key;size:100;uniqueIndex;comment:配置键"`
	Value      string `json:"value" gorm:"size:500;comment:配置值"`
	Comment    string `json:"comment" gorm:"size:200;comment:说明"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *SysConfig) TableComment() string {
	return "系统配置表"
}

func (c *SysConfig) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	c.ID = j.Get("id").Uint()
	c.Key = j.Get("key").String()
	c.Value = j.Get("value").String()
	c.Comment = j.Get("comment").String()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &SysConfig{})
}
