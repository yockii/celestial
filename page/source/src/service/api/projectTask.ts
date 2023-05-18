import {request} from "../request";
import type {ProjectTaskWorkTimeStatistics} from "@/types/project"

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