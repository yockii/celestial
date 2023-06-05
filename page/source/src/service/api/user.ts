import { request } from "../request"
import { Paginate } from "@/types/common"
import { User, UserCondition } from "@/types/user"

/**
 * 用户列表
 * @param condition - 查询条件
 */
export const getUserList = (condition: UserCondition) => {
  return request.get<Paginate<User>>("/user/list", { params: condition })
}

/**
 * 新增用户
 * @param user - 用户信息
 */
export const addUser = (user: User) => {
  return request.post("/user/add", user)
}

/**
 * 更新用户
 * @param user - 用户信息
 */
export const updateUser = (user: User) => {
  return request.put("/user/updateUser", user)
}

/**
 * 删除用户
 * @param id - 用户ID
 */
export const deleteUser = (id: string) => {
  return request.delete("/user/delete", { params: { id } })
}

/**
 * 给用户分配角色
 * @param userId - 用户ID
 * @param roleIdList - 角色ID集合
 */
export const assignRole = (userId: string, roleIdList: string[]) => {
  return request.put("/user/assignRole", { userId, roleIdList })
}

/**
 * 获取用户角色ID列表
 * @param id - 用户ID
 */
export const getUserRoleIdList = (id: string) => {
  return request.get<string[]>("/user/roleIdList", { params: { id } })
}
