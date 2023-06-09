package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
	"gorm.io/gorm"
)

const (
	ProjectModuleStatusPendingReview = 1
	ProjectModuleStatusPendingDev    = 2
	ProjectModuleStatusCompleted     = 9
	ProjectModuleStatusRejected      = -1
)

type ProjectModule struct {
	ID            uint64         `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID     uint64         `json:"projectId,omitempty,string" gorm:"index;comment:项目ID"`
	ParentID      uint64         `json:"parentId,omitempty,string" gorm:"index;comment:父模块ID"`
	Name          string         `json:"name,omitempty" gorm:"size:50;comment:模块名称"`
	Alias         string         `json:"alias,omitempty" gorm:"size:50;comment:模块别名"`
	Remark        string         `json:"remark,omitempty" gorm:"size:200;comment:备注"`
	ChildrenCount int            `json:"childrenCount,omitempty" gorm:"comment:子模块数量"`
	FullPath      string         `json:"fullPath,omitempty" gorm:"size:1000;comment:全路径"`
	CreatorID     uint64         `json:"creatorId,omitempty,string" gorm:"comment:创建人ID"`
	Status        int            `json:"status,omitempty" gorm:"comment:状态 1-待评审 2-评审通过待开发 9-已完成 -1-评审不通过"`
	CreateTime    int64          `json:"createTime" gorm:"autoCreateTime:milli"`
	DeleteTime    gorm.DeletedAt `json:"deleteTime" gorm:"index"`
}

func (_ *ProjectModule) TableComment() string {
	return "项目模块表"
}

func (pm *ProjectModule) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	pm.ID = j.Get("id").Uint()
	pm.ProjectID = j.Get("projectId").Uint()
	pm.ParentID = j.Get("parentId").Uint()
	pm.Name = j.Get("name").String()
	pm.Alias = j.Get("alias").String()
	pm.Remark = j.Get("remark").String()
	pm.ChildrenCount = int(j.Get("childrenCount").Int())
	pm.FullPath = j.Get("fullPath").String()
	pm.CreatorID = j.Get("creatorId").Uint()
	pm.CreateTime = j.Get("createTime").Int()
	pm.Status = int(j.Get("status").Int())

	return nil
}

func init() {
	constant.Models = append(constant.Models, &ProjectModule{})
}
