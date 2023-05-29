import { Paginate } from "@/types/common"
import { request } from "../request"
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
 * 获取项目需求列表
 * @param condition - 查询条件
 */
export function getProjectTaskList(condition: ProjectTaskCondition) {
  return request.get<Paginate<ProjectTask>>("/projectTask/list", {
    params: condition
  })
}

/**
 * 新增项目需求
 * @param projectTask - 项目需求信息
 */
export function addProjectTask(projectTask: ProjectTask) {
  return request.post<boolean>("/projectTask/add", projectTask)
}

/**
 * 获取项目需求详情
 * @param id - 项目需求id
 */
export function getProjectTask(id: string) {
  return request.get<ProjectTask>("/projectTask/instance", {
    params: { id }
  })
}

/**
 * 删除项目需求
 * @param id - 项目需求id
 */
export function deleteProjectTask(id: string) {
  return request.delete<boolean>("/projectTask/delete", {
    params: { id }
  })
}

/**
 * 更新项目需求
 * @param projectTask - 项目需求信息
 */
export function updateProjectTask(projectTask: ProjectTask) {
  return request.put<boolean>("/projectTask/update", projectTask)
}
