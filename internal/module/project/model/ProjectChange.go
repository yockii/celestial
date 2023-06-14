package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type ProjectChange struct {
	ID             uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID      uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	Title          string         `json:"title,omitempty" gorm:"comment:变化名称"`
	Type           int            `json:"type,omitempty" gorm:"comment:变化类型 1-时间节点调整 2-需求变更 3-资源变动 9-其他变更"`
	Level          int            `json:"level,omitempty" gorm:"comment:变化级别 1-一般 2-重大"`
	Reason         string         `json:"reason,omitempty" gorm:"comment:变更原因"`
	Plan           string         `json:"plan,omitempty" gorm:"comment:变更方案"`
	Review         string         `json:"review,omitempty" gorm:"comment:变更评审结果"`
	Risk           string         `json:"risk,omitempty" gorm:"comment:变更风险"`
	Status         int            `json:"status,omitempty" gorm:"comment:状态 1-待评审 2-已批准 -1-已拒绝 9-关闭"`
	ApplyUserID    uint64         `json:"applyUserId,omitempty,string" gorm:"comment:申请人ID"`
	ReviewerIDList string         `json:"reviewerIdList,omitempty" gorm:"comment:评审人ID列表 以,分割"`
	Result         string         `json:"result,omitempty" gorm:"comment:结果说明"`
	ReviewTime     int64          `json:"reviewTime" gorm:"comment:评审时间"`
	CreateTime     int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime     int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime     gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectChange) TableComment() string {
	return "项目变更表"
}

func (p *ProjectChange) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.ProjectID = j.Get("projectId").Uint()
	p.Title = j.Get("title").String()
	p.Type = int(j.Get("type").Int())
	p.Level = int(j.Get("level").Int())
	p.Reason = j.Get("reason").String()
	p.Plan = j.Get("plan").String()
	p.Review = j.Get("review").String()
	p.Risk = j.Get("risk").String()
	p.Status = int(j.Get("status").Int())
	p.ApplyUserID = j.Get("applyUserId").Uint()
	if j.Get("reviewerIdList").IsArray() {
		for _, idJson := range j.Get("reviewerIdList").Array() {
			p.ReviewerIDList += "," + idJson.String()
		}
		if len(p.ReviewerIDList) > 0 {
			p.ReviewerIDList = p.ReviewerIDList[1:]
		}
	} else {
		p.ReviewerIDList = j.Get("reviewerIdList").String()
	}
	p.Result = j.Get("result").String()
	p.ReviewTime = j.Get("reviewTime").Int()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectChange{})
}
