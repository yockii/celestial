import { ThirdSource, ThirdSourceCondition, ThirdSourcePublic } from "@/types/thirdSource"
import { request } from "../../request"
import type { Paginate } from "@/types/common"
import { LoginResponse } from "@/types/user"

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

/**
 * 获取第三方源公共信息
 */
export function getThirdSourcePublic() {
  return request.get<ThirdSourcePublic[]>("/thirdSource/publicList")
}

/**
 * 根据钉钉authCode从第三方源登录并获取用户信息
 * @param thirdSourceId - 三方源id
 * @param authCode - 钉钉authCode
 */
export function loginByDingTalk(thirdSourceId: string, code: string) {
  return request.post<LoginResponse>("/loginByDingtalkCode", {
    thirdSourceId,
    code
  })
}

/**
 * 同步三方源数据
 * @param id - 三方源id
 */
export function syncThirdSourceData(id: string) {
  return request.post<boolean>(`/thirdSource/sync?id=${id}`, null, {
    timeout: 5 * 60 * 1000
  })
}
