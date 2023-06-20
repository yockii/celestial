package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

const (
	ThirdSourceStatusEnabled = 1
)

const (
	ThirdSourceCodeDingtalk = "dingtalk"
)

type ThirdSource struct {
	ID            uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name          string `json:"name,omitempty" gorm:"size:50;comment:第三方来源名称"`              // 第三方来源名称
	Code          string `json:"code,omitempty" gorm:"size:50;comment:第三方来源代码"`              // 第三方来源代码
	CorpId        string `json:"corpId,omitempty" gorm:"size:50;comment:第三方企业ID"`            // 第三方企业ID
	Configuration string `json:"configuration,omitempty" gorm:"size:1000;comment:第三方json配置"` // 第三方json配置
	MatchConfig   string `json:"matchConfig,omitempty" gorm:"size:1000;comment:第三方匹配配置"`     // 第三方匹配配置
	Status        int    `json:"status,omitempty" gorm:"comment:状态 1-启用 其他禁用"`
	CreateTime    int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *ThirdSource) TableComment() string {
	return "第三方来源表"
}

func (s *ThirdSource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	s.ID = j.Get("id").Uint()
	s.Name = j.Get("name").String()
	s.Code = j.Get("code").String()
	s.CorpId = j.Get("corpId").String()
	s.Configuration = j.Get("configuration").String()
	s.MatchConfig = j.Get("matchConfig").String()
	s.Status = int(j.Get("status").Int())
	s.CreateTime = j.Get("createTime").Int()
	return nil
}

type ThirdUser struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	SourceID   uint64 `json:"sourceId,omitempty,string" gorm:"comment:第三方来源ID"`
	UserID     uint64 `json:"userId,omitempty,string" gorm:"comment:关联的用户ID"`
	SourceCode string `json:"sourceCode,omitempty" gorm:"size:50;comment:第三方来源代码"` // 第三方来源代码
	OpenID     string `json:"openId,omitempty" gorm:"size:50;comment:第三方openId"`   // 第三方openId
	UnionID    string `json:"unionId,omitempty" gorm:"size:50;comment:第三方unionId"` // 第三方unionId
	Info       string `json:"info,omitempty" gorm:"size:1000;comment:第三方用户json信息"` // 第三方用户信息
	Status     int    `json:"status,omitempty" gorm:"comment:状态 1-正常"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *ThirdUser) TableComment() string {
	return "第三方用户表"
}

func (u *ThirdUser) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	u.ID = j.Get("id").Uint()
	u.SourceID = j.Get("sourceId").Uint()
	u.UserID = j.Get("userId").Uint()
	u.SourceCode = j.Get("sourceCode").String()
	u.OpenID = j.Get("openId").String()
	u.UnionID = j.Get("unionId").String()
	u.Info = j.Get("info").String()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &ThirdSource{}, &ThirdUser{})
}
