import { Condition } from "./common"

export type WorkTime = {
  id: string
  userId?: string
  name?: string
  projectId: string
  topParentId?: string
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
    start?: string
    end?: string
  }
  startDateCondition?: {
    start?: string
    end?: string
  }
  endDateCondition?: {
    start?: string
    end?: string
  }
}

export type WorkTimeStatisticsCondition = {
  departmentId?: string
  dateCondition?: {
    start?: number
    end?: number
  }
}
