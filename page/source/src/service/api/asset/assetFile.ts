import { request } from "../../request"
import { Paginate } from "@/types/common"
import { File } from "@/types/asset"

/**
 * 获取资产文件详情
 * @param id 资产文件id
 */
export const getAssetFile = (id: string) => {
  return request.get<File>("/assetFile/instance", {
    params: { id }
  })
}
