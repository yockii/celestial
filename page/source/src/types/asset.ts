import { Condition } from "./common"

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
