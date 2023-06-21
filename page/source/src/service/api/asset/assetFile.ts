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

/**
 * 新增资产文件
 * @param file 资产文件
 */
export const addAssetFile = (file: File) => {
  return request.post<File>("/assetFile/add", file)
}

/**
 * 修改资产文件
 * @param file 资产文件
 */
export const updateAssetFile = (file: File) => {
  return request.put<File>("/assetFile/update", file)
}

/**
 * 下载资产文件
 * @param id 资产文件id
 */
export const downloadAssetFile = (id: string) => {
  return request.get<Blob>("/assetFile/download", {
    params: { id },
    responseType: "blob"
  })
}
