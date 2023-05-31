import { request } from "../../request"
import { Paginate } from "@/types/common"
import { AssetCategory, AssetCategoryCondition } from "@/types/asset"

/**
 * 新增资产目录
 * @param AssetCategory - 资产目录信息
 */
export function addAssetCategory(AssetCategory: AssetCategory) {
  return request.post<AssetCategory>("/assetCategory/add", AssetCategory)
}

/**
 * 修改资产目录
 * @param AssetCategory - 资产目录信息
 */
export function updateAssetCategory(AssetCategory: AssetCategory) {
  return request.put<boolean>("/assetCategory/update", AssetCategory)
}

/**
 * 删除资产目录
 * @param id - 资产目录id
 */
export function deleteAssetCategory(id: string) {
  return request.delete<boolean>("/assetCategory/delete", {
    params: { id }
  })
}

/**
 * 获取资产目录列表
 * @param condition - 查询条件
 */
export function getAssetCategoryList(condition: AssetCategoryCondition) {
  return request.get<Paginate<AssetCategory>>("/assetCategory/list", {
    params: condition
  })
}

/**
 * 获取资产目录详情
 * @param id - 资产目录id
 */
export function getAssetCategoryDetail(id: string) {
  return request.get<AssetCategory>("/assetCategory/instance", {
    params: { id }
  })
}
