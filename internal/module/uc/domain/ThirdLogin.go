package domain

import (
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type ThirdSourceListRequest struct {
	model.ThirdSource
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ThirdSourcePublic struct {
	ID     uint64 `json:"id,string"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	AppKey string `json:"appKey"`
	CorpID string `json:"corpId"`
}
