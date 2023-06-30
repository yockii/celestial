import { request } from "../../request"
import { Paginate } from "@/types/common"
import { Department, DepartmentCondition } from "@/types/user"

/**
 * 新增部门
 * @param Department - 部门信息
 */
export function addDepartment(Department: Department) {
  return request.post<Department>("/department/add", Department)
}

/**
 * 修改部门
 * @param Department - 部门信息
 */
export function updateDepartment(Department: Department) {
  return request.put<boolean>("/department/update", Department)
}

/**
 * 删除部门
 * @param id - 部门id
 */
export function deleteDepartment(id: string) {
  return request.delete<boolean>("/department/delete", {
    params: { id }
  })
}

/**
 * 获取部门列表
 * @param condition - 查询条件
 */
export function getDepartmentList(condition: DepartmentCondition) {
  return request.get<Paginate<Department>>("/department/list", {
    params: condition
  })
}

/**
 * 获取部门详情
 * @param id - 部门id
 */
export function getDepartmentDetail(id: string) {
  return request.get<Department>("/department/instance", {
    params: { id }
  })
}

/**
 * 部门添加用户
 * @param departmentId - 部门id
 * @param userIds - 用户id列表
 */
export function addDepartmentUser(departmentId: string, userIds: string[]) {
  return request.post<boolean>("/department/addUser", {
    departmentId,
    userIds
  })
}

/**
 * 部门移除用户
 * @param departmentId - 部门id
 * @param userId - 用户id
 */
export function removeDepartmentUser(departmentId: string, userId: string) {
  return request.delete<boolean>("/department/removeUser", {
    params: {
      departmentId,
      userId
    }
  })
}
