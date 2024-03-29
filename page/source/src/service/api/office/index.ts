import { OfficeConfig } from "@/types/office"
import { request } from "../../request"

/**
 * 获取编辑器url信息
 */
export function getEditorUrl() {
  return request.get<string>("/office/editorUrl")
}

/**
 * 获取配置信息
 * @param fileId - 文件 ID
 * @param id - 版本id
 */
export function getFileConfig(fileId: string, id?: string) {
  return request.get<OfficeConfig>("/office/config", {
    params: {
      fileId,
      id,
      baseUrl: window.location.origin
    }
  })
}
