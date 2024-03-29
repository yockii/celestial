import { Condition } from "./common"
import { User } from "./user"

export type AssetCategory = {
  id: string
  parentId?: string
  name: string
  type: number
  creatorId?: string
  createTime?: number
  childrenCount?: number
  children?: AssetCategory[]
  isLeaf?: boolean
}

export type AssetCategoryCondition = Condition & {
  parentId?: string
  name?: string
  type?: number
  onlyParent?: boolean
}

export type CommonTestCase = {
  id: string
  name: string
  remark?: string
  creatorId?: string
  createTime?: number
  items?: CommonTestCaseItem[]
}

export type CommonTestCaseItem = {
  id: string
  testCaseId: string
  content: string
  remark?: string
  creatorId?: string
  createTime?: number
}

export type CommonTestCaseCondition = Condition & {
  id?: string
  name?: string
  categoryId?: string
}

export type File = {
  id: string
  categoryId: string
  ossConfigId?: string
  name: string
  suffix?: string
  size?: number
  objName?: string
  creatorId?: string
  createTime?: number
  creator?: User
  permission?: number
}

export type FileCondition = Condition & {
  id?: string
  categoryId?: string
  name?: string
  suffix?: string
}

export type FileUser = {
  id?: string
  fileId: string
  userId?: string
  permission?: number
  realName?: string
}

export type FileVersion = {
  id: string
  fileId: string
  size: number
  createTime?: number
  creatorId?: string
  creator?: User
}

export type FileVersionCondition = Condition & {
  fileId: string
}
