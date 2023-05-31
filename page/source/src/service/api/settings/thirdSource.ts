import { ThirdSource, ThirdSourceCondition } from "@/types/thirdSource"
import { request } from "../../request"
import type { Paginate } from "@/types/common"

/**
 * 新增三方源
 * @param thirdSource - 三方源信息
 */
export function addThirdSource(thirdSource: ThirdSource) {
  return request.post<ThirdSource>("/thirdSource/add", thirdSource)
}

/**
 * 修改三方源
 * @param thirdSource - 三方源信息
 */
export function updateThirdSource(thirdSource: ThirdSource) {
  return request.put<boolean>("/thirdSource/update", thirdSource)
}

/**
 * 删除三方源
 * @param id - 三方源id
 */
export function deleteThirdSource(id: string) {
  return request.delete<boolean>("/thirdSource/delete", {
    params: { id }
  })
}

/**
 * 获取三方源列表
 * @param condition - 查询条件
 */
export function getThirdSourceList(condition: ThirdSourceCondition) {
  return request.get<Paginate<ThirdSource>>("/thirdSource/list", {
    params: condition
  })
}

/**
 * 获取三方源详情
 * @param id - 三方源id
 */
export function getThirdSourceDetail(id: string) {
  return request.get<ThirdSource>("/thirdSource/instance", {
    params: { id }
  })
}
