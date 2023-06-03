<script setup lang="ts">
import { getProjectModuleList } from "@/service/api/projectModule"
import { deleteProjectRequirement, getProjectRequirementList } from "@/service/api/projectRequirement"
import { useProjectStore } from "@/store/project"
import { ProjectModule, ProjectRequirement, ProjectRequirementCondition, ProjectTask } from "@/types/project"
import { useMessage, NButton, NButtonGroup, NPopconfirm, PaginationProps, DataTableFilterState, DataTableBaseColumn } from "naive-ui"
import { storeToRefs } from "pinia"
import { Refresh } from "@vicons/tabler"
import Drawer from "./drawer/index.vue"
import TaskDrawer from "../task/drawer/index.vue"

const message = useMessage()
const projectStore = useProjectStore()
const { project } = storeToRefs(useProjectStore())

const condition = ref<ProjectRequirementCondition>({
  projectId: project.value.id
})
const list = ref<ProjectRequirement[]>([])

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
  getProjectRequirementList(condition.value)
    .then((res) => {
      if (res) {
        list.value = res.items
        paginationReactive.itemCount = res.total
        paginationReactive.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
        paginationReactive.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
        statusColumn.filterOptionValue = condition.value.status || 0
        feasibilityColumn.filterOptionValue = condition.value.feasibility || 0
        typeColumn.filterOptionValue = condition.value.type || 0
      }
    })
    .finally(() => {
      loading.value = false
    })
}
const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: (row: ProjectRequirement) => row.detail && row.detail !== "",
  renderExpand: (row: ProjectRequirement) => h("div", {}, { default: () => row.detail })
})
const typeColumn = reactive({
  title: "类型",
  key: "type",
  render: (row: ProjectRequirement) => {
    switch (row.type) {
      case 1:
        return "功能"
      case 2:
        return "接口"
      case 3:
        return "性能"
      case 4:
        return "安全"
      case 5:
        return "体验"
      case 6:
        return "改进"
      case 7:
        return "其他"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "功能",
      value: 1
    },
    {
      label: "接口",
      value: 2
    },
    {
      label: "性能",
      value: 3
    },
    {
      label: "安全",
      value: 4
    },
    {
      label: "体验",
      value: 5
    },
    {
      label: "改进",
      value: 6
    },
    {
      label: "其他",
      value: 7
    }
  ]
})
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectRequirement) => {
    switch (row.status) {
      case 1:
        return "待评审"
      case 2:
        return "已评审"
      case 3:
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
      label: "待评审",
      value: 1
    },
    {
      label: "已评审",
      value: 2
    },
    {
      label: "已完成",
      value: 3
    }
  ]
})
const feasibilityColumn = reactive({
  title: "可行性",
  key: "feasibility",
  render: (row: ProjectRequirement) => {
    switch (row.feasibility) {
      case -1:
        return "不可行"
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
      label: "不可行",
      value: -1
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
const columns = [
  expandColumn,
  {
    title: "需求名称",
    key: "name"
  },
  typeColumn,
  feasibilityColumn,
  statusColumn,
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: ProjectRequirement) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            onClick: () => {
              handleAddProjectTask(row)
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
      case "feasibility":
        if (filters["feasibility"] instanceof Array) {
          feasibilityColumn.filterOptionValue = (filters["feasibility"][0] as number) || 0
        } else {
          feasibilityColumn.filterOptionValue = filters["feasibility"] as number
        }
        condition.value.feasibility = feasibilityColumn.filterOptionValue
        break
      case "type":
        if (filters["type"] instanceof Array) {
          typeColumn.filterOptionValue = (filters["type"][0] as number) || 0
        } else {
          typeColumn.filterOptionValue = filters["type"] as number
        }
        condition.value.type = typeColumn.filterOptionValue
        break
    }

    refresh()
  }
}

const drawerActive = ref(false)
const currentData = ref<ProjectRequirement>({
  id: "",
  name: "",
  projectId: project.value.id
})

const handleEditData = (row: ProjectRequirement) => {
  currentData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectRequirement(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
const handleAddProjectRequirement = () => {
  currentData.value = {
    id: "",
    name: "",
    projectId: project.value.id
  }
  drawerActive.value = true
}

// 需求任务分解
const taskDrawerActive = ref(false)
const currentTaskData = ref<ProjectTask>({
  id: "",
  name: "",
  projectId: ""
})
const handleAddProjectTask = (requirement: ProjectRequirement) => {
  currentTaskData.value = {
    id: "",
    name: "",
    projectId: project.value.id,
    moduleId: requirement.moduleId,
    requirementId: requirement.id
  }
  taskDrawerActive.value = true
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
      modules.value = res.items
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
          <n-button type="primary" @click="handleAddProjectRequirement">新增需求</n-button>
        </n-gi>
        <n-gi>
          <n-data-table
            size="small"
            remote
            :data="list"
            :loading="loading"
            :columns="columns"
            :pagination="paginationReactive"
            :row-key="(row: ProjectRequirement) => row.id"
            :on-update:page="handlePageChange"
            :on-update:page-size="handlePageSizeChange"
            :on-update:filters="handleFiltersChange"
          />
        </n-gi>
      </n-grid>
    </n-gi>
  </n-grid>

  <drawer :project-id="project.id" v-model:drawer-active="drawerActive" v-model:data="currentData" @refresh="refresh" />

  <task-drawer :project-id="project.id" v-model:drawer-active="taskDrawerActive" v-model:data="currentTaskData" />
</template>

<style scoped></style>
