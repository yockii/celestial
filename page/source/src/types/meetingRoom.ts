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
