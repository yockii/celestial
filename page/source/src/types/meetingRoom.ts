import { Condition } from "./common"

export type MeetingRoom = {
  id: string
  name?: string
  position?: string
  capacity?: number
  devices?: string
  remark?: string
  status: number
  creatorId?: string
  createTime?: number
}

export type MeetingRoomCondition = Condition & {
  name?: string
  position?: string
  devices?: string
  status?: number
}

export type MeetingRoomReservation = {
  id: string
  meetingRoomId: string
  startTime: number
  endTime: number
  subject: string
  participants?: string
  bookerId?: string
  createTime?: number
}
