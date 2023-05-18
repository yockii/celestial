package domain

import (
	"github.com/yockii/celestial/internal/module/task/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectTaskListTask struct {
	model.ProjectTask
	StartTimeCondition       *server.TimeCondition `json:"startTimeCondition"`
	EndTimeCondition         *server.TimeCondition `json:"endTimeCondition"`
	ActualStartTimeCondition *server.TimeCondition `json:"actualStartTimeCondition"`
	ActualEndTimeCondition   *server.TimeCondition `json:"actualEndTimeCondition"`
	CreateTimeCondition      *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition      *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy                  string                `json:"orderBy"`
}

type ProjectTaskWorkTimeStatistics struct {
	ProjectID uint64 `json:"projectId,string"`
	// 任务总数
	TaskCount int `json:"taskCount"`
	// 任务预计总工时
	EstimateDuration int64 `json:"estimateDuration"`
	// 任务实际总工时
	ActualDuration int64 `json:"actualDuration"`
}
