import { Condition } from "./common"

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
