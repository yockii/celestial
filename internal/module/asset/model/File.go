package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type File struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	CategoryID  uint64         `json:"categoryId,omitempty,string" gorm:"index;comment:分类ID"`
	OssConfigID uint64         `json:"ossConfigId,omitempty,string" gorm:"comment:OSS配置ID"`
	OriginName  string         `json:"originName,omitempty" gorm:"comment:原始文件名"`
	Name        string         `json:"name,omitempty" gorm:"comment:文件名"`
	Suffix      string         `json:"suffix,omitempty" gorm:"comment:文件后缀"`
	Size        int64          `json:"size,omitempty" gorm:"comment:文件大小"`
	ObjName     string         `json:"objName,omitempty" gorm:"comment:存储的对象名称"`
	CreatorID   uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建者ID"`
	CreateTime  int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*File) TableComment() string {
	return `资产文件表`
}

func (af *File) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	af.ID = j.Get("id").Uint()
	af.CategoryID = j.Get("categoryId").Uint()
	af.OssConfigID = j.Get("ossConfigId").Uint()
	af.OriginName = j.Get("originName").String()
	af.Name = j.Get("name").String()
	af.Suffix = j.Get("suffix").String()
	af.Size = j.Get("size").Int()
	af.ObjName = j.Get("objName").String()
	af.CreatorID = j.Get("creatorId").Uint()
	af.CreateTime = j.Get("createTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &File{})
}
