package domain

import (
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTaskMemberListRequest struct {
	model.ProjectTaskMember
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectTaskMemberWithRealName struct {
	model.ProjectTaskMember
	RealName string `json:"realName"`
}
