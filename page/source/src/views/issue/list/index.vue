<script setup lang="ts">
import { ProjectIssue, ProjectIssueCondition } from "@/types/project"
import dayjs from "dayjs"
import {
  DataTableFilterState,
  DataTableSortState,
  FormInst,
  NButton,
  NButtonGroup,
  NGrid,
  NGridItem,
  NIcon,
  NInput,
  NInputNumber,
  NPopconfirm,
  NPopselect,
  NTooltip,
  PaginationProps
} from "naive-ui"
import {
  addProjectIssue,
  assignProjectIssue,
  closeProjectIssue,
  deleteProjectIssue,
  finishProjectIssue,
  getProjectIssue,
  getProjectIssueList,
  rejectProjectIssue,
  reopenProjectIssue,
  startProjectIssue,
  updateProjectIssue,
  verifyProjectIssue
} from "@/service/api"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import { useUserStore } from "@/store/user"
import { Delete, Edit } from "@vicons/carbon"
import { AssignmentIndOutlined, CancelOutlined, CheckCircleFilled, PlayCircleFilled } from "@vicons/material"
import { CloseCircleFilled, ReloadOutlined } from "@vicons/antd"

const projectStore = useProjectStore()
const { project } = storeToRefs(projectStore)
const userStore = useUserStore()

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => projectStore.hasResourceCode("project:detail:plan:instance"),
  renderExpand: (row: ProjectIssue) => {
    if (!row.content || !row.issueCause || !row.solveMethod) {
      getProjectIssue(row.id).then((res) => {
        row.content = res.content
        row.issueCause = res.issueCause
        row.solveMethod = res.solveMethod
      })
    }
    return h(
      NGrid,
      {
        cols: 1,
        yGap: 8
      },
      () => [
        h(NGridItem, {}, { default: () => row.content }),
        h(NGridItem, {}, { default: () => "问题原因：" + (row.issueCause || "") }),
        h(NGridItem, {}, { default: () => "解决方法：" + (row.solveMethod || "") })
      ]
    )
  }
})
const startTimeColumn = reactive({
  title: "开始解决时间",
  key: "startTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => {
    return h(
      NTooltip,
      {
        placement: "top",
        trigger: "hover"
      },
      {
        default: () => (row.startTime ? dayjs(row.startTime).format("YYYY-MM-DD HH:mm:ss") : "未开始"),
        trigger: () => (row.startTime ? dayjs(row.startTime).format("YYYY-MM-DD") : "未开始")
      }
    )
  },
  // 排序
  sorter: true,
  sortOrder: false
})
const endTimeColumn = reactive({
  title: "解决完成时间",
  key: "endTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => {
    return h(
      NTooltip,
      {
        placement: "top",
        trigger: "hover"
      },
      {
        default: () => (row.endTime ? dayjs(row.endTime).format("YYYY-MM-DD HH:mm:ss") : "未解决"),
        trigger: () => (row.endTime ? dayjs(row.endTime).format("YYYY-MM-DD") : "未解决")
      }
    )
  },
  // 排序
  sorter: true,
  sortOrder: false
})
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectIssue) => {
    switch (row.status) {
      case -1:
        return h(
          NTooltip,
          {
            placement: "top",
            trigger: "hover"
          },
          {
            trigger: () => "已驳回",
            default: () => "已驳回原因：" + (row.rejectedReason || "")
          }
        )
      case 1:
        return "新建"
      case 2:
        return "已指定"
      case 3:
        return "处理中"
      case 4:
        return "待验证"
      case 5:
        return "已解决"
      case 9:
        return "已关闭"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "已驳回",
      value: -1
    },
    {
      label: "新建",
      value: 1
    },
    {
      label: "已指定",
      value: 2
    },
    {
      label: "处理中",
      value: 3
    },
    {
      label: "待验证",
      value: 4
    },
    {
      label: "已解决",
      value: 5
    },
    {
      label: "已关闭",
      value: 9
    }
  ]
})
const dealUserColumn = reactive({
  title: "当前处理人",
  key: "dealUser",
  render: (row: ProjectIssue) => {
    if (row.assigneeId && project.value?.members) {
      const user = project.value?.members.find((item) => item.userId === row.assigneeId)
      if (user) {
        return user.realName
      }
    }
    return "未知"
  }
})
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const doneSolveDuration = ref<number>(0)
const rejectReason = ref<string>("")
const operationColumn = reactive({
  title: "操作",
  key: "operation",
  render: (row: ProjectIssue) => {
    const btnGroup: VNode[] = []
    if (
      row.status !== 9 &&
      row.status !== 5 &&
      projectStore.hasResourceCode("project:detail:issue:assign") &&
      ((row.status === 1 && row.creatorId === userStore.user.id) || row.assigneeId === userStore.user.id)
    ) {
      btnGroup.push(
        h(
          NPopselect,
          {
            options: projectStore.memberList?.map((item) => ({
              label: item.realName,
              value: item.userId
            })),
            value: row.assigneeId || "",
            valueField: "userId",
            onUpdateValue: (value: string) => {
              assignProjectIssue(row.id, value).then((res) => {
                if (res) {
                  message.success("指派成功")
                  refresh()
                }
              })
            }
          },
          () =>
            h(
              NTooltip,
              {},
              {
                default: () => "指派处理人",
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: "small",
                      secondary: true,
                      type: "primary",
                      disabled:
                        !projectStore.hasResourceCode("project:detail:issue:assign") ||
                        (row.creatorId !== userStore.user.id && row.assigneeId !== userStore.user.id)
                    },
                    {
                      default: () => h(NIcon, { component: AssignmentIndOutlined })
                    }
                  )
              }
            )
        )
      )
    }

    // 如果是已指派 状态，且当前用户是处理人, 则可以点击开始解决，或者拒绝缺陷
    if (projectStore.hasResourceCode("project:detail:issue:start") && row.status === 2 && row.assigneeId === userStore.user.id) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "开始解决",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  secondary: true,
                  type: "primary",
                  disabled: !projectStore.hasResourceCode("project:detail:issue:start"),
                  onClick: () => handleStartData(row)
                },
                {
                  default: () => h(NIcon, { component: PlayCircleFilled })
                }
              )
          }
        )
      )

      // 拒绝缺陷
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "拒绝该缺陷",
            trigger: () =>
              h(
                NPopconfirm,
                {
                  onPositiveClick: () => handleRejectData(row)
                },
                {
                  default: () =>
                    h(NInput, {
                      value: rejectReason.value,
                      clearable: true,
                      onUpdateValue: (v: string | [string, string]) => {
                        if (Array.isArray(v)) {
                          rejectReason.value = v[0]
                        } else {
                          rejectReason.value = v
                        }
                      }
                    }),
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        secondary: true,
                        type: "primary",
                        disabled: !projectStore.hasResourceCode("project:detail:issue:reject")
                        // onClick: () => handleFinishData(row)
                      },
                      {
                        default: () => h(NIcon, { component: CancelOutlined })
                      }
                    )
                }
              )
          }
        )
      )
    }

    // 若是解决中的状态，且当前用户是处理人，则可以点击完成，并更新原因和解决方案
    if (projectStore.hasResourceCode("project:detail:issue:done") && row.status === 3 && row.assigneeId === userStore.user.id) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "处理完成",
            trigger: () =>
              h(
                NPopconfirm,
                {
                  onPositiveClick: () => handleFinishData(row)
                },
                {
                  default: () =>
                    h(
                      NInputNumber,
                      {
                        value: doneSolveDuration.value,
                        clearable: true,
                        precision: 2,
                        onUpdateValue: (v: number | null) => {
                          doneSolveDuration.value = v || 0
                        }
                      },
                      {
                        suffix: () => "小时"
                      }
                    ),
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        secondary: true,
                        type: "primary",
                        disabled: !projectStore.hasResourceCode("project:detail:issue:finish")
                        // onClick: () => handleFinishData(row)
                      },
                      {
                        default: () => h(NIcon, { component: CheckCircleFilled })
                      }
                    )
                }
              )
          }
        )
      )
    }

    // 待验证状态，且当前用户是处理人，则可以点击验证
    if (projectStore.hasResourceCode("project:detail:issue:verify") && row.status === 4 && row.assigneeId === userStore.user.id) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "缺陷验证",
            trigger: () =>
              h(
                NPopconfirm,
                {
                  positiveText: "验证通过",
                  negativeText: "验证不通过",
                  onPositiveClick: () => handleVerifyData(row, 5),
                  onNegativeClick: () => handleVerifyData(row, 2)
                },
                {
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        secondary: true,
                        type: "primary",
                        disabled: !projectStore.hasResourceCode("project:detail:issue:verify")
                      },
                      {
                        default: () => h(NIcon, { component: CheckCircleFilled })
                      }
                    ),
                  default: () => "缺陷验证结果"
                }
              )
          }
        )
      )
    }

    // 如果是已解决状态，且当前用户是处理人，则可以点击关闭
    if (projectStore.hasResourceCode("project:detail:issue:close") && row.status === 5 && row.assigneeId === userStore.user.id) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "关闭缺陷",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  secondary: true,
                  type: "primary",
                  disabled: !projectStore.hasResourceCode("project:detail:issue:close"),
                  onClick: () => handleCloseData(row)
                },
                {
                  default: () => h(NIcon, { component: CloseCircleFilled })
                }
              )
          }
        )
      )
    }

    // 如果是已关闭或已驳回状态，且当前用户是处理人，则可以点击重新打开
    if (projectStore.hasResourceCode("project:detail:issue:reopen") && (row.status === 9 || row.status === -1) && row.assigneeId === userStore.user.id) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "重新打开",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  secondary: true,
                  type: "primary",
                  disabled: !projectStore.hasResourceCode("project:detail:issue:reopen"),
                  onClick: () => handleReopenData(row)
                },
                {
                  default: () => h(NIcon, { component: ReloadOutlined })
                }
              )
          }
        )
      )
    }

    if (projectStore.hasResourceCode("project:detail:issue:update") && row.status !== 9) {
      btnGroup.push(
        h(
          NTooltip,
          {},
          {
            default: () => "编辑",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  secondary: true,
                  type: "primary",
                  disabled: !projectStore.hasResourceCode("project:detail:issue:update"),
                  onClick: () => handleEditData(row)
                },
                {
                  default: () => h(NIcon, { component: Edit })
                }
              )
          }
        )
      )
    }
    if (projectStore.hasResourceCode("project:detail:issue:delete")) {
      btnGroup.push(
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteData(row.id)
          },
          {
            default: () => "确认删除",
            trigger: () =>
              h(
                NTooltip,
                {},
                {
                  default: () => "删除",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "error",
                        disabled: !projectStore.hasResourceCode("project:detail:issue:delete")
                      },
                      {
                        default: () => h(NIcon, { component: Delete })
                      }
                    )
                }
              )
          }
        )
      )
    }
    return h(NButtonGroup, {}, () => btnGroup)
  }
})
const columns = [
  expandColumn,
  {
    title: "标题",
    key: "title"
  },
  startTimeColumn,
  endTimeColumn,
  statusColumn,
  dealUserColumn,
  createTimeColumn,
  operationColumn
]
const condition = ref<ProjectIssueCondition>({
  projectId: project.value.id
})
const list = ref<ProjectIssue[]>([])
const loading = ref(false)

