package controller

import (
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	// 会议室管理
	{
		meetingRoom := server.Group("/api/v1/meetingRoom")
		meetingRoom.Post("/add", middleware.NeedAuthorization(constant.ResourceMeetingRoomAdd), MeetingRoomController.Add)
		meetingRoom.Delete("/delete", middleware.NeedAuthorization(constant.ResourceMeetingRoomDelete), MeetingRoomController.Delete)
		meetingRoom.Put("/update", middleware.NeedAuthorization(constant.ResourceMeetingRoomUpdate), MeetingRoomController.Update)
		meetingRoom.Get("/list", middleware.NeedAuthorization(constant.ResourceMeetingRoomList), MeetingRoomController.List)
		meetingRoom.Get("/instance", middleware.NeedAuthorization(constant.ResourceMeetingRoomInstance), MeetingRoomController.Instance)
		meetingRoom.Get("/reservationList", middleware.NeedAuthorization(constant.ResourceMeetingRoomInstance), MeetingRoomReservationController.ReservationList)
		meetingRoom.Post("/reserve", middleware.NeedAuthorization(constant.ResourceMeetingRoomReserve), MeetingRoomReservationController.Reserve)
		meetingRoom.Put("/reservationUpdate", middleware.NeedAuthorization(constant.ResourceMeetingRoomReserve), MeetingRoomReservationController.Update)
		meetingRoom.Delete("/reservationDelete", middleware.NeedAuthorization(constant.ResourceMeetingRoomReserve), MeetingRoomReservationController.Delete)

		// 对于禁用put和delete方法时的处理
		meetingRoom.Post("/delete", middleware.NeedAuthorization(constant.ResourceMeetingRoomDelete), MeetingRoomController.Delete)
		meetingRoom.Post("/update", middleware.NeedAuthorization(constant.ResourceMeetingRoomUpdate), MeetingRoomController.Update)
		meetingRoom.Post("/reservationUpdate", middleware.NeedAuthorization(constant.ResourceMeetingRoomReserve), MeetingRoomReservationController.Update)
		meetingRoom.Post("/reservationDelete", middleware.NeedAuthorization(constant.ResourceMeetingRoomReserve), MeetingRoomReservationController.Delete)
	}

}
