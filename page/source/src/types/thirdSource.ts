import { Condition } from "./common"

export type ThirdSourcePublic = {
  id: string
  name: string
  appKey: string
  corpId: string
}

export type ThirdSource = {
  id: string
  name: string
  code: string
  corpId?: string
  configuration?: string
  matchConfig?: string
  createTime?: number
}

export type ThirdSourceCondition = Condition & {
  name?: string
  code?: string
  corpId?: string
}
