import { Condition } from "./common"

export type ProjectMember = {
  userId: string
  username: string
  realName: string
  roleId: string
}

export type Project = {
  id: string
  name: string
  code: string
  description: string
  stageId: string
  createTime?: number
  members?: ProjectMember[]
}
export type ProjectCondition = Condition & {
  name: string
  stageId: string
}

export type ProjectStageStatistics = {
  stageId: string
  count: number
}

export type ProjectPlan = {
  id: string
  projectId: string
  stageId?: string
  planName: string
  planDesc?: string
  startTime: number
  endTime: number
  target?: string
  scope?: string
  schedule?: string
  resource?: string
  budget?: number
  status: number
  actualStartTime?: number
  actualEndTime?: number
  createTime?: number
  updateTime?: number
}
export type ProjectPlanCondition = Condition & {
  projectId: string
  planName?: string
  stageId?: string
  status?: number
}

export type ProjectRisk = {
  id: string
  projectId: string
  stageId: string
  riskName: string
  riskProbability: number
  riskImpact: number
  riskLevel: number
  status: number
  response: string
  startTime: number
  endTime: number
  result: string
  createTime?: number
}
export type ProjectRiskCondition = Condition & {
  name: string
  stageId: string
}

export type ProjectTaskWorkTimeStatistics = {
  projectId: string
  taskCount: number
  estimateDuration: number
  actualDuration: number
}
export type Stage = {
  id: string
  name: string
  orderNum: number
  status: number
  createTime?: number
}

export type StageCondition = Condition & {
  name: string
  status: number
}

export type ProjectModule = {
  id: string
  projectId: string
  parentId?: string
  name: string
  alias?: string
  remark?: string
  childrenCount?: number
  children?: ProjectModule[]
  fullPath?: string
  status: number
  createTime?: number
}

export type ProjectModuleCondition = Condition & {
  projectId: string
  parentId?: string
  name?: string
  status?: number
}

export type ProjectRequirement = {
  id: string
  projectId: string
  moduleId?: string
  name: string
  type?: number
  detail?: string
  priority?: number
  stageId?: string
  source?: number
  ownerId?: string
  feasibility?: number
  status?: number
  createTime?: number
}

export type ProjectRequirementCondition = Condition & {
  projectId: string
  status?: number
  feasibility?: number
  type?: number
  priority?: number
  fullPath?: string
}

export type ProjectTask = {
  id: string
  projectId: string
  moduleId?: string
  requirementId?: string
  name: string
  stageId?: string
  parentId?: string
  startTime?: number
  endTime?: number
  taskDesc?: string
  priority?: number
  ownerId?: string
  status?: number
  actualStartTime?: number
  actualEndTime?: number
  estimateDuration?: number
  actualDuration?: number
  createTime?: number
  updateTime?: number
  creatorId?: string
  members?: ProjectMember[]
}

export type ProjectTaskCondition = Condition & {
  projectId: string
  onlyParent?: boolean
  name?: string
  moduleId?: string
  requirementId?: string
  status?: number
  priority?: number
  ownerId?: string
  stageId?: string
  parentId?: string
  fullPath?: string
}
