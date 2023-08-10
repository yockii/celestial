package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type WorkTime struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	UserID       uint64 `json:"userId,omitempty,string" gorm:"index;comment:用户ID"`
	ProjectID    uint64 `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	WorkTime     int64  `json:"workTime,omitempty" gorm:"comment:工时,单位:秒"`
	StartDate    int64  `json:"startDate,omitempty" gorm:"comment:开始日期"`
	EndDate      int64  `json:"endDate,omitempty" gorm:"comment:结束日期"`
	WorkContent  string `json:"workContent,omitempty" gorm:"size:500;comment:工作内容"`
	ReviewerID   uint64 `json:"reviewerID,omitempty,string" gorm:"index;comment:审核人ID"`
	ReviewTime   int64  `json:"reviewTime,omitempty" gorm:"comment:审核时间"`
	Status       int    `json:"status,omitempty" gorm:"comment:状态 0-未提交 1-已提交 2-已审核 3-已驳回 4-已取消"`
	RejectReason string `json:"rejectReason,omitempty" gorm:"size:500;comment:驳回原因"`
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *WorkTime) TableComment() string {
	return "工时表"
}

func (wt *WorkTime) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	wt.ID = j.Get("id").Uint()
	wt.UserID = j.Get("userId").Uint()
	wt.ProjectID = j.Get("projectId").Uint()
	wt.WorkTime = j.Get("workTime").Int()
	wt.StartDate = j.Get("startDate").Int()
	wt.EndDate = j.Get("endDate").Int()
	wt.WorkContent = j.Get("workContent").String()
	wt.ReviewerID = j.Get("reviewerId").Uint()
	wt.ReviewTime = j.Get("reviewTime").Int()
	wt.Status = int(j.Get("status").Int())
	wt.RejectReason = j.Get("rejectReason").String()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &WorkTime{})
}
