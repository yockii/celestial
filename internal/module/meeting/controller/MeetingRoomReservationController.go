package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/core/helper"
	"github.com/yockii/celestial/internal/module/meeting/model"
	"github.com/yockii/celestial/internal/module/meeting/service"
	"github.com/yockii/ruomu-core/server"
	"time"
)

var MeetingRoomReservationController = new(meetingRoomReservationController)

type meetingRoomReservationController struct{}

func (c *meetingRoomReservationController) ReservationList(ctx *fiber.Ctx) error {
	condition := new(model.MeetingRoomReservation)
	if err := ctx.QueryParser(condition); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if condition.MeetingRoomID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}
	if condition.StartTime == 0 {
		condition.StartTime = time.Now().UnixMilli()
	}
	list, err := service.MeetingRoomReservationService.List(condition)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: list,
	})
}

func (c *meetingRoomReservationController) Reserve(ctx *fiber.Ctx) error {
	reservation := new(model.MeetingRoomReservation)
	if err := ctx.BodyParser(reservation); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if reservation.MeetingRoomID == 0 || reservation.StartTime == 0 || reservation.EndTime == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}
	if reservation.StartTime > reservation.EndTime {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 开始时间不得晚于结束时间",
		})
	}

	if meetingRoom, err := service.MeetingRoomService.Instance(reservation.MeetingRoomID); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	} else if meetingRoom == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 会议室不存在",
		})
	} else if meetingRoom.Status != model.MeetingRoomStatusEnabled {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 会议室已被禁用",
		})
	}

	if uid, err := helper.GetCurrentUserID(ctx); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	} else {
		reservation.BookerID = uid
	}

	if inUse, instance, err := service.MeetingRoomReservationService.Reserve(reservation); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	} else if inUse {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 该时间段已被占用",
		})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: instance,
		})
	}
}

func (c *meetingRoomReservationController) Update(ctx *fiber.Ctx) error {
	reservation := new(model.MeetingRoomReservation)
	if err := ctx.BodyParser(reservation); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if reservation.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}

	old, err := service.MeetingRoomReservationService.Instance(reservation.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}
	if reservation.StartTime == 0 {
		reservation.StartTime = old.StartTime
	}
	if reservation.EndTime == 0 {
		reservation.EndTime = old.EndTime
	}
	if reservation.StartTime > reservation.EndTime {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 开始时间不得晚于结束时间",
		})
	}
	if uid, err := helper.GetCurrentUserID(ctx); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	} else {
		reservation.BookerID = uid
	}

	if inUse, err := service.MeetingRoomReservationService.UpdateReservation(old, reservation); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	} else if inUse {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotMatch,
			Msg:  server.ResponseMsgDataNotMatch + ", 该时间段已被占用",
		})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: true,
		})
	}
}

func (c *meetingRoomReservationController) Delete(ctx *fiber.Ctx) error {
	reservation := new(model.MeetingRoomReservation)
	if err := ctx.QueryParser(reservation); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if reservation.ID == 0 {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough,
		})
	}

	old, err := service.MeetingRoomReservationService.Instance(reservation.ID)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	}
	if old == nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDataNotExists,
			Msg:  server.ResponseMsgDataNotExists,
		})
	}

	if err = service.MeetingRoomReservationService.DeleteReservation(reservation); err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase,
		})
	} else {
		return ctx.JSON(&server.CommonResponse{
			Data: true,
		})
	}
}
