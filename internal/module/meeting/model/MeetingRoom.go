package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	MeetingRoomStatusEnabled  = 1
	MeetingRoomStatusDisabled = -1
)

type MeetingRoom struct {
	ID         uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name       string         `json:"name" gorm:"size:50;comment:会议室名称"`
	Position   string         `json:"position" gorm:"size:200;comment:会议室位置"`
	Capacity   int            `json:"capacity" gorm:"comment:最大容纳人数"`
	Devices    string         `json:"devices" gorm:"size:2000;comment:设备列表"`
	Status     int            `json:"status" gorm:"comment:状态,1:正常,-1:禁用"`
	Remark     string         `json:"remark" gorm:"size:2000;comment:备注"`
	CreatorID  uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建者ID"`
	CreateTime int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*MeetingRoom) TableComment() string {
	return "会议室表"
}

func (mr *MeetingRoom) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	mr.ID = j.Get("id").Uint()
	mr.Name = j.Get("name").String()
	mr.Position = j.Get("position").String()
	mr.Capacity = int(j.Get("capacity").Int())
	mr.Devices = j.Get("devices").String()
	mr.Status = int(j.Get("status").Int())
	mr.Remark = j.Get("remark").String()
	mr.CreatorID = j.Get("creatorId").Uint()
	mr.CreateTime = j.Get("createTime").Int()

	return nil
}

func init() {
	constant.Models = append(constant.Models, &MeetingRoom{})
}
