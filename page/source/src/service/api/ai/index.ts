import { request } from "../../request"

/**
 * 调用ai接口获取数据
 */
export function getAiData(question: string) {
    return request.get<string>("/ai/dataSearch", {
        params: { question },
        timeout: 5 * 60 * 1000
    })
}