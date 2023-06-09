package domain

import (
	"github.com/yockii/celestial/internal/module/project/model"
	ucModel "github.com/yockii/celestial/internal/module/uc/model"
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

type ProjectAsset struct {
	model.ProjectAsset
	Permission uint8         `json:"permission"`
	Creator    *ucModel.User `json:"creator"`
}
