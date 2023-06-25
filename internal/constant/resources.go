package constant

const NeedLogin = "user"

// 首页
const (
	ResourceHome = "home"

	// 仪表盘
	ResourceDashboard = ResourceHome + ":dashboard"
)

// 项目
const (
	ResourceProject = "project"

	// 项目列表
	ResourceProjectList = ResourceProject + ":list"
	ResourceProjectAdd  = ResourceProject + ":add"

	// 项目详情
	ResourceProjectInstance = ResourceProject + ":detail"
	ResourceProjectUpdate   = ResourceProject + ":update"
	ResourceProjectDelete   = ResourceProject + ":delete"

	// 项目成员
	ResourceProjectMember       = ResourceProjectInstance + ":member"
	ResourceProjectMemberAdd    = ResourceProjectMember + ":add"
	ResourceProjectMemberDelete = ResourceProjectMember + ":delete"
	ResourceProjectMemberList   = ResourceProjectMember + ":list"
	ResourceProjectMemberUpdate = ResourceProjectMember + ":update"

	// 项目计划
	ResourceProjectPlan         = ResourceProjectInstance + ":plan"
	ResourceProjectPlanAdd      = ResourceProjectPlan + ":add"
	ResourceProjectPlanDelete   = ResourceProjectPlan + ":delete"
	ResourceProjectPlanUpdate   = ResourceProjectPlan + ":update"
	ResourceProjectPlanList     = ResourceProjectPlan + ":list"
	ResourceProjectPlanInstance = ResourceProjectPlan + ":instance"

	// 项目模块
	ResourceProjectModule       = ResourceProjectInstance + ":module"
	ResourceProjectModuleAdd    = ResourceProjectModule + ":add"
	ResourceProjectModuleDelete = ResourceProjectModule + ":delete"
	ResourceProjectModuleUpdate = ResourceProjectModule + ":update"
	ResourceProjectModuleList   = ResourceProjectModule + ":list"
	ResourceProjectModuleReview = ResourceProjectModule + ":review"

	// 项目需求
	ResourceProjectRequirement         = ResourceProjectInstance + ":requirement"
	ResourceProjectRequirementAdd      = ResourceProjectRequirement + ":add"
	ResourceProjectRequirementDelete   = ResourceProjectRequirement + ":delete"
	ResourceProjectRequirementUpdate   = ResourceProjectRequirement + ":update"
	ResourceProjectRequirementList     = ResourceProjectRequirement + ":list"
	ResourceProjectRequirementInstance = ResourceProjectRequirement + ":instance"
	// 设计完毕
	ResourceProjectRequirementStatusDesign = ResourceProjectRequirement + ":statusDesign"
	// 评审完毕
	ResourceProjectRequirementStatusReview = ResourceProjectRequirement + ":statusReview"
	// 已完成
	ResourceProjectRequirementStatusCompleted = ResourceProjectRequirement + ":statusCompleted"

	// 项目任务
	ResourceProjectTask         = ResourceProjectInstance + ":task"
	ResourceProjectTaskAdd      = ResourceProjectTask + ":add"
	ResourceProjectTaskDelete   = ResourceProjectTask + ":delete"
	ResourceProjectTaskUpdate   = ResourceProjectTask + ":update"
	ResourceProjectTaskList     = ResourceProjectTask + ":list"
	ResourceProjectTaskInstance = ResourceProjectTask + ":instance"
	ResourceProjectTaskCancel   = ResourceProjectTask + ":cancel"
	ResourceProjectTaskConfirm  = ResourceProjectTask + ":confirm"
	ResourceProjectTaskStart    = ResourceProjectTask + ":start"
	ResourceProjectTaskDone     = ResourceProjectTask + ":done"
	ResourceProjectTaskRestart  = ResourceProjectTask + ":restart"

	// 测试
	ResourceProjectTest         = ResourceProjectInstance + ":test"
	ResourceProjectTestAdd      = ResourceProjectTest + ":add"
	ResourceProjectTestDelete   = ResourceProjectTest + ":delete"
	ResourceProjectTestUpdate   = ResourceProjectTest + ":update"
	ResourceProjectTestClose    = ResourceProjectTestAdd
	ResourceProjectTestList     = ResourceProjectTest + ":list"
	ResourceProjectTestInstance = ResourceProjectTest + ":instance"
	// 测试用例
	ResourceProjectTestCase         = ResourceProjectTest + ":testCase"
	ResourceProjectTestCaseAdd      = ResourceProjectTestCase + ":add"
	ResourceProjectTestCaseDelete   = ResourceProjectTestCase + ":delete"
	ResourceProjectTestCaseUpdate   = ResourceProjectTestCase + ":update"
	ResourceProjectTestCaseList     = ResourceProjectTestCase + ":list"
	ResourceProjectTestCaseInstance = ResourceProjectTestCase + ":instance"
	// 测试用例项
	ResourceProjectTestCaseItem             = ResourceProjectTestCase + ":item"
	ResourceProjectTestCaseItemAdd          = ResourceProjectTestCaseItem + ":add"
	ResourceProjectTestCaseItemDelete       = ResourceProjectTestCaseItem + ":delete"
	ResourceProjectTestCaseItemUpdate       = ResourceProjectTestCaseItem + ":update"
	ResourceProjectTestCaseItemUpdateStatus = ResourceProjectTestCaseItem + ":updateStatus"
	ResourceProjectTestCaseItemList         = ResourceProjectTestCaseItem + ":list"
	ResourceProjectTestCaseItemInstance     = ResourceProjectTestCaseItem + ":instance"
	// 测试用例项步骤
	ResourceProjectTestCaseItemStep         = ResourceProjectTestCaseItem + ":step"
	ResourceProjectTestCaseItemStepAdd      = ResourceProjectTestCaseItemStep + ":add"
	ResourceProjectTestCaseItemStepDelete   = ResourceProjectTestCaseItemStep + ":delete"
	ResourceProjectTestCaseItemStepUpdate   = ResourceProjectTestCaseItemStep + ":update"
	ResourceProjectTestCaseItemStepList     = ResourceProjectTestCaseItemStep + ":list"
	ResourceProjectTestCaseItemStepInstance = ResourceProjectTestCaseItemStep + ":instance"

	// 项目变更
	ResourceProjectChange         = ResourceProjectInstance + ":change"
	ResourceProjectChangeAdd      = ResourceProjectChange + ":add"
	ResourceProjectChangeDelete   = ResourceProjectChange + ":delete"
	ResourceProjectChangeUpdate   = ResourceProjectChange + ":update"
	ResourceProjectChangeList     = ResourceProjectChange + ":list"
	ResourceProjectChangeInstance = ResourceProjectChange + ":instance"

	// 项目缺陷
	ResourceProjectIssue         = ResourceProjectInstance + ":issue"
	ResourceProjectIssueAdd      = ResourceProjectIssue + ":add"
	ResourceProjectIssueDelete   = ResourceProjectIssue + ":delete"
	ResourceProjectIssueUpdate   = ResourceProjectIssue + ":update"
	ResourceProjectIssueList     = ResourceProjectIssue + ":list"
	ResourceProjectIssueInstance = ResourceProjectIssue + ":instance"
	ResourceProjectIssueAssign   = ResourceProjectIssue + ":assign"
	ResourceProjectIssueStart    = ResourceProjectIssue + ":start"
	ResourceProjectIssueDone     = ResourceProjectIssue + ":done"
	ResourceProjectIssueVerify   = ResourceProjectIssue + ":verify"
	ResourceProjectIssueClose    = ResourceProjectIssue + ":close"

	// 项目风险
	ResourceProjectRisk         = ResourceProjectInstance + ":risk"
	ResourceProjectRiskAdd      = ResourceProjectRisk + ":add"
	ResourceProjectRiskDelete   = ResourceProjectRisk + ":delete"
	ResourceProjectRiskUpdate   = ResourceProjectRisk + ":update"
	ResourceProjectRiskList     = ResourceProjectRisk + ":list"
	ResourceProjectRiskInstance = ResourceProjectRisk + ":instance"

	// 项目资产
	ResourceProjectAsset         = ResourceProjectInstance + ":asset"
	ResourceProjectAssetAdd      = ResourceProjectAsset + ":add"
	ResourceProjectAssetDelete   = ResourceProjectAsset + ":delete"
	ResourceProjectAssetUpdate   = ResourceProjectAsset + ":update"
	ResourceProjectAssetList     = ResourceProjectAsset + ":list"
	ResourceProjectAssetInstance = ResourceProjectAsset + ":instance"
)

