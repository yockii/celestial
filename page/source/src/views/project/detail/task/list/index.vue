<script setup lang="ts">
import {
  addProjectIssue,
  getProjectTask,
  finishProjectTask,
  getProjectTaskList,
  startProjectTask,
  testPassProjectTask,
  testingProjectTask
} from "@/service/api"
import { ProjectTask, ProjectTaskCondition, ProjectTaskMember, ProjectIssue } from "@/types/project"
import {
  NAvatarGroup,
  NBadge,
  NPopconfirm,
  DataTableBaseColumn,
  DataTableFilterState,
  NButtonGroup,
  NButton,
  NIcon,
  PaginationProps,
  NTooltip,
  NDropdown,
  NAvatar,
  NSpace,
  NText,
  FormInst
} from "naive-ui"
import { ArrowsSplit } from "@vicons/tabler"
import { Debug, Delete, Edit, FaceDissatisfied, FaceSatisfied, ParentChild } from "@vicons/carbon"
import NameAvatar from "@/components/NameAvatar.vue"
import { useUserStore } from "@/store/user"
import { useProjectStore } from "@/store/project"
import { BoxCheckmark20Regular, Checkmark20Regular, Info20Regular, Play20Regular } from "@vicons/fluent"
import dayjs from "dayjs"
import WorkTimeDrawer from "../workTimeDrawer/index.vue"
import IssueForm from "@/components/project/issue/IssueForm.vue"

const router = useRouter()
const userStore = useUserStore()
const projectStore = useProjectStore()
const message = useMessage()

const props = defineProps<{
  condition: ProjectTaskCondition
  useTree?: boolean
}>()
const emit = defineEmits(["showChild", "newChild", "edit", "delete", "showDetail"])

const condition = toRef(props, "condition")
const list = ref<ProjectTask[]>([])

