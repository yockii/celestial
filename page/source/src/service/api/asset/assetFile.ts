import { request } from "../../request"
import { Paginate } from "@/types/common"
import { File, FileCondition, FileUser, FileVersion, FileVersionCondition } from "@/types/asset"

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
    responseType: "blob",
    timeout: 10 * 60 * 1000
  })
}

/**
 * 获取资产文件的用户权限列表
 * @param params 查询条件
 */
export const getAssetFileUserList = (params: { fileId: string; permission?: number; realName?: string }) => {
  return request.get<FileUser[]>("/assetFile/filePermissionUsers", {
    params
  })
}

/**
 * 修改资产文件的用户权限
 * @param fileUser 资产文件用户权限
 */
export const updateAssetFileUser = (fileUser: FileUser) => {
  return request.post<FileUser>("/assetFile/updateFileUserPermission", fileUser)
}

/**
 * 删除资产文件的用户权限
 * @param id 资产文件用户权限ID
 */
export const deleteAssetFileUser = (id: string) => {
  return request.delete<FileUser>("/assetFile/removeFileUserPermission", {
    params: { id }
  })
}

/**
 * 获取资产文件的版本记录
 * @param fileId 资产文件ID
 */
export const getAssetFileVersionList = (condition: FileVersionCondition) => {
  return request.get<Paginate<FileVersion>>("/assetFile/versionList", {
    params: condition
  })
}
