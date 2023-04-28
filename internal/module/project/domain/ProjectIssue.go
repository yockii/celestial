package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectIssueListRequest struct {
	model.ProjectIssue
	StartTimeCondition  *server.TimeCondition `json:"startTimeCondition"`
	EndTimeCondition    *server.TimeCondition `json:"endTimeCondition"`
	SolveTimeCondition  *server.TimeCondition `json:"solveTimeCondition"`
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
