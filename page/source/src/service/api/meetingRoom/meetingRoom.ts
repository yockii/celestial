import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { MeetingRoom, MeetingRoomCondition } from "@/types/meetingRoom"

/**
 * 新增会议室
 * @param meetingRoom - 会议室信息
 */
export function addMeetingRoom(meetingRoom: MeetingRoom) {
  return request.post<MeetingRoom>("/meetingRoom/add", meetingRoom)
}

/**
 * 修改会议室
 * @param meetingRoom - 会议室信息
 */
export function updateMeetingRoom(meetingRoom: MeetingRoom) {
  return request.put<boolean>("/meetingRoom/update", meetingRoom)
}

/**
 * 删除会议室
 * @param id - 会议室id
 */
export function deleteMeetingRoom(id: string) {
  return request.delete<boolean>("/meetingRoom/delete", {
    params: { id }
  })
}

/**
 * 获取会议室列表
 * @param condition - 查询条件
 */
export function getMeetingRoomList(condition: MeetingRoomCondition) {
  return request.get<Paginate<MeetingRoom>>("/meetingRoom/list", {
    params: condition
  })
}

/**
 * 获取会议室详情
 * @param id - 会议室id
 */
export function getMeetingRoomDetail(id: string) {
  return request.get<MeetingRoom>("/meetingRoom/instance", {
    params: { id }
  })
}
