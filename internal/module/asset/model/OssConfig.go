package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/celestial/internal/constant"
)

type OssConfig struct {
	ID              uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Type            string `json:"type,omitempty" gorm:"size:50;comment:类型 oss minio ks3 obs azure"`
	Name            string `json:"name,omitempty" gorm:"index;size:50;comment:名称"`
	Endpoint        string `json:"endpoint,omitempty" gorm:"size:200;comment:Endpoint"`
	AccessKeyID     string `json:"accessKeyId,omitempty" gorm:"size:100;comment:AccessKeyID"`
	SecretAccessKey string `json:"secretAccessKey,omitempty" gorm:"size:100;comment:secretAccessKey"`
	Bucket          string `json:"bucket,omitempty" gorm:"size:50;comment:桶"`
	Region          string `json:"region,omitempty" gorm:"size:50;comment:Region"`
	Secure          int    `json:"secure,omitempty" gorm:"comment:是否使用HTTPS 1-是 2-否"`
	SelfDomain      int    `json:"selfDomain,omitempty" gorm:"comment:是否自定义域名 1-是 2-否"`
	Status          int    `json:"status,omitempty" gorm:"comment:状态 1-启用 其他-禁用"`
	CreateTime      int64  `json:"createTime" gorm:"autoCreateTime:milli"`
}

func (_ *OssConfig) TableComment() string {
	return "对象存储配置表"
}

func (p *OssConfig) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.Name = j.Get("name").String()
	p.Type = j.Get("type").String()
	p.Endpoint = j.Get("endpoint").String()
	p.AccessKeyID = j.Get("accessKeyId").String()
	p.SecretAccessKey = j.Get("secretAccessKey").String()
	p.Bucket = j.Get("bucket").String()
	p.Region = j.Get("region").String()
	p.Secure = int(j.Get("secure").Int())
	p.SelfDomain = int(j.Get("selfDomain").Int())
	p.Status = int(j.Get("status").Int())
	p.CreateTime = j.Get("createTime").Int()
	return nil
}

func init() {
	constant.Models = append(constant.Models, &OssConfig{})
}
