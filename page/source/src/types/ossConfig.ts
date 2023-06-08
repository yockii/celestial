import { Condition } from "./common"

export type OssConfig = {
  id: string
  type: string
  name: string
  endpoint?: string
  accessKeyId?: string
  secretAccessKey?: string
  bucket?: string
  region?: string
  secure?: number
  selfDomain?: number
  createTime?: number
}

export type OssConfigCondition = Condition & {
  type?: string
  name?: string
  bucket?: string
}