const refresh = () => {
  getProjectIssueList(condition.value).then((res) => {
    list.value = res.items || []
    paginationReactive.itemCount = res.total
    paginationReactive.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
    paginationReactive.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
    statusColumn.filterOptionValues = [condition.value.status || 0]
  })
}

const paginationReactive = reactive({
  itemCount: 0,
  page: 1,
  pageCount: 1,
  pageSize: 10,
  prefix({ itemCount }: PaginationProps) {
    return `共${itemCount}条`
  }
})

const handlePageChange = (page: number) => {
  condition.value.offset = (page - 1) * (condition.value.limit || 10)
  refresh()
}
const handlePageSizeChange = (pageSize: number) => {
  condition.value.limit = pageSize
  refresh()
}
const handleFiltersChange = (filters: DataTableFilterState) => {
  if (!loading.value) {
    const filterValues = filters.status || []
    if (filterValues instanceof Array) {
      condition.value.status = (filterValues[0] as number) || 0
    }
    refresh()
  }
}
const handleSorterChange = (sorter: DataTableSortState) => {
  if (!loading.value) {
    const { columnKey, order } = sorter
    let field = "start_time"
    if (columnKey === "createTime") {
      createTimeColumn.sortOrder = order === "ascend"
      field = "create_time"
    } else if (columnKey === "startTime") {
      startTimeColumn.sortOrder = order === "ascend"
      field = "start_time"
    } else if (columnKey === "endTime") {
      endTimeColumn.sortOrder = order === "ascend"
      field = "end_time"
    }
    condition.value.orderBy = field + (order === "ascend" ? " asc" : " desc")
    refresh()
  }
}

