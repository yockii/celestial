import {request} from "../request";
import {Paginate} from "@/types/common";
import {Stage, StageCondition} from "@/types/project"

/**
 * 新增阶段
 * @param stage - 阶段信息
 */
export function addStage(stage: Stage) {
    return request.post<Stage>("stage/add", stage)
}

/**
 * 修改阶段
 * @param stage - 阶段信息
 */
export function updateStage(stage: Stage) {
    return request.put<boolean>("stage/update", stage)
}

/**
 * 删除阶段
 * @param id - 阶段id
 */
export function deleteStage(id: string) {
    return request.delete<boolean>("stage/delete", {
        params: {id}
    })
}

/**
 * 获取阶段列表
 * @param condition - 查询条件
 */
export function getStageList(condition: StageCondition) {
    return request.get<Paginate<Stage>>("stage/list", {
        params: condition
    })
}

/**
 * 获取阶段详情
 * @param id - 阶段id
 */
export function getStageDetail(id: string) {
    return request.get<Stage>("stage/instance", {
        params: {id}
    })
}