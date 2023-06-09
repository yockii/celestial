import { request } from "../../request"
import { Paginate } from "@/types/common"
import type { Resource, ResourceCondition } from "@/types/user"

/**
 * 获取资源列表
 * @param condition - 查询条件
 */
export function getResourceList(condition: ResourceCondition) {
  return request.get<Paginate<Resource>>("resource/list", {
    params: condition
  })
}
