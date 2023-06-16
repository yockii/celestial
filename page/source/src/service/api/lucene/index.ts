import { LuceneDocument } from "@/types/lecene"
import { request } from "../../request"
import type { Paginate } from "@/types/common"

/**
 * 进行全文检索
 * @param keyword - 关键字
 * @param limit - 每页数量
 * @param offset - 偏移量
 */
export function luceneSearch(keyword: string, limit: number, offset: number) {
  return request.get<Paginate<LuceneDocument>>("/search", {
    params: {
      keyword,
      limit,
      offset
    }
  })
}
