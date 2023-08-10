import { Condition } from "./common"

export type WorkTime = {
  id: string
  userId?: string
  projectId: string
  workTime: number
  startDate: number
  endDate: number
  workContent?: string
  reviewerId?: string
  reviewTime?: number
  status?: number
  rejectReason?: string
  createTime?: number
}

export type WorkTimeCondition = Condition & {
  id?: string
  userId?: string
  projectId?: string
  createTimeCondition?: {
    start?: number
    end?: number
  }
  startDateCondition?: {
    start?: number
    end?: number
  }
  endDateCondition?: {
    start?: number
    end?: number
  }
}
