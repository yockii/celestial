import { ProjectTest, ProjectTestCase, ProjectTestCaseCondition, ProjectTestCaseItem } from "@/types/project"
import { request } from "../../request"
import { Paginate } from "@/types/common"

/**
 * 获取项目测试轮列表
 * @param projectId - 项目id
 */
export function getProjectTestList(projectId: string) {
  return request.get<ProjectTest[]>(`/project/test/list?projectId=${projectId}`)
}

/**
 * 新增项目测试轮
 * @param projectTest - 项目测试轮信息
 */
export function addProjectTest(projectTest: ProjectTest) {
  return request.post<ProjectTest>("/project/test/add", projectTest)
}

/**
 * 关闭项目测试轮
 * @param id - 项目测试轮id
 */
export function closeProjectTest(id: string) {
  return request.put<boolean>("/project/test/close", {
    id
  })
}

/**
 * 获取测试轮详情
 * @param id - 项目测试轮id
 */
export function getProjectTest(id: string) {
  return request.get<ProjectTest>(`/project/test/instance?id=${id}`)
}

/**
 * 获取带items的测试用例列表
 * @param condition - 查询条件
 */
export function getProjectTestCaseListWithItems(condition: ProjectTestCaseCondition) {
  return request.get<Paginate<ProjectTestCase>>("/project/testCase/listWithItems", {
    params: condition
  })
}

/**
 * 新增项目测试用例
 * @param projectTestCase - 项目测试用例
 */
export function addProjectTestCase(projectTestCase: ProjectTestCase) {
  return request.post<ProjectTestCase>("/project/testCase/add", projectTestCase)
}

/**
 * 更新项目测试用例
 * @param projectTestCase - 项目测试用例
 */
export function updateProjectTestCase(projectTestCase: ProjectTestCase) {
  return request.put<boolean>("/project/testCase/update", projectTestCase)
}

/**
 * 批量提交项目测试用例
 * @param projectTestCase[] - 项目测试用例
 */
export function batchSubmitProjectTestCase(projectTestCase: ProjectTestCase[]) {
  return request.post<boolean>("/project/testCase/batchAdd", projectTestCase)
}

/**
 * 删除项目测试用例
 * @param id - 项目测试用例id
 */
export function deleteProjectTestCase(id: string) {
  return request.delete<boolean>(`/project/testCase/delete?id=${id}`)
}

/**
 * 添加测试用例项
 * @param projectTestCaseItem - 测试用例项
 */
export function addProjectTestCaseItem(projectTestCaseItem: ProjectTestCaseItem) {
  return request.post<ProjectTestCaseItem>("/project/testCaseItem/add", projectTestCaseItem)
}

/**
 * 更新测试用例项
 * @param projectTestCaseItem - 测试用例项
 */
export function updateProjectTestCaseItem(projectTestCaseItem: ProjectTestCaseItem) {
  return request.put<boolean>("/project/testCaseItem/update", projectTestCaseItem)
}

/**
 * 批量提交测试用例项
 * @param projectTestCaseItem[] - 测试用例项
 */
export function batchSubmitProjectTestCaseItem(projectTestCaseItem: ProjectTestCaseItem[]) {
  return request.post<boolean>("/project/testCaseItem/batchAdd", projectTestCaseItem)
}

/**
 * 删除测试用例项
 * @param id - 测试用例项id
 */
export function deleteProjectTestCaseItem(id: string) {
  return request.delete<boolean>(`/project/testCaseItem/delete?id=${id}`)
}

/**
 * 更新测试用例项状态
 * @param id - 测试用例项id
 * @param status - 状态
 */
export function updateProjectTestCaseItemStatus(id: string, status: number) {
  return request.put<boolean>("/project/testCaseItem/updateStatus", {
    id,
    status
  })
}
