import {request} from "../request";
import {Paginate} from "../request/request";

export type Role = {
    id: string;
    name: string;
    desc: string;
    type: number;
    style?: string;
    dataPermission: number;
    status: number;
    createTime?: number;
}

export type RoleCondition = {
    name: string;
    type: number;
    dataPermission: number;
    status: number;
    offset: number;
    limit: number;
    orderBy?: string;
}

/**
 * 新增阶段
 * @param role - 阶段信息
 */
export function addRole(role: Role) {
    return request.post<Role>("role/add", role)
}

/**
 * 修改阶段
 * @param role - 阶段信息
 */
export function updateRole(role: Role) {
    return request.put<boolean>("role/update", role)
}

/**
 * 删除阶段
 * @param id - 阶段id
 */
export function deleteRole(id: string) {
    return request.delete<boolean>("role/delete", {
        params: {id}
    })
}

/**
 * 获取阶段列表
 * @param condition - 查询条件
 */
export function getRoleList(condition: RoleCondition) {
    return request.get<Paginate<Role>>("role/list", {
        params: condition
    })
}