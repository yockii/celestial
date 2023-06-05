import { request } from "../request"
import { Paginate } from "@/types/common"
import type { Role, RoleCondition } from "@/types/user"

/**
 * 新增角色
 * @param role - 角色信息
 */
export function addRole(role: Role) {
  return request.post<Role>("role/add", role)
}

/**
 * 修改角色
 * @param role - 角色信息
 */
export function updateRole(role: Role) {
  return request.put<boolean>("role/update", role)
}

/**
 * 删除角色
 * @param id - 角色id
 */
export function deleteRole(id: string) {
  return request.delete<boolean>("role/delete", {
    params: { id }
  })
}

/**
 * 获取角色列表
 * @param condition - 查询条件
 */
export function getRoleList(condition: RoleCondition) {
  return request.get<Paginate<Role>>("role/list", {
    params: condition
  })
}

/**
 * 设置默认角色
 * @param id - 角色id
 */
export function setDefaultRole(id: string) {
  return request.put<boolean>(`role/setDefaultRole?id=${id}`)
}

/**
 * 分配角色资源
 * @param roleId - 角色id
 * @param resourceCodeList - 资源编码列表
 */
export function assignResource(roleId: string, resourceCodeList: string[]) {
  return request.put<boolean>("role/assignResource", {
    roleId,
    resourceCodeList
  })
}

/**
 * 获取角色资源编码列表
 * @param id - 角色id
 */
export function getRoleResourceCodeList(id: string) {
  return request.get<string[]>(`role/resourceCodeList?id=${id}`)
}
