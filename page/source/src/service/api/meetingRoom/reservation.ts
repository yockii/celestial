import { request } from "../../request"

import { MeetingRoomReservation } from "@/types/meetingRoom"

/**
 * 获取会议室预定列表
 * @param meetingRoomId - 会议室id
 */
export function getMeetingRoomReservationList(meetingRoomId: string) {
  return request.get<MeetingRoomReservation[]>("/meetingRoom/reservationList", {
    params: { meetingRoomId }
  })
}

/**
 * 预定会议室
 * @param reservation - 会议室预定信息
 */
export function reserveMeetingRoom(reservation: MeetingRoomReservation) {
  return request.post<MeetingRoomReservation>("/meetingRoom/reserve", reservation)
}

/**
 * 更新会议室预定信息
 * @param reservation - 会议室预定信息
 */
export function updateMeetingRoomReservation(reservation: MeetingRoomReservation) {
  return request.put<boolean>("/meetingRoom/reservationUpdate", reservation)
}

/**
 * 删除会议室预定信息
 * @param reservationId - 会议室预定id
 */
export function deleteMeetingRoomReservation(reservationId: string) {
  return request.delete<boolean>("/meetingRoom/reservationDelete", {
    params: { id: reservationId }
  })
}
