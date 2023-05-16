package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectRiskListRisk struct {
	model.ProjectRisk
	StartTimeCondition  *server.TimeCondition `json:"startTimeCondition"`
	EndTimeCondition    *server.TimeCondition `json:"endTimeCondition"`
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}

type ProjectRiskCoefficient struct {
	RiskCoefficient float64            `json:"riskCoefficient"`
	MaxRisk         *model.ProjectRisk `json:"maxRisk,omitempty"`
}
