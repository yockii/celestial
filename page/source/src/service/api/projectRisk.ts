import { request } from "../request"
import type { Paginate } from "@/types/common"
import { ProjectRisk, ProjectRiskCondition } from "@/types/project"

/**
 * 获取项目风险列表
 * @param condition - 查询条件
 */
export function getProjectRiskList(condition: ProjectRiskCondition) {
  return request.get<Paginate<ProjectRisk>>("/projectRisk/list", {
    params: condition
  })
}

/**
 * 新增项目风险
 * @param projectRisk - 项目风险信息
 */
export function addProjectRisk(projectRisk: ProjectRisk) {
  return request.post<boolean>("/projectRisk/add", projectRisk)
}

/**
 * 获取项目风险详情
 * @param id - 项目风险id
 */
export function getProjectRiskDetail(id: string) {
  return request.get<ProjectRisk>(`/projectRisk/instance?id=${id}`)
}

export type ProjectRiskCoefficient = {
  riskCoefficient: number
  maxRisk?: ProjectRisk
}
/**
 * 获取项目风险系数
 */
export function getProjectRiskCoefficient(projectId: string) {
  return request.get<ProjectRiskCoefficient>("/projectRisk/coefficient", {
    params: {
      projectId
    }
  })
}
