import { Condition } from "./common"

export type User = {
  id: string
  username: string
  realName?: string
  password?: string
  email?: string
  mobile?: string
  status: number
  createTime?: number
}

export type UserCondition = Condition & {
  username?: string
  realName?: string
  email?: string
  mobile?: string
  status?: number
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
  children?: Resource[]
}

export type ResourceCondition = Condition & {
  resourceName?: string
  resourceCode?: string
}
