package domain

import (
	"github.com/yockii/celestial/internal/module/meeting/model"
	"github.com/yockii/ruomu-core/server"
)

type MeetingRoomListRequest struct {
	model.MeetingRoom
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition"`
	OrderBy             string                `json:"orderBy"`
}
