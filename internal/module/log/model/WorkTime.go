package model

import "github.com/tidwall/gjson"

type WorkTime struct {
	ID           uint64 `json:"ID,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	TargetID     uint64 `json:"targetID,omitempty,string" gorm:"index;comment:目标ID"`
	TargetType   int    `json:"targetType,omitempty" gorm:"comment:目标类型 1-项目 2-任务 3-测试 9-其他"`
	UserID       uint64 `json:"userID,omitempty,string" gorm:"index;comment:用户ID"`
	WorkTime     int64  `json:"workTime,omitempty" gorm:"comment:工时,单位:秒"`
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
	wt.TargetID = j.Get("targetId").Uint()
	wt.TargetType = int(j.Get("targetType").Int())
	wt.UserID = j.Get("userId").Uint()
	wt.WorkTime = j.Get("workTime").Int()
	wt.WorkContent = j.Get("workContent").String()
	wt.ReviewerID = j.Get("reviewerId").Uint()
	wt.ReviewTime = j.Get("reviewTime").Int()
	wt.Status = int(j.Get("status").Int())
	wt.RejectReason = j.Get("rejectReason").String()

	return nil
}

func init() {
	Models = append(Models, &WorkTime{})
}
