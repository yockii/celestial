import { WorkTime, WorkTimeCondition, WorkTimeStatisticsCondition } from "@/types/log"
import { request } from "../../request"
import { Paginate } from "@/types/common"

/**
 * 新增工时
 * @param data
 */
export function addWorkTime(data: WorkTime) {
  return request.post<WorkTime>("/workTime/add", data)
}

/**
 * 删除工时
 * @param id
 */
export function deleteWorkTime(id: string) {
  return request.delete<boolean>("/workTime/delete", {
    params: { id }
  })
}

/**
 * 修改工时
 * @param data
 */
export function updateWorkTime(data: WorkTime) {
  return request.put<boolean>("/workTime/update", data)
}

/**
 * 获取工时列表
 * @param condition
 */
export function getWorkTimeList(condition: WorkTimeCondition) {
  return request.get<Paginate<WorkTime>>("/workTime/list", {
    params: condition
  })
}

/**
 * 获取工时统计信息
 * @param condition
 */
export function getWorkTimeStatistics(condition: WorkTimeStatisticsCondition) {
  return request.get<WorkTime[]>("/workTime/statistics", {
    params: condition
  })
}
