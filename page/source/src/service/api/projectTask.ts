import {request} from "../request";
import type {Paginate} from "../request/request"

export type ProjectTaskWorkTimeStatistics = {
    projectId: string;
    taskCount: number;
    estimateDuration: number;
    actualDuration: number;
}
/**
 * 获取项目工时统计
 * @param projectId - 项目id
 */
export function getProjectWorkTimeStatistics(projectId: string) {
    return request.get<ProjectTaskWorkTimeStatistics>("/projectTask/statisticsByProject", {
        params: {
            projectId,
        },
    })
}