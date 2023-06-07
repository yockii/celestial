<script setup lang="ts">
import { getProjectModuleList } from "@/service/api/projectModule"
import { deleteProjectTask, getProjectTask, getProjectTaskList } from "@/service/api/projectTask"
import { useProjectStore } from "@/store/project"
import { ProjectModule, ProjectTask, ProjectTaskCondition } from "@/types/project"
import {
  useMessage,
  NButton,
  NButtonGroup,
  NPopconfirm,
  PaginationProps,
  DataTableFilterState,
  DataTableBaseColumn,
  NBadge,
  NAvatarGroup,
  NTooltip,
  NDropdown,
  NAvatar
} from "naive-ui"
import { storeToRefs } from "pinia"
import { Refresh } from "@vicons/tabler"
import Drawer from "./drawer/index.vue"
import { useUserStore } from "@/store/user"
import NameAvatar from "@/components/NameAvatar.vue"

const message = useMessage()
const userStore = useUserStore()
const projectStore = useProjectStore()
const { project } = storeToRefs(useProjectStore())
const condition = ref<ProjectTaskCondition>({
  projectId: project.value.id,
  onlyParent: false
})
const list = ref<ProjectTask[]>([])

const { modules, moduleTree } = storeToRefs(projectStore)
const treeSelected = (keys: string[], modules: ProjectModule[]) => {
  if (modules.length > 0) {
    condition.value.fullPath = modules[0].fullPath
    refresh()
  }
}
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
const columns = [
  expandColumn,
  nameColumn,
  priorityColumn,
  membersColumn,
  statusColumn,
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: ProjectTask) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            disabled: !userStore.hasResourceCode("project:detail:task:add"),
            onClick: () => {
              console.log(1)
            }
          },
          {
            default: () => "任务分解"
          }
        ),
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            type: "primary",
            disabled: !userStore.hasResourceCode("project:detail:task:edit"),
            onClick: () => handleEditData(row)
          },
          {
            default: () => "编辑"
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteData(row.id)
          },
          {
            default: () => "确认删除",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  disabled: !userStore.hasResourceCode("project:detail:task:delete"),
                  type: "error"
                },
                {
                  default: () => "删除"
                }
              )
          }
        )
      ])
    }
  }
]
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

const drawerActive = ref(false)
const currentData = ref<ProjectTask>({
  id: "",
  name: "",
  projectId: project.value.id
})

const handleEditData = (row: ProjectTask) => {
  currentData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectTask(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
const handleAddProjectTask = () => {
  currentData.value = {
    id: "",
    name: "",
    projectId: project.value.id
  }
  drawerActive.value = true
}

// 加载页面
onMounted(() => {
  refresh()
  // 如果功能模块列表为空, 则加载
  if (!moduleTree.value.length) {
    loadModules()
  }
})
const loadModules = () => {
  getProjectModuleList({
    projectId: project.value.id,
    offset: -1,
    limit: -1
  }).then((res) => {
    if (res) {
      modules.value = res.items || []
    }
  })
}
</script>

<template>
  <n-grid :cols="6" x-gap="16">
    <n-gi>
      <n-grid :cols="1" y-gap="16">
        <n-gi class="flex justify-between">
          <n-text class="font-bold">功能模块</n-text>
          <n-button circle size="small" @click="loadModules">
            <template #icon>
              <n-icon>
                <Refresh />
              </n-icon>
            </template>
          </n-button>
        </n-gi>
        <n-gi>
          <n-tree :data="moduleTree" key-field="id" label-field="name" children-field="children" :on-update:selected-keys="treeSelected" selectable />
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi :span="5">
      <n-grid :cols="1" y-gap="16">
        <n-gi>
          <n-space justify="space-between">
            <n-switch v-model:value="condition.onlyParent" :round="false" @update:value="refresh" size="small" :loading="loading">
              <template #checked>仅显示主任务</template>
              <template #unchecked>显示所有任务</template>
            </n-switch>
            <n-button type="primary" @click="handleAddProjectTask" v-resource-code="'project:detail:task:add'">新增任务</n-button>
          </n-space>
        </n-gi>
        <n-gi v-resource-code="'project:detail:task:list'">
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
    </n-gi>
  </n-grid>

  <drawer :project-id="project.id" v-model:drawer-active="drawerActive" v-model:data="currentData" @refresh="refresh" />
</template>
