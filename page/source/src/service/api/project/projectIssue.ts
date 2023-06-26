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

/**
 * 分配项目缺陷处理人
 * @param id - 项目缺陷id
 * @param userId - 处理人id
 */
export function assignProjectIssue(id: string, userId: string) {
  return request.put<boolean>("/projectIssue/assign", {
    id,
    assigneeId: userId
  })
}

/**
 * 开始处理项目缺陷
 * @param id - 项目缺陷id
 */
export function startProjectIssue(id: string) {
  return request.put<boolean>("/projectIssue/start", {
    id
  })
}

/**
 * 完成项目缺陷
 * @param id - 项目缺陷id
 */
export function finishProjectIssue(id: string) {
  return request.put<boolean>("/projectIssue/done", {
    id
  })
}

/**
 * 验证项目缺陷
 * @param id - 项目缺陷id
 * @param status - 验证结果
 */
export function verifyProjectIssue(id: string, status: number) {
  return request.put<boolean>("/projectIssue/verify", {
    id,
    status
  })
}

/**
 * 关闭缺陷
 * @param id - 项目缺陷id
 */
export function closeProjectIssue(id: string) {
  return request.put<boolean>("/projectIssue/close", {
    id
  })
}