// 任务
const (
	ResourceTask = "task"
)

// 测试
const (
	ResourceTest = "test"
)

// 资产
const (
	ResourceAsset = "asset"

	// 文件
	ResourceFile         = ResourceAsset + ":file"
	ResourceFileAdd      = ResourceFile + ":add"
	ResourceFileDelete   = ResourceFile + ":delete"
	ResourceFileUpdate   = ResourceFile + ":update"
	ResourceFileList     = ResourceFile + ":list"
	ResourceFileInstance = ResourceFile + ":instance"
	ResourceFileDownload = ResourceFile + ":download"

	// 测试用例库
	ResourceCommonTestCase           = ResourceAsset + ":commonTestCase"
	ResourceCommonTestCaseAdd        = ResourceCommonTestCase + ":add"
	ResourceCommonTestCaseDelete     = ResourceCommonTestCase + ":delete"
	ResourceCommonTestCaseUpdate     = ResourceCommonTestCase + ":update"
	ResourceCommonTestCaseList       = ResourceCommonTestCase + ":list"
	ResourceCommonTestCaseInstance   = ResourceCommonTestCase + ":instance"
	ResourceCommonTestCaseAddItem    = ResourceCommonTestCase + ":addItem"
	ResourceCommonTestCaseDeleteItem = ResourceCommonTestCase + ":deleteItem"
	ResourceCommonTestCaseUpdateItem = ResourceCommonTestCase + ":updateItem"
)

