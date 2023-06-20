import { Paginate } from "@/types/common"
import { request } from "../../request"
import type { ProjectTaskWorkTimeStatistics, ProjectTask, ProjectTaskCondition } from "@/types/project"

/**
 * 获取项目工时统计
 * @param projectId - 项目id
 */
export function getProjectWorkTimeStatistics(projectId: string) {
  return request.get<ProjectTaskWorkTimeStatistics>("/projectTask/statisticsByProject", {
    params: {
      projectId
    }
  })
}

/**
 * 获取项目任务列表
 * @param condition - 查询条件
 */
export function getProjectTaskList(condition: ProjectTaskCondition) {
  return request.get<Paginate<ProjectTask>>("/projectTask/list", {
    params: condition
  })
}

/**
 * 新增项目任务
 * @param projectTask - 项目任务信息
 */
export function addProjectTask(projectTask: ProjectTask) {
  return request.post<boolean>("/projectTask/add", projectTask)
}

/**
 * 获取项目任务详情
 * @param id - 项目任务id
 */
export function getProjectTask(id: string) {
  return request.get<ProjectTask>("/projectTask/instance", {
    params: { id }
  })
}

/**
 * 删除项目任务
 * @param id - 项目任务id
 */
export function deleteProjectTask(id: string) {
  return request.delete<boolean>("/projectTask/delete", {
    params: { id }
  })
}

/**
 * 更新项目任务
 * @param projectTask - 项目任务信息
 */
export function updateProjectTask(projectTask: ProjectTask) {
  return request.put<boolean>("/projectTask/update", projectTask)
}

/**
 * 取消任务
 * @param id - 项目任务id
 */
export function cancelProjectTask(id: string) {
  return request.put<boolean>("/projectTask/cancel", {
    id
  })
}

/**
 * 确认任务
 * @param id - 项目任务id
 * @param estimateDuration - 预计工时
 */
export function confirmProjectTask(id: string, estimateDuration: number) {
  return request.put<boolean>("/projectTask/confirm", {
    id,
    estimateDuration
  })
}

/**
 * 开始任务
 * @param id - 项目任务id
 */
export function startProjectTask(id: string) {
  return request.put<boolean>("/projectTask/start", {
    id
  })
}

/**
 * 完成任务
 * @param id - 项目任务id
 * @param actualDuration - 实际工时
 */
export function finishProjectTask(id: string, actualDuration: number) {
  return request.put<boolean>("/projectTask/done", {
    id,
    actualDuration
  })
}

/**
 * 重新开始任务
 * @param id - 项目任务id
 */
export function restartProjectTask(id: string) {
  return request.put<boolean>("/projectTask/restart", {
    id
  })
}
