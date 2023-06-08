import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { OssConfig, OssConfigCondition } from "@/types/ossConfig"

/**
 * 新增oss配置
 * @param ossConfig - oss配置信息
 */
export function addOssConfig(ossConfig: OssConfig) {
  return request.post<OssConfig>("/ossConfig/add", ossConfig)
}

/**
 * 修改oss配置
 * @param ossConfig - oss配置信息
 */
export function updateOssConfig(ossConfig: OssConfig) {
  return request.put<boolean>("/ossConfig/update", ossConfig)
}

/**
 * 删除oss配置
 * @param id - oss配置id
 */
export function deleteOssConfig(id: string) {
  return request.delete<boolean>("/ossConfig/delete", {
    params: { id }
  })
}

/**
 * 获取oss配置列表
 * @param condition - 查询条件
 */
export function getOssConfigList(condition: OssConfigCondition) {
  return request.get<Paginate<OssConfig>>("/ossConfig/list", {
    params: condition
  })
}

/**
 * 获取oss配置详情
 * @param id - oss配置id
 */
export function getOssConfigDetail(id: string) {
  return request.get<OssConfig>("/ossConfig/instance", {
    params: { id }
  })
}

/**
 * 更新OSS配置状态
 * @param id - oss配置id
 * @param status - 状态
 */
export function updateOssConfigStatus(id: string, status: number) {
  return request.put<boolean>("/ossConfig/updateStatus", {
    id,
    status
  })
}
