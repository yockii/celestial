import { request } from "../../request"

export type SysConfig = {
  id?: string
  key: string
  value: string
  comment?: string
  updateTime?: number
}

/**
 * 系统配置列表
 */
export const getSysConfigList = () => {
  return request.get<SysConfig[]>("/sysConfig/list")
}

/**
 * 更新系统配置
 * @param key - 配置键
 * @param value - 配置值
 */
export const updateSysConfig = (key: string, value: string) => {
  return request.put<boolean>("/sysConfig/update", { key, value })
}
