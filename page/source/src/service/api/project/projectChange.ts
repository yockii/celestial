import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { ProjectChange, ProjectChangeCondition } from "@/types/project"

/**
 * 获取项目变更列表
 * @param condition - 查询条件
 */
export function getProjectChangeList(condition: ProjectChangeCondition) {
  return request.get<Paginate<ProjectChange>>("/projectChange/list", {
    params: condition
  })
}

/**
 * 新增项目变更
 * @param projectChange - 项目变更信息
 */
export function addProjectChange(projectChange: ProjectChange) {
  return request.post<boolean>("/projectChange/add", projectChange)
}

/**
 * 获取项目变更详情
 * @param id - 项目变更id
 */
export function getProjectChange(id: string) {
  return request.get<ProjectChange>("/projectChange/instance", {
    params: { id }
  })
}

/**
 * 删除项目变更
 * @param id - 项目变更id
 */
export function deleteProjectChange(id: string) {
  return request.delete<boolean>("/projectChange/delete", {
    params: { id }
  })
}

/**
 * 更新项目变更
 * @param projectChange - 项目变更信息
 */
export function updateProjectChange(projectChange: ProjectChange) {
  return request.put<boolean>("/projectChange/update", projectChange)
}
