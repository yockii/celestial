import {request} from "../request";
import type {Paginate} from "../request/request"

export type Stage = {
    id: string;
    name: string;
}

/**
 * 获取可筛选阶段列表
 */
export function getStageList() {
    return request.get<Paginate<Stage>>("/stage/list", {
        params: {
            offset: 0,
            limit: 100,
            status: 1,
        }
    })
}