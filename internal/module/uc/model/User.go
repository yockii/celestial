package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

const (
	UserStatusNormal = 1
	UserStatusLeaved = -9
)

type User struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Username   string `json:"username,omitempty" gorm:"size:30;index;comment:用户名"`
	Password   string `json:"password,omitempty" gorm:"size:100;comment:密码"`
	RealName   string `json:"realName,omitempty" gorm:"size:50;comment:真实姓名"`
	Email      string `json:"email,omitempty" gorm:"size:50;comment:邮箱"`
	Mobile     string `json:"mobile,omitempty" gorm:"size:50;comment:手机号"`
	Status     int    `json:"status,omitempty" gorm:"comment:状态 1-正常"`
	ExtType    int    `json:"extType,omitempty" gorm:"comment:扩展类型 -1-无任何扩展类型 1-需要工时填报"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (_ *User) TableComment() string {
	return "用户表"
}

func (u *User) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	u.ID = j.Get("id").Uint()
	u.Username = j.Get("username").String()
	u.Password = j.Get("password").String()
	u.RealName = j.Get("realName").String()
	u.Email = j.Get("email").String()
	u.Mobile = j.Get("mobile").String()
	u.Status = int(j.Get("status").Int())
	u.ExtType = int(j.Get("extType").Int())

	return nil
}

func init() {
	constant.Models = append(constant.Models, &User{})
}
