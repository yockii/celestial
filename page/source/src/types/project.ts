import { Condition } from "./common"

export type ProjectMember = {
  userId: string
  username: string
  realName: string
  roleId: string
  status?: number
}

export type Project = {
  id: string
  name: string
  code: string
  description: string
  stageId: string
  ownerId?: string
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
  projectId?: string
  stageId?: string
  planName?: string
  planDesc?: string
  startTime?: number
  endTime?: number
  target?: string
  scope?: string
  schedule?: string
  resource?: string
  budget?: number
  status?: number
  actualStartTime?: number
  actualEndTime?: number
  createTime?: number
  updateTime?: number
}
export type ProjectPlanCondition = Condition & {
  id?: string
  projectId: string
  planName?: string
  stageId?: string
  status?: number
}

export type ProjectRisk = {
  id: string
  projectId: string
  stageId?: string
  riskName: string
  riskDesc?: string
  riskProbability?: number
  riskImpact?: number
  riskLevel?: number
  status: number
  response?: string
  startTime?: number
  endTime?: number
  result?: string
  createTime?: number
}
export type ProjectRiskCondition = Condition & {
  id?: string
  name?: string
  projectId: string
  status?: number
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
  id?: string
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
  id?: string
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
  members?: ProjectTaskMember[]
  children?: ProjectTask[]
  isLeaf?: boolean
  childrenCount?: number
}

export type ProjectTaskCondition = Condition & {
  id?: string
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

export type ProjectTaskMember = {
  id: string
  projectId: string
  userId: string
  taskId: string
  realName?: string
  roleId?: string
  estimateDuration?: number
  actualDuration?: number
  status?: number
}

export type ProjectAsset = {
  id: string
  projectId: string
  fileId: string
  name: string
  type: number
  status: number
  remark?: string
  createTime?: number
}

export type ProjectAssetCondition = Condition & {
  projectId: string
  type: number
  status: number
}

export type ProjectTest = {
  id?: string
  projectId: string
  round?: number
  testRecord?: string
  remark?: string
  startTime?: number
  endTime?: number
}

export type ProjectTestCase = {
  id: string
  projectId: string
  relatedId?: string
  relatedType?: number
  name: string
  remark?: string
  creatorId?: string
  createTime?: number
  items?: ProjectTestCaseItem[]
}
export type ProjectTestCaseCondition = Condition & {
  projectId: string
  name?: string
}

export type ProjectTestCaseItem = {
  id?: string
  projectId?: string
  testCaseId?: string
  name?: string
  type?: number
  content?: string
  status?: number
  steps?: ProjectTestCaseItemStep[]
}

export type ProjectTestCaseItemStep = {
  id?: string
  caseItemId?: string
  orderNum?: number
  content?: string
  expect?: string
  status?: number
}

export type ProjectChange = {
  id: string
  projectId: string
  title: string
  type?: number
  level?: number
  reason?: string
  plan?: string
  review?: string
  risk?: string
  status?: number
  applyUserId?: string
  reviewerIdList?: string
  result?: string
  reviewTime?: number
  createTime?: number
}

export type ProjectChangeCondition = Condition & {
  id?: string
  projectId: string
  type?: number
  level?: number
  status?: number
}

export type ProjectIssue = {
  id: string
  projectId: string
  title: string
  type?: number
  level?: number
  content?: string
  status?: number
  assigneeId?: string
  creatorId?: string
  createTime?: number
  updateTime?: number
  startTime?: number
  endTime?: number
  resolvedTime?: number
  solveDuration?: number
  rejectedReason?: string
  issueCause?: string
  solveMethod?: string
}

export type ProjectIssueCondition = Condition & {
  id?: string
  projectId: string
  type?: number
  level?: number
  status?: number
  assigneeId?: string
  creatorId?: string
}
