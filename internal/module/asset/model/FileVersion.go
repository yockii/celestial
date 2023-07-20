package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

type FileVersion struct {
	ID          uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	FileID      uint64         `json:"fileId,omitempty,string" gorm:"index;comment:文件ID"`
	OssConfigID uint64         `json:"ossConfigId,omitempty,string" gorm:"comment:OSS配置ID"`
	Size        int64          `json:"size,omitempty" gorm:"comment:文件大小"`
	ObjName     string         `json:"objName,omitempty" gorm:"size:200;comment:存储的对象名称"`
	CreatorID   uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建者ID"`
	CreateTime  int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

func (*FileVersion) TableComment() string {
	return `资产文件版本表`
}

func (af *FileVersion) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	af.ID = j.Get("id").Uint()
	af.FileID = j.Get("fileId").Uint()
	af.OssConfigID = j.Get("ossConfigId").Uint()
	af.Size = j.Get("size").Int()
	af.ObjName = j.Get("objName").String()
	af.CreatorID = j.Get("creatorId").Uint()
	af.CreateTime = j.Get("createTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &FileVersion{})
}