// 系统
const (
	ResourceSystem = "system"

	// 用户
	ResourceUser                  = ResourceSystem + ":user"
	ResourceUserList              = ResourceUser + ":list"
	ResourceUserAdd               = ResourceUser + ":add"
	ResourceUserDelete            = ResourceUser + ":delete"
	ResourceUserUpdateUser        = ResourceUser + ":updateUser"
	ResourceUserUpdate            = ResourceUser + ":update"
	ResourceUserInstance          = ResourceUser + ":instance"
	ResourceUserDispatchRoles     = ResourceUser + ":dispatchRoles"
	ResourceUserRoles             = ResourceUser + ":roles"
	ResourceUserResetUserPassword = ResourceUser + ":resetUserPassword"

	// 角色
	ResourceRole       = ResourceSystem + ":role"
	ResourceRoleAdd    = ResourceRole + ":add"
	ResourceRoleDelete = ResourceRole + ":delete"
	ResourceRoleUpdate = ResourceRole + ":update"
	//ResourceRoleList              = ResourceRole + ":list" // 直接给权限
	ResourceRoleInstance          = ResourceRole + ":instance"
	ResourceRoleDispatchResources = ResourceRole + ":dispatchResources"
	ResourceRoleResources         = ResourceRole + ":resources"

	// 阶段
	ResourceStage         = ResourceSystem + ":stage"
	ResourceStageAdd      = ResourceStage + ":add"
	ResourceStageDelete   = ResourceStage + ":delete"
	ResourceStageUpdate   = ResourceStage + ":update"
	ResourceStageList     = ResourceStage + ":list"
	ResourceStageInstance = ResourceStage + ":instance"

	// 资产目录
	ResourceAssetCategory         = ResourceSystem + ":assetCategory"
	ResourceAssetCategoryAdd      = ResourceAssetCategory + ":add"
	ResourceAssetCategoryDelete   = ResourceAssetCategory + ":delete"
	ResourceAssetCategoryUpdate   = ResourceAssetCategory + ":update"
	ResourceAssetCategoryList     = ResourceAssetCategory
	ResourceAssetCategoryInstance = ResourceAssetCategory + ":instance"

	// 三方源
	ResourceThirdSource         = ResourceSystem + ":thirdSource"
	ResourceThirdSourceAdd      = ResourceThirdSource + ":add"
	ResourceThirdSourceDelete   = ResourceThirdSource + ":delete"
	ResourceThirdSourceUpdate   = ResourceThirdSource + ":update"
	ResourceThirdSourceList     = ResourceThirdSource + ":list"
	ResourceThirdSourceInstance = ResourceThirdSource + ":instance"
	ResourceThirdSourceSync     = ResourceThirdSource + ":sync"

	// oss配置
	ResourceOssConfig         = ResourceSystem + ":ossConfig"
	ResourceOssConfigAdd      = ResourceOssConfig + ":add"
	ResourceOssConfigDelete   = ResourceOssConfig + ":delete"
	ResourceOssConfigUpdate   = ResourceOssConfig + ":update"
	ResourceOssConfigList     = ResourceOssConfig + ":list"
	ResourceOssConfigInstance = ResourceOssConfig + ":instance"

	// 资源列表
	ResourceResourceList = ResourceSystem + ":resourceList"

	// 部门资源
	ResourceDepartment             = ResourceSystem + ":department"
	ResourceDepartmentAdd          = ResourceDepartment + ":add"
	ResourceDepartmentDelete       = ResourceDepartment + ":delete"
	ResourceDepartmentUpdate       = ResourceDepartment + ":update"
	ResourceDepartmentUpdateName   = ResourceDepartment + ":updateName"
	ResourceDepartmentChangeParent = ResourceDepartment + ":changeParent"
	ResourceDepartmentList         = ResourceDepartment + ":list"
	ResourceDepartmentInstance     = ResourceDepartment + ":instance"
	ResourceDepartmentAddUser      = ResourceDepartment + ":addUser"
	ResourceDepartmentRemoveUser   = ResourceDepartment + ":removeUser"
)
