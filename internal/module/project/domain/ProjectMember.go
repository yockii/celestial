package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectMemberListRequest struct {
	model.ProjectMember
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectMemberLite struct {
	UserID   uint64 `json:"userId,string"`
	Username string `json:"username"`
	RealName string `json:"realName"`
}
