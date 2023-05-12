package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectListRequest struct {
	model.Project
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectWithMembers struct {
	model.Project
	Members []*ProjectMemberLite `json:"members"`
}

type ProjectCountByStage struct {
	StageID uint64 `json:"stageId,string"`
	Count   int64  `json:"count"`
}