const loading = ref(false)
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
const refresh = () => {
  loading.value = true
  getProjectTaskList(condition.value)
    .then((res) => {
      if (res) {
        list.value = res.items || []
        if (props.useTree) {
          list.value.forEach((item) => {
            item.isLeaf = !item.childrenCount || item.childrenCount === 0
          })
        }
        paginationReactive.itemCount = res.total
        paginationReactive.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
        paginationReactive.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
        statusColumn.filterOptionValue = condition.value.status || 0
        priorityColumn.filterOptionValue = condition.value.priority || 0
      }
    })
    .finally(() => {
      loading.value = false
    })
}

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => projectStore.hasResourceCode("project:detail:task:instance"),
  renderExpand: (row: ProjectTask) => {
    if (!row.taskDesc) {
      getProjectTask(row.id).then((res) => {
        if (res.taskDesc) {
          row.taskDesc = res.taskDesc
        } else {
          row.taskDesc = "暂无描述"
        }
      })
    }
    return h("div", {}, { default: () => row.taskDesc })
  }
})
const nameColumn = reactive({
  title: "任务名称",
  key: "name",
  render: (row: ProjectTask) => {
    const group: VNode[] = []
    group.push(
      h(
        NBadge,
        {
          value: row.parentId ? "子" : "主",
          offset: [10, 0],
          type: row.parentId ? "info" : "success"
        },
        {
          default: () => row.name
        }
      )
    )
    if (row.issue) {
      group.push(
        h(
          NTooltip,
          {},
          {
            default: () => "相关缺陷:" + row.issue?.title || "",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  type: row.issue?.status && row.issue.status >= 5 ? "default" : "error",
                  onClick: () => {
                    router.push(`/project/detail/${row.projectId}/issue?id=${row.issue?.id}`)
                  }
                },
                {
                  default: () => h(NIcon, { component: Debug })
                }
              )
          }
        )
      )
    }
    return h(NSpace, { justify: "space-between" }, () => group)
  }
})
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectTask) => {
    let msg = "未知"
    let hasMe = false
    let imConfirmed = false
    let imStarted = false
    let hasConfirmedMember = false
    if (row.members) {
      for (let i = 0; i < row.members.length; i++) {
        const member = row.members[i]
        if (member.status && member.status === 2) {
          hasConfirmedMember = true
        }
        if (member.userId === userStore.user.id) {
          hasMe = true
          imConfirmed = !!member.status && member.status === 2
          imStarted = !!member.status && member.status === 3
        }
        if (hasMe && hasConfirmedMember) {
          break
        }
      }
    }

    switch (row.status) {
      case -1:
        msg = "已取消"
        break
      case 1:
        msg = "未开始"
        if (hasConfirmedMember) {
          msg = "部分确认"
        }
        break
      case 2:
        msg = "已确认"
        break
      case 3:
        msg = "进行中"
        break
      case 4:
        msg = "已提测"
        break
      case 5:
        msg = "测试打回"
        break
      case 6:
        msg = "测试中"
        break
      case 7:
        msg = "测试通过"
        break
      case 9:
        msg = "已完成"
        break
      default:
        msg = "未知"
    }

    const group: VNode[] = []
    group.push(h(NText, {}, () => msg))
    const btnGroup: VNode[] = []
    if (hasMe) {
      if (projectStore.hasResourceCodes(["project:detail:task:dev", "project:detail:task:test"])) {
        if (row.status === 1) {
          if (!imConfirmed) {
            // 显示确认按钮
            btnGroup.push(
              h(
                NTooltip,
                {},
                {
                  default: () => "确认签收",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "success",
                        onClick: () => {
                          currentTask.value = row
                          const member = row.members?.find((m) => m.userId === userStore.user.id)
                          if (member) {
                            currentTaskMember.value = member
                            showWorkTimeDrawer.value = true
                          }
                        }
                      },
                      {
                        default: () => h(NIcon, { component: BoxCheckmark20Regular })
                      }
                    )
                }
              )
            )
          }
        }
        if (projectStore.hasResourceCode("project:detail:task:dev")) {
          if (row.status === 2) {
            // 都确认后，可以开始执行任务
            btnGroup.push(
              h(
                NTooltip,
                {},
                {
                  default: () => "开始执行",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "success",
                        onClick: () => {
                          startProjectTask(row.id).then(() => {
                            refresh()
                          })
                        }
                      },
                      {
                        default: () => h(NIcon, { component: Play20Regular })
                      }
                    )
                }
              )
            )
          } else if (row.status === 3 || row.status === 5) {
            // 进行中， 可能我已经开始任务或者还未开始任务
            if (imConfirmed) {
              btnGroup.push(
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "开始执行",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "info",
                          onClick: () => {
                            startProjectTask(row.id).then((res) => {
                              if (res) {
                                message.success("开始执行成功")
                                refresh()
                              }
                            })
                          }
                        },
                        {
                          default: () => h(NIcon, { component: Play20Regular })
                        }
                      )
                  }
                )
              )
            } else if (imStarted || row.status === 5) {
              // 可以完成任务，需要填报工时
              btnGroup.push(
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "完成开发并提测",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "success",
                          onClick: () => {
                            // 需要填报工时
                            currentTask.value = row
                            const member = row.members?.find((m) => m.userId === userStore.user.id)
                            if (member) {
                              currentTaskMember.value = member
                              showWorkTimeDrawer.value = true
                            }
                          }
                        },
                        {
                          default: () => h(NIcon, { component: Checkmark20Regular })
                        }
                      )
                  }
                )
              )
            }
          }
        }
        if (projectStore.hasResourceCode("project:detail:task:test")) {
          if (row.status === 4) {
            // 可以开始测试
            btnGroup.push(
              h(
                NTooltip,
                {},
                {
                  default: () => "开始测试",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "success",
                        onClick: () => {
                          testingProjectTask(row.id).then(() => {
                            refresh()
                          })
                        }
                      },
                      {
                        default: () => h(NIcon, { component: Play20Regular })
                      }
                    )
                }
              )
            )
          } else if (row.status === 6) {
            // 测试中，可以拒绝或通过
            btnGroup.push(
              h(
                NTooltip,
                {},
                {
                  default: () => "测试通过",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "success",
                        onClick: () => {
                          testPassProjectTask(row.id).then(() => {
                            refresh()
                          })
                        }
                      },
                      {
                        default: () => h(NIcon, { component: FaceSatisfied })
                      }
                    )
                }
              )
            )
            if (!row.issue) {
              btnGroup.push(
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "测试打回",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "error",
                          onClick: () => {
                            issueInstance.value.title = ""
                            issueInstance.value.projectId = row.projectId
                            issueInstance.value.content = ""
                            issueInstance.value.issueCause = ""
                            issueInstance.value.taskId = row.id
                            taskRejectDrawerActive.value = true
                          }
                        },
                        {
                          default: () => h(NIcon, { component: FaceDissatisfied })
                        }
                      )
                  }
                )
              )
            } else {
              btnGroup.push(
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "请到相关缺陷处进行操作",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "error",
                          disabled: true
                        },
                        {
                          default: () => h(NIcon, { component: FaceDissatisfied })
                        }
                      )
                  }
                )
              )
            }
          }
        }
        if (projectStore.hasResourceCode("project:detail:task:done") && row.status === 7) {
          // 测试通过，可以完成任务
          btnGroup.push(
            h(
              NTooltip,
              {},
              {
                default: () => "完成任务",
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: "small",
                      type: "success",
                      onClick: () => {
                        finishProjectTask(row.id).then((res) => {
                          if (res) {
                            message.success("任务成功完成")
                            refresh()
                          }
                        })
                      }
                    },
                    {
                      default: () => h(NIcon, { component: Checkmark20Regular })
                    }
                  )
              }
            )
          )
        }
      }
    }
    group.push(
      h(
        NButtonGroup,
        {
          size: "small"
        },
        () => btnGroup
      )
    )

    return h(NSpace, { justify: "space-between" }, () => group)
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "已取消",
      value: -1
    },
    {
      label: "未开始",
      value: 1
    },
    {
      label: "已确认",
      value: 2
    },
    {
      label: "进行中",
      value: 3
    },
    {
      label: "已提测",
      value: 4
    },
    {
      label: "测试打回",
      value: 5
    },
    {
      label: "测试中",
      value: 6
    },
    {
      label: "测试通过",
      value: 7
    },
    {
      label: "已完成",
      value: 9
    }
  ]
})
const priorityColumn = reactive({
  title: "优先级",
  key: "priority",
  render: (row: ProjectTask) => {
    switch (row.priority) {
      case 0:
        return "未评估"
      case 1:
        return "低"
      case 2:
        return "中"
      case 3:
        return "高"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "未评估",
      value: 0
    },
    {
      label: "低",
      value: 1
    },
    {
      label: "中",
      value: 2
    },
    {
      label: "高",
      value: 3
    }
  ]
})
const membersColumn = reactive({
  title: "参与人",
  key: "members",
  render: (row: ProjectTask) => {
    return h(
      NAvatarGroup,
      {
        size: 20,
        max: 5,
        options: row.members?.map((item) => {
          return {
            name: item.realName || "",
            src: item.realName || ""
          }
        })
      },
      {
        avatar: (info: { option: { src: string } }) =>
          h(
            NTooltip,
            {},
            {
              trigger: () =>
                h(NameAvatar, {
                  name: info.option.src
                }),
              default: () => info.option.src
            }
          ),
        rest: (info: { options: Array<{ src: string }>; rest: number }) =>
          h(
            NDropdown,
            {
              options: info.options.map((option) => ({
                key: option.src,
                label: option.src
              }))
            },
            {
              default: () =>
                h(
                  NAvatar,
                  {},
                  {
                    default: () => `+${info.rest}`
                  }
                )
            }
          )
      }
    )
  }
})
const ownerColumn = reactive({
  title: "负责人",
  key: "owner",
  render: (row: ProjectTask) => {
    if (row.owner && row.owner.realName) {
      return h(
        NTooltip,
        {},
        {
          trigger: () =>
            h(NameAvatar, {
              size: 20,
              name: row.owner?.realName || ""
            }),
          default: () => row.owner?.realName || ""
        }
      )
    }
  }
})
const actionColumn = reactive({
  title: "操作",
  key: "operation",
  // 返回VNode, 用于渲染操作按钮
  render: (row: ProjectTask) => {
    const bg: VNode[] = []
    if (props.useTree && projectStore.hasResourceCode("project:detail:task:instance")) {
      // 任务树，可以查看详情抽屉
      bg.push(
        h(
          NTooltip,
          {},
          {
            default: () => "查看详情",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  onClick: () => emit("showDetail", row)
                },
                {
                  default: () => h(NIcon, { component: Info20Regular })
                }
              )
          }
        )
      )
    }
    if (condition.value.onlyParent && row.childrenCount && row.childrenCount > 0) {
      bg.push(
        h(
          NTooltip,
          {},
          {
            default: () => "查看子任务",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  onClick: () => emit("showChild", row)
                },
                {
                  default: () => h(NIcon, { component: ParentChild })
                }
              )
          }
        )
      )
    }
    if (row.status !== 9 && projectStore.hasResourceCode("project:detail:task:add")) {
      bg.push(
        h(
          NTooltip,
          {},
          {
            default: () => "任务分解",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  onClick: () => emit("newChild", row)
                },
                {
                  default: () => h(NIcon, { component: ArrowsSplit })
                }
              )
          }
        )
      )
    }
    if (projectStore.hasResourceCode("project:detail:task:edit")) {
      bg.push(
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
                  onClick: () => emit("edit", row)
                },
                {
                  default: () => h(NIcon, { component: Edit })
                }
              )
          }
        )
      )
    }
    if (projectStore.hasResourceCode("project:detail:task:delete")) {
      bg.push(
        h(
          NPopconfirm,
          {
            onPositiveClick: () => emit("delete", row.id)
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
                        type: "error"
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

    return h(NButtonGroup, {}, () => bg)
  }
})
const estimateTimeColumn = reactive({
  title: "预计时间",
  key: "estimateTime",
  render: (row: ProjectTask) => {
    let content = "暂无预计时间"
    let tip = "暂无预计时间"

    if (row.startTime) {
      content = dayjs(row.startTime).format("YYYY-MM-DD") + " ~ " + (row.endTime ? dayjs(row.endTime).format("YYYY-MM-DD") : "?")
      tip = dayjs(row.startTime).format("YYYY-MM-DD HH:mm:ss") + " ~ " + (row.endTime ? dayjs(row.endTime).format("YYYY-MM-DD HH:mm:ss") : "未指定")
    } else if (row.endTime) {
      content = "? ~ " + dayjs(row.endTime).format("YYYY-MM-DD")
      tip = "截止时间: " + dayjs(row.endTime).format("YYYY-MM-DD HH:mm:ss")
    }
    return h(
      NTooltip,
      {},
      {
        default: () => tip,
        trigger: () => content
      }
    )
  }
})
const actualTimeColumn = reactive({
  title: "实际时间",
  key: "actualTime",
  render: (row: ProjectTask) => {
    let content = "暂无实际时间"
    let tip = "暂无实际时间"

    if (row.actualStartTime) {
      content = dayjs(row.actualStartTime).format("YYYY-MM-DD") + " ~ " + (row.actualEndTime ? dayjs(row.actualEndTime).format("YYYY-MM-DD") : "?")
      tip =
        dayjs(row.actualStartTime).format("YYYY-MM-DD HH:mm:ss") +
        " ~ " +
        (row.actualEndTime ? dayjs(row.actualEndTime).format("YYYY-MM-DD HH:mm:ss") : "进行中")
    } else if (row.actualEndTime) {
      content = "? ~ " + dayjs(row.actualEndTime).format("YYYY-MM-DD")
      tip = "截止时间: " + dayjs(row.actualEndTime).format("YYYY-MM-DD HH:mm:ss")
    }
    return h(
      NTooltip,
      {},
      {
        default: () => tip,
        trigger: () => content
      }
    )
  }
})
const columns = computed(() => {
  const result = []
  if (!props.useTree) {
    result.push(expandColumn)
  }
  result.push(nameColumn, priorityColumn, estimateTimeColumn, actualTimeColumn, ownerColumn, membersColumn, statusColumn, actionColumn)
  return result
})

const handleFiltersChange = (filters: DataTableFilterState, sourceColumn: DataTableBaseColumn) => {
  if (!loading.value) {
    switch (sourceColumn.key) {
      case "status":
        if (filters["status"] instanceof Array) {
          statusColumn.filterOptionValue = (filters["status"][0] as number) || 0
        } else {
          statusColumn.filterOptionValue = filters["status"] as number
        }
        condition.value.status = statusColumn.filterOptionValue
        break
      case "priority":
        if (filters["priority"] instanceof Array) {
          priorityColumn.filterOptionValue = (filters["priority"][0] as number) || 0
        } else {
          priorityColumn.filterOptionValue = filters["priority"] as number
        }
        condition.value.priority = priorityColumn.filterOptionValue
        break
    }

    refresh()
  }
}

const onLoad = (row: ProjectTask) => {
  if (props.useTree) {
    return new Promise<void>((resolve, reject) => {
      getProjectTaskList({
        projectId: row.projectId,
        parentId: row.id,
        offset: -1,
        limit: -1
      })
        .then((res) => {
          row.children = res.items
          resolve()
        })
        .catch(() => {
          reject()
        })
    })
  } else {
    return Promise.reject()
  }
}

// 工时抽屉
const showWorkTimeDrawer = ref(false)
const currentTask = ref<ProjectTask>({ id: "", projectId: "", name: "" })
const currentTaskMember = ref<ProjectTaskMember>({ id: "", projectId: "", taskId: "", userId: "" })

// 测试打回时需要建立缺陷，抽屉
const taskRejectDrawerActive = ref(false)
const issueInstance = ref<ProjectIssue>({
  id: "",
  projectId: props.condition.projectId,
  title: "",
  type: 1,
  content: "",
  issueCause: ""
})
const issueFormRef = ref<typeof IssueForm>()
const submitNewIssue = (e: MouseEvent) => {
  e.preventDefault()
  const fi = issueFormRef.value?.formRef as FormInst
  if (!fi) {
    return
  }
  fi.validate((errors) => {
    if (!errors) {
      addProjectIssue(issueInstance.value).then((res) => {
        if (res) {
          message.success("保存成功")
          taskRejectDrawerActive.value = false
          refresh()
        }
      })
    }
  })
}

const refreshIfNoData = () => {
  if (list.value.length === 0) {
    refresh()
  }
}
defineExpose({
  refresh,
  refreshIfNoData
})

// onMounted(() => {
//   refresh()
// })
</script>

<template>
  <n-data-table
    size="small"
    remote
    :data="list"
    :loading="loading"
    :columns="columns"
    :pagination="paginationReactive"
    :row-key="(row: ProjectTask) => row.id"
    :on-update:page="handlePageChange"
    :on-update:page-size="handlePageSizeChange"
    :on-update:filters="handleFiltersChange"
    :on-load="onLoad"
  />

  <work-time-drawer v-model:drawer-active="showWorkTimeDrawer" :task="currentTask" :data="currentTaskMember" @refresh="refresh" />

  <n-drawer v-model:show="taskRejectDrawerActive" :default-height="600" resizable placement="bottom">
    <n-drawer-content>
      <template #header>
        <n-text>新增缺陷</n-text>
        <n-button
          class="absolute right-8px mt--4px"
          type="primary"
          size="small"
          @click="submitNewIssue"
          v-if="projectStore.hasResourceCodes(['project:detail:issue:add', 'project:detail:issue:update'])"
          >提交</n-button
        >
      </template>
      <issue-form ref="issueFormRef" v-model:value="issueInstance" />
    </n-drawer-content>
  </n-drawer>
</template>
