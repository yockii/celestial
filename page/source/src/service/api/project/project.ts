import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { Project, ProjectCondition, ProjectStageStatistics, ProjectMember } from "@/types/project"

/**
 * 获取项目列表
 * @param condition - 查询条件
 */
export function getProjectList(condition: ProjectCondition) {
  return request.get<Paginate<Project>>("/project/list", {
    params: condition
  })
}

/**
 * 新增项目
 * @param project - 项目信息
 */
export function addProject(project: Project) {
  return request.post<boolean>("/project/add", project)
}

/**
 * 删除项目
 * @param id - 项目id
 */
export function deleteProject(id: string) {
  return request.delete<boolean>("/project/delete", {
    params: {
      id
    }
  })
}

/**
 * 更新项目信息
 * @param project - 项目信息
 */
export function updateProject(project: Project) {
  return request.put<boolean>("/project/update", project)
}

/**
 * 获取项目详情
 * @param id - 项目id
 */
export function getProjectDetail(id: string) {
  return request.get<Project>(`/project/instance?id=${id}`)
}

/**
 * 获取项目根据阶段统计信息
 */
export function getProjectStageStatistics() {
  return request.get<ProjectStageStatistics[]>("/project/statisticsByStage")
}

export function changeCharger(projectId: string, targetUserId: string) {
  return request.post<boolean>("/project/updateCharger", {
    projectId,
    targetUserId
  })
}

/**
 * 批量添加项目成员
 * @param projectId - 项目id
 * @param roleIdList - 用户id列表
 * @param userIdList - 角色id列表
 */
export function addProjectMembers(projectId: string, roleIdList: string[], userIdList: string[]) {
  return request.post<boolean>("/projectMember/batchAdd", {
    projectId,
    roleIdList,
    userIdList
  })
}

/**
 * 获取项目成员
 * @param projectId - 项目id
 */
export function getProjectMembers(projectId: string) {
  return request.get<ProjectMember[]>("/projectMember/listByProject", {
    params: {
      projectId
    }
  })
}

/**
 * 获取在项目中的权限资源代码
 * @param id - 项目id
 */
export function getProjectResourceCode(id: string) {
  return request.get<string[]>(`/project/resourceCode?id=${id}`)
}

/**
 * 获取我的项目列表
 * @param condition - 查询条件
 */
export function getMyProjectList(condition: ProjectCondition) {
  return request.get<Project[]>("/project/myProjectList", {
    params: condition
  })
}

export function getTopProjects(query: string) {
  return request.get<Project[]>(`/project/topList?name=${query}`)
}

/**
 * 工时统计获取所有项目列表
 */
export function getAllProjectListForWorkTimeStatistics() {
  return request.get<Project[]>("/project/listAllForWorkTimeStatistics")
}
