<script setup lang="ts">
import { getMyProjectTaskList, getProjectTask, getProjectTaskList, startProjectTask } from "@/service/api"
import { ProjectTask, ProjectTaskCondition, ProjectTaskMember } from "@/types/project"
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
  NText
} from "naive-ui"
import { ArrowsSplit } from "@vicons/tabler"
import { Delete, Edit, ParentChild } from "@vicons/carbon"
import NameAvatar from "@/components/NameAvatar.vue"
import { useUserStore } from "@/store/user"
import { useProjectStore } from "@/store/project"
import { BoxCheckmark20Regular, Checkmark20Regular, Play20Regular } from "@vicons/fluent"
import dayjs from "dayjs"
import WorkTimeDrawer from "@/views/project/detail/task/workTimeDrawer/index.vue"
const userStore = useUserStore()
const projectStore = useProjectStore()

const props = defineProps<{
  selectedProjectId: string
}>()
const emit = defineEmits(["showChild", "newChild", "edit", "delete", "showDetail"])

const condition = ref<ProjectTaskCondition>({
  projectId: props.selectedProjectId,
  offset: 0,
  limit: 10
})
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

const onlyMine = ref(true)
const handleOnlyMineSwitch = (value: boolean) => {
  onlyMine.value = value
  refresh()
}
const refresh = () => {
  loading.value = true
  if (onlyMine.value) {
    getMyProjectTaskList(condition.value)
      .then((res) => {
        if (res) {
          list.value = res.items || []
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
  } else {
    getProjectTaskList(condition.value)
      .then((res) => {
        if (res) {
          list.value = res.items || []
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
    return h(
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
      case 9:
        msg = "已完成"
        break
      default:
        msg = "未知"
    }

    const group: VNode[] = []
    group.push(h(NText, {}, () => msg))
    if (hasMe) {
      if (row.status === 1) {
        if (!imConfirmed) {
          // 显示确认按钮
          group.push(
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
      } else if (row.status === 2) {
        // 都确认后，可以开始执行任务
        group.push(
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
      } else if (row.status === 3) {
        // 进行中， 可能我已经开始任务或者还未开始任务
        if (imConfirmed) {
          group.push(
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
        } else if (imStarted) {
          // 可以完成任务，需要填报工时
          group.push(
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
const actionColumn = reactive({
  title: "操作",
  key: "operation",
  // 返回VNode, 用于渲染操作按钮
  render: (row: ProjectTask) => {
    const bg: VNode[] = []
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
  result.push(expandColumn)
  result.push(nameColumn, priorityColumn, estimateTimeColumn, actualTimeColumn, membersColumn, statusColumn, actionColumn)
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

// 工时抽屉
const showWorkTimeDrawer = ref(false)
const currentTask = ref<ProjectTask>({ id: "", projectId: "", name: "" })
const currentTaskMember = ref<ProjectTaskMember>({ id: "", projectId: "", taskId: "", userId: "" })

defineExpose({
  projectSelected: (id: string | undefined) => {
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
      <n-space>
        <n-switch :value="onlyMine" @update:value="handleOnlyMineSwitch" size="small">
          <template #checked>只看我的</template>
          <template #unchecked>所有可见的</template>
        </n-switch>
      </n-space>
    </n-gi>
    <n-gi>
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
      />
    </n-gi>
  </n-grid>

  <work-time-drawer v-model:drawer-active="showWorkTimeDrawer" :task="currentTask" :data="currentTaskMember" @refresh="refresh" />
</template>
