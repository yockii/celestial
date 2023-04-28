package domain

import (
	"github.com/yockii/celestial/internal/module/uc/model"
	"github.com/yockii/ruomu-core/server"
)

type RoleListRequest struct {
	model.Role
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type RoleDispatchResourcesRequest struct {
	RoleID           uint64   `json:"roleId,string"`
	ResourceCodeList []string `json:"resourceCodeList"`
}
