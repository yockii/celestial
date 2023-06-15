import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { ProjectIssue, ProjectIssueCondition } from "@/types/project"

/**
 * 获取项目缺陷列表
 * @param condition - 查询条件
 */
export function getProjectIssueList(condition: ProjectIssueCondition) {
  return request.get<Paginate<ProjectIssue>>("/projectIssue/list", {
    params: condition
  })
}

/**
 * 新增项目缺陷
 * @param projectIssue - 项目缺陷信息
 */
export function addProjectIssue(projectIssue: ProjectIssue) {
  return request.post<boolean>("/projectIssue/add", projectIssue)
}

/**
 * 获取项目缺陷详情
 * @param id - 项目缺陷id
 */
export function getProjectIssue(id: string) {
  return request.get<ProjectIssue>("/projectIssue/instance", {
    params: { id }
  })
}

/**
 * 删除项目缺陷
 * @param id - 项目缺陷id
 */
export function deleteProjectIssue(id: string) {
  return request.delete<boolean>("/projectIssue/delete", {
    params: { id }
  })
}

/**
 * 更新项目缺陷
 * @param projectIssue - 项目缺陷信息
 */
export function updateProjectIssue(projectIssue: ProjectIssue) {
  return request.put<boolean>("/projectIssue/update", projectIssue)
}
