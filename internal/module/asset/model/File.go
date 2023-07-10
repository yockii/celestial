package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	FilePermissionReadOnly = 1
	FilePermissionEditable = 2
	FilePermissionDownload = 3
	FilePermissionManage   = 4
)

type File struct {
	ID           uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	CategoryID   uint64         `json:"categoryId,omitempty,string" gorm:"index;comment:分类ID"`
	OssConfigID  uint64         `json:"ossConfigId,omitempty,string" gorm:"comment:OSS配置ID"`
	OriginName   string         `json:"originName,omitempty" gorm:"size:100;comment:原始文件名"`
	Name         string         `json:"name,omitempty" gorm:"size:100;comment:文件名"`
	Suffix       string         `json:"suffix,omitempty" gorm:"size:20;comment:文件后缀"`
	Size         int64          `json:"size,omitempty" gorm:"comment:文件大小"`
	ObjName      string         `json:"objName,omitempty" gorm:"size:200;comment:存储的对象名称"`
	CategoryPath string         `json:"categoryPath,omitempty" gorm:"size:200;comment:分类路径"`
	CreatorID    uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建者ID"`
	CreateTime   int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime   gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
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
	af.CategoryPath = j.Get("categoryPath").String()
	af.CreatorID = j.Get("creatorId").Uint()
	af.CreateTime = j.Get("createTime").Int()
	af.UpdateTime = j.Get("updateTime").Int()
	return nil
}

// FilePermission 文件权限
type FilePermission struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	FileID     uint64 `json:"assetId,omitempty,string" gorm:"index;comment:资产ID"`
	UserID     uint64 `json:"userId,omitempty,string" gorm:"index;comment:用户ID"`
	Permission uint8  `json:"permission,omitempty" gorm:"comment:权限 1-可读(在线预览) 2-可编辑 3-可下载 4-可管理(编辑/下载/删除)"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime:milli"`
}

func (*FilePermission) TableComment() string {
	return `资产文件表`
}

func (af *FilePermission) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	af.ID = j.Get("id").Uint()
	af.FileID = j.Get("fileId").Uint()
	af.UserID = j.Get("userId").Uint()
	af.Permission = uint8(j.Get("permission").Uint())
	af.CreateTime = j.Get("createTime").Int()
	af.UpdateTime = j.Get("updateTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &File{}, &FilePermission{})
}
