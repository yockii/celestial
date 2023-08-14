import { Condition } from "./common"

export type User = {
  id: string
  username?: string
  realName?: string
  password?: string
  email?: string
  mobile?: string
  status?: number
  extType?: number
  createTime?: number
}

export type UserCondition = Condition & {
  username?: string
  realName?: string
  email?: string
  mobile?: string
  status?: number
  departmentId?: string
}

export type LoginResponse = {
  token: string
  user: User
}

export type Role = {
  id: string
  name: string
  desc: string
  type: number
  style?: string
  dataPermission: number
  defaultRole?: number
  status: number
  createTime?: number
}

export type RoleCondition = Condition & {
  name?: string
  type?: number
  dataPermission?: number
  status?: number
}

export type Resource = {
  id: string
  resourceName: string
  resourceCode: string
  type: number
  children?: Resource[]
}

export type ResourceCondition = Condition & {
  resourceName?: string
  resourceCode?: string
}

export type UserPermission = {
  isSuperAdmin: boolean
  resourceCodeList: string[]
  dataPermission: number
}

export type Department = {
  id: string
  externalId?: string
  name: string
  parentId?: string
  fullPath?: string
  orderNum?: number
  childCount?: number
  createTime?: number
  isLeaf?: boolean
  children?: Department[]
}

export type DepartmentCondition = Condition & {
  name?: string
  parentId?: string
  fullPath?: string
  onlyParent?: boolean
}
