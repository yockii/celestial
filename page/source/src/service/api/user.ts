import {request} from "../request";
import {Paginate} from "../request/request";

export type User = {
    id: string;
    username: string;
    realName: string;
    password: string;
    email: string;
    mobile: string;
    status: number;
    createTime: number;
}

export type UserCondition = {
    username: string;
    realName: string;
    email: string;
    mobile: string;
    status: number;
    offset: number;
    limit: number;
    orderBy?: string;
}

/**
 * 用户列表
 * @param condition - 查询条件
 */
export const getUserList = (condition: UserCondition) => {
    return request.get<Paginate<User>>('/user/list', {params: condition})
}

/**
 * 新增用户
 * @param user - 用户信息
 */
export const addUser = (user: User) => {
    return request.post('/user/add', user)
}

/**
 * 更新用户
 * @param user - 用户信息
 */
export const updateUser = (user: User) => {
    return request.put('/user/updateUser', user)
}

/**
 * 删除用户
 * @param id - 用户ID
 */
export const deleteUser = (id: string) => {
    return request.delete('/user/delete', {params: {id}})
}