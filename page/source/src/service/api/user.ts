import {request} from "../request";
import {Paginate} from "@/types/common";
import {User, UserCondition} from "@/types/user";

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