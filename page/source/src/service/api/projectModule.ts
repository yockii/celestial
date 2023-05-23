import {request} from "../request";
import type {Paginate} from "@/types/common"
import {ProjectModule, ProjectModuleCondition} from "@/types/project";

/**
 * 获取项目模块列表
 * @param condition - 查询条件
 */
export function getProjectModuleList(condition: ProjectModuleCondition) {
    return request.get<Paginate<ProjectModule>>("/projectModule/list", {
        params: condition,
    })
}

/**
 * 新增项目模块
 * @param projectModule - 项目模块信息
 */
export function addProjectModule(projectModule: ProjectModule) {
    return request.post<boolean>("/projectModule/add", projectModule)
}

/**
 * 获取项目模块详情
 * @param id - 项目模块id
 */
export function getProjectModule(id: string) {
    return request.get<ProjectModule>("/projectModule/instance", {
        params: {id},
    })
}

/**
 * 根据projectId获取当前项目执行中模块
 * @param projectId - 项目id
 */
export function getExecutingProjectModuleByProjectId(projectId: string) {
    return request.get<ProjectModule>("/projectModule/executing", {
        params: {projectId},
    })
}

/**
 * 删除项目模块
 * @param id - 项目模块id
 */
export function deleteProjectModule(id: string) {
    return request.delete<boolean>("/projectModule/delete", {
        params: {id},
    })
}

/**
 * 更新项目模块
 * @param projectModule - 项目模块信息
 */
export function updateProjectModule(projectModule: ProjectModule) {
    return request.put<boolean>("/projectModule/update", projectModule)
}