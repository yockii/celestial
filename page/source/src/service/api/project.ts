import {request} from "../request";
import type {Paginate} from "../request/request"

export type Project = {
    id: string;
    name: string;
    stageId: string;
}
export type ProjectCondition = {
    name?: string;
    stageId?: string;
    offset?: number;
    limit?: number;
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