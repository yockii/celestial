import { request } from "../../request"
import { Paginate } from "@/types/common"
import { File, FileCondition } from "@/types/asset"

/**
 * 获取资产文件详情
 * @param id 资产文件id
 */
export const getAssetFile = (id: string) => {
  return request.get<File>("/assetFile/instance", {
    params: { id }
  })
}

/**
 * 获取资产文件列表
 * @param condition 查询条件
 */
export const getAssetFileList = (condition: FileCondition) => {
  return request.get<Paginate<File>>("/assetFile/list", {
    params: condition
  })
}
