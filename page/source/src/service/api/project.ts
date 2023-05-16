import {request} from "../request";
import type {Paginate} from "../request/request"

export type ProjectMember = {
    userId:string;
    username: string;
    realName: string;
    roleId: string;
}

export type Project = {
    id: string;
    name: string;
    code: string;
    description: string;
    stageId: string;
    createTime?: number;
    members?: ProjectMember[];
}
export type ProjectCondition = {
    name: string;
    stageId: string;
    offset?: number;
    limit?: number;
}

export type ProjectStageStatistics = {
    stageId: string;
    count: number;
}

/**
 * 获取项目列表
 * @param condition - 查询条件
 */
export function getProjectList(condition: ProjectCondition) {
    return request.get<Paginate<Project>>("/project/list", {
        params: condition,
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