// 抽屉部分，新建、编辑内容
const instance = ref<ProjectIssue>({ endTime: 0, id: "", title: "", projectId: "", startTime: 0, status: 0 })
const drawerActive = ref(false)
const isUpdate = computed(() => !!instance.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "编辑缺陷" : "新建缺陷"))
const resetInstance = (origin: ProjectIssue | undefined = undefined) => {
  if (origin) {
    instance.value = JSON.parse(JSON.stringify(origin))
  } else {
    instance.value = { endTime: 0, id: "", title: "", projectId: project.value.id, startTime: 0, status: 0 }
  }
}
const newInstance = () => {
  resetInstance()
  drawerActive.value = true
}
const planRules = {
  title: [
    { required: true, message: "请输入缺陷名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: [{ type: "number", required: true, message: "请选择缺陷类型", trigger: ["blur", "change"] }]
}
const handleEditData = (row: ProjectIssue) => {
  instance.value = row
  drawerActive.value = true
}
const handleReopenData = (row: ProjectIssue) => {
  reopenProjectIssue(row.id).then((res) => {
    if (res) {
      message.success("重新打开成功")
      refresh()
    }
  })
}
const handleDeleteData = (id: string) => {
  deleteProjectIssue(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
const handleStartData = (row: ProjectIssue) => {
  startProjectIssue(row.id).then((res) => {
    if (res) {
      message.success("启动成功")
      refresh()
    }
  })
}
const handleFinishData = (row: ProjectIssue) => {
  if (doneSolveDuration.value <= 0) {
    message.error("请输入处理时长")
    return
  }
  finishProjectIssue(row.id, doneSolveDuration.value).then((res) => {
    if (res) {
      message.success("处理成功")
      doneSolveDuration.value = 0
      refresh()
    }
  })
}
const handleRejectData = (row: ProjectIssue) => {
  rejectProjectIssue(row.id, rejectReason.value).then((res) => {
    if (res) {
      message.success("驳回成功")
      rejectReason.value = ""
      refresh()
    }
  })
}
const handleVerifyData = (row: ProjectIssue, status: number) => {
  verifyProjectIssue(row.id, status).then((res) => {
    if (res) {
      message.success("验证成功")
      refresh()
    }
  })
}
const handleCloseData = (row: ProjectIssue) => {
  closeProjectIssue(row.id).then((res) => {
    if (res) {
      message.success("关闭成功")
      refresh()
    }
  })
}
const formRef = ref<FormInst>()
const message = useMessage()
const submit = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (isUpdate.value) {
        updateProjectIssue(instance.value).then((res) => {
          if (res) {
            message.success("保存成功")
            drawerActive.value = false
            refresh()
          }
        })
      } else {
        addProjectIssue(instance.value).then((res) => {
          if (res) {
            message.success("保存成功")
            drawerActive.value = false
            refresh()
          }
        })
      }
    }
  })
}

defineExpose({
  projectSelected: (id: string) => {
    condition.value.projectId = id || ""
    refresh()
  }
})

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: ProjectIssue) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :on-update:sorter="handleSorterChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>
</template>
