<script setup lang="ts">
import { getProjectTask, getProjectTaskList } from "@/service/api/project/projectTask"
import { ProjectTask, ProjectTaskCondition } from "@/types/project"
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
  DataTableColumn
} from "naive-ui"
import { ArrowsSplit } from "@vicons/tabler"
import { Delete, Edit, ParentChild } from "@vicons/carbon"
import NameAvatar from "@/components/NameAvatar.vue"
import { useUserStore } from "@/store/user"
const userStore = useUserStore()

const props = defineProps<{
  condition: ProjectTaskCondition
  useTree?: boolean
}>()
const emit = defineEmits(["showChild", "newChild", "edit", "delete"])

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

const expandColumn = reactive<DataTableColumn>({
  key: "expand",
  type: "expand",
  expandable: () => userStore.hasResourceCode("project:detail:task:instance"),
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
    switch (row.status) {
      case -1:
        return "已取消"
      case 1:
        return "未开始"
      case 2:
        return "已确认"
      case 3:
        return "进行中"
      case 9:
        return "已完成"
      default:
        return "未知"
    }
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
          console.log
          return {
            name: item.realName,
            src: item.realName
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
    if (condition.value.onlyParent && row.childrenCount > 0) {
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
    if (userStore.hasResourceCode("project:detail:task:add")) {
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
    if (userStore.hasResourceCode("project:detail:task:edit")) {
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
    if (userStore.hasResourceCode("project:detail:task:delete")) {
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
const columns = computed(() => {
  const result = []
  if (!props.useTree) {
    result.push(expandColumn)
  }
  result.push(nameColumn, priorityColumn, membersColumn, statusColumn, actionColumn)
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

defineExpose({
  refresh
})

onMounted(() => {
  refresh()
})
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
</template>
