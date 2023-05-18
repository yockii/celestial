import {request} from "../request";
import type {Paginate} from "@/types/common"
import {ProjectPlan, ProjectPlanCondition} from "@/types/project";

/**
 * 获取项目计划列表
 * @param condition - 查询条件
 */
export function getProjectPlanList(condition: ProjectPlanCondition) {
    return request.get<Paginate<ProjectPlan>>("/projectPlan/list", {
        params: condition,
    })
}

/**
 * 新增项目计划
 * @param projectPlan - 项目计划信息
 */
export function addProjectPlan(projectPlan: ProjectPlan) {
    return request.post<boolean>("/projectPlan/add", projectPlan)
}

/**
 * 获取项目计划详情
 * @param id - 项目计划id
 */
export function getProjectPlan(id: string) {
    return request.get<ProjectPlan>("/projectPlan/instance", {
        params: {id},
    })
}

/**
 * 根据projectId获取当前项目执行中计划
 * @param projectId - 项目id
 */
export function getExecutingProjectPlanByProjectId(projectId: string) {
    return request.get<ProjectPlan>("/projectPlan/executing", {
        params: {projectId},
    })
}

/**
 * 删除项目计划
 * @param id - 项目计划id
 */
export function deleteProjectPlan(id: string) {
    return request.delete<boolean>("/projectPlan/delete", {
        params: {id},
    })
}

/**
 * 更新项目计划
 * @param projectPlan - 项目计划信息
 */
export function updateProjectPlan(projectPlan: ProjectPlan) {
    return request.put<boolean>("/projectPlan/update", projectPlan)
}