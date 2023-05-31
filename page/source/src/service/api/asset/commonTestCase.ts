import { request } from "../../request"
import { Paginate } from "@/types/common"
import { CommonTestCase, CommonTestCaseCondition, CommonTestCaseItem } from "@/types/asset"

/**
 * 新增通用测试用例
 * @param CommonTestCase - 通用测试用例信息
 */
export function addCommonTestCase(CommonTestCase: CommonTestCase) {
  return request.post<CommonTestCase>("/commonTestCase/add", CommonTestCase)
}

/**
 * 新增通用测试用例项
 * @param CommonTestCaseItem - 通用测试用例项信息
 */
export function addCommonTestCaseItem(CommonTestCaseItem: CommonTestCaseItem) {
  return request.post<CommonTestCaseItem>("/commonTestCaseItem/add", CommonTestCaseItem)
}

/**
 * 修改通用测试用例
 * @param CommonTestCase - 通用测试用例信息
 */
export function updateCommonTestCase(CommonTestCase: CommonTestCase) {
  return request.put<boolean>("/commonTestCase/update", CommonTestCase)
}

/**
 * 修改通用测试用例项
 * @param CommonTestCaseItem - 通用测试用例项信息
 */
export function updateCommonTestCaseItem(CommonTestCaseItem: CommonTestCaseItem) {
  return request.put<boolean>("/commonTestCaseItem/update", CommonTestCaseItem)
}

/**
 * 删除通用测试用例
 * @param id - 通用测试用例id
 */
export function deleteCommonTestCase(id: string) {
  return request.delete<boolean>("/commonTestCase/delete", {
    params: { id }
  })
}

/**
 * 删除通用测试用例项
 * @param id - 通用测试用例项id
 */
export function deleteCommonTestCaseItem(id: string) {
  return request.delete<boolean>("/commonTestCaseItem/delete", {
    params: { id }
  })
}

/**
 * 获取通用测试用例列表
 * @param condition - 查询条件
 */
export function getCommonTestCaseList(condition: CommonTestCaseCondition) {
  return request.get<Paginate<CommonTestCase>>("/commonTestCase/listWithItem", {
    params: condition
  })
}

/**
 * 获取通用测试用例详情
 * @param id - 通用测试用例id
 */
export function getCommonTestCaseDetail(id: string) {
  return request.get<CommonTestCase>("/commonTestCase/instance", {
    params: { id }
  })
}
