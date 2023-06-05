package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type Log struct {
	ID            uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	TargetID      uint64 `json:"targetId,omitempty,string" gorm:"index;comment:目标ID"`
	TargetType    int    `json:"targetType,omitempty" gorm:"comment:目标类型 1-项目 2-需求 3-设计 4-代码 5-测试 6-缺陷 7-任务 8-文档 9-登录登出 99-其他"`
	Action        string `json:"action,omitempty" gorm:"size:50;comment:操作"`
	Content       string `json:"content,omitempty" gorm:"comment:原内容记录"`
	ChangeContent string `json:"changeContent,omitempty" gorm:"comment:变更内容记录"`
	Status        int    `json:"status,omitempty" gorm:"comment:状态 1-成功 -1-失败"`
	OperatorID    uint64 `json:"operatorId,omitempty,string" gorm:"index;comment:操作人ID"`
	CreateTime    int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *Log) TableComment() string {
	return "日志表"
}

func (p *Log) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.TargetID = j.Get("targetId").Uint()
	p.TargetType = int(j.Get("targetType").Int())
	p.Action = j.Get("action").String()
	p.Content = j.Get("content").String()
	p.ChangeContent = j.Get("changeContent").String()
	p.Status = int(j.Get("status").Int())
	p.OperatorID = j.Get("operatorId").Uint()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &Log{})
}
