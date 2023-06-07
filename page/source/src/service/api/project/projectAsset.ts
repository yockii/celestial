import { request } from "@/service/request"
import type { Paginate } from "@/types/common"
import { ProjectAsset, ProjectAssetCondition } from "@/types/project"

/**
 * 获取项目资产列表
 * @param condition - 查询条件
 */
export function getProjectAssetList(condition: ProjectAssetCondition) {
  return request.get<Paginate<ProjectAsset>>("/projectAsset/list", {
    params: condition
  })
}

/**
 * 新增项目资产
 * @param projectAsset - 项目资产信息
 */
export function addProjectAsset(projectAsset: ProjectAsset) {
  return request.post<boolean>("/projectAsset/add", projectAsset)
}

/**
 * 删除项目资产
 * @param id - 项目资产id
 */
export function deleteProjectAsset(id: string) {
  return request.delete<boolean>("/projectAsset/delete", {
    params: {
      id
    }
  })
}

/**
 * 更新项目资产信息
 * @param projectAsset - 项目资产信息
 */
export function updateProjectAsset(projectAsset: ProjectAsset) {
  return request.put<boolean>("/projectAsset/update", projectAsset)
}

/**
 * 获取项目资产详情
 * @param id - 项目资产id
 */
export function getProjectAssetDetail(id: string) {
  return request.get<ProjectAsset>(`/projectAsset/instance?id=${id}`)
}
