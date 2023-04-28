package model

import (
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type ProjectAsset struct {
	ID           uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID    uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	Name         string         `json:"name,omitempty" gorm:"size:50;comment:资产名称"`
	Type         int            `json:"type,omitempty" gorm:"comment:资产类型 1-需求 2-设计 3-代码 4-测试 9-其他"`
	Version      string         `json:"version,omitempty" gorm:"size:50;comment:版本号"`
	FileID       uint64         `json:"fileId,omitempty,string" gorm:"index;comment:文件ID"`
	CreatorID    uint64         `json:"creatorId,omitempty,string" gorm:"index;comment:创建人ID"`
	Remark       string         `json:"remark,omitempty" gorm:"size:500;comment:备注"`
	Status       int            `json:"status,omitempty" gorm:"comment:状态 1-草稿 2-已审核 3-已发布 9-已归档"`
	VerifyUserID uint64         `json:"verifyUserId,omitempty,string" gorm:"index;comment:审核人ID"`
	VerifyTime   int64          `json:"verifyTime,omitempty" gorm:"comment:审核时间"`
	ReleaseTime  int64          `json:"releaseTime,omitempty" gorm:"comment:发布时间"`
	ArchiveTime  int64          `json:"archiveTime,omitempty" gorm:"comment:归档时间"`
	CreateTime   int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	UpdateTime   int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
	DeleteTime   gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectAsset) TableComment() string {
	return "项目资产表"
}

func (p *ProjectAsset) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.ProjectID = j.Get("projectId").Uint()
	p.Name = j.Get("name").String()
	p.Type = int(j.Get("type").Int())
	p.Version = j.Get("version").String()
	p.FileID = j.Get("fileId").Uint()
	p.CreatorID = j.Get("creatorId").Uint()
	p.Remark = j.Get("remark").String()
	p.Status = int(j.Get("status").Int())
	p.VerifyUserID = j.Get("verifyUserId").Uint()
	p.VerifyTime = j.Get("verifyTime").Int()
	p.ReleaseTime = j.Get("releaseTime").Int()
	p.ArchiveTime = j.Get("archiveTime").Int()

	return nil
}

func init() {
	Models = append(Models, &ProjectAsset{})
}
