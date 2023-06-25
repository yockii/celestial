import { request } from "../../request"
import { Paginate } from "@/types/common"
import { sm2 } from "sm-crypto"
import { User, UserCondition, UserPermission } from "@/types/user"

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

/**
 * 获取用户权限
 */
export const getUserPermissions = () => {
  return request.get<UserPermission>("/user/permissions")
}

/**
 * 重置用户密码
 * @param id - 用户ID
 * @param password - 密码
 */
export const resetPassword = (id: string, password: string) => {
  const encryptedPassword = "04" + sm2.doEncrypt(password, import.meta.env.VITE_SM2_PK)
  return request.put("/user/resetPassword", { id, password: encryptedPassword })
}
