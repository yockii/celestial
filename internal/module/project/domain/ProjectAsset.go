package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	"github.com/yockii/ruomu-core/server"
)

type ProjectAssetListRequest struct {
	model.ProjectAsset
	VerifyTimeCondition  *server.TimeCondition `json:"verifyTimeCondition"`
	ReleaseTimeCondition *server.TimeCondition `json:"releaseTimeCondition"`
	ArchiveTimeCondition *server.TimeCondition `json:"archiveTimeCondition"`
	CreateTimeCondition  *server.TimeCondition `json:"createTimeCondition"`
	UpdateTimeCondition  *server.TimeCondition `json:"updateTimeCondition"`
	OrderBy              string                `json:"orderBy"`
}
