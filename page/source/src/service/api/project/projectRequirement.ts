import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { ProjectRequirement, ProjectRequirementCondition } from "@/types/project"

/**
 * 获取项目需求列表
 * @param condition - 查询条件
 */
export function getProjectRequirementList(condition: ProjectRequirementCondition) {
  return request.get<Paginate<ProjectRequirement>>("/projectRequirement/list", {
    params: condition
  })
}

/**
 * 新增项目需求
 * @param projectRequirement - 项目需求信息
 */
export function addProjectRequirement(projectRequirement: ProjectRequirement) {
  return request.post<boolean>("/projectRequirement/add", projectRequirement)
}

/**
 * 获取项目需求详情
 * @param id - 项目需求id
 */
export function getProjectRequirement(id: string) {
  return request.get<ProjectRequirement>("/projectRequirement/instance", {
    params: { id }
  })
}

/**
 * 删除项目需求
 * @param id - 项目需求id
 */
export function deleteProjectRequirement(id: string) {
  return request.delete<boolean>("/projectRequirement/delete", {
    params: { id }
  })
}

/**
 * 更新项目需求
 * @param projectRequirement - 项目需求信息
 */
export function updateProjectRequirement(projectRequirement: ProjectRequirement) {
  return request.put<boolean>("/projectRequirement/update", projectRequirement)
}

/**
 * 需求设计完毕
 * @param id - 项目需求id
 */
export function requirementDesigned(id: string) {
  return request.put<boolean>("/projectRequirement/designed", {
    id
  })
}

/**
 * 需求评审
 * @param id - 项目需求id
 * @param status - 评审状态
 */
export function requirementReview(id: string, status: number) {
  return request.put<boolean>("/projectRequirement/review", {
    id,
    status
  })
}

/**
 * 需求完成
 * @param id - 项目需求id
 */
export function requirementCompleted(id: string) {
  return request.put<boolean>("/projectRequirement/completed", {
    id
  })
}
