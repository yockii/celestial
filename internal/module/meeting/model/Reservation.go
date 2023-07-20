package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type MeetingRoomReservation struct {
	ID            uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	MeetingRoomID uint64 `json:"meetingRoomId,omitempty,string" gorm:"comment:会议室ID"`
	StartTime     int64  `json:"startTime" gorm:"comment:开始时间"`
	EndTime       int64  `json:"endTime" gorm:"comment:结束时间"`
	Subject       string `json:"subject" gorm:"size:200;comment:会议主题"`
	Participants  string `json:"participants" gorm:"size:2000;comment:参会人员"`
	BookerID      uint64 `json:"bookerId,omitempty,string" gorm:"comment:预订人ID"`
	CreateTime    int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime    int64  `json:"deleteTime,omitempty" gorm:"index"`
}

func (*MeetingRoomReservation) TableComment() string {
	return "会议室预订表"
}

func (r *MeetingRoomReservation) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ID = j.Get("id").Uint()
	r.MeetingRoomID = j.Get("meetingRoomId").Uint()
	r.StartTime = j.Get("startTime").Int()
	r.EndTime = j.Get("endTime").Int()
	r.Subject = j.Get("subject").String()
	r.Participants = j.Get("participants").String()
	r.BookerID = j.Get("bookerId").Uint()
	r.CreateTime = j.Get("createTime").Int()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &MeetingRoomReservation{})
}
