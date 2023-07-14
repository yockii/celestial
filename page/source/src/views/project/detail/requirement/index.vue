<script setup lang="ts">
import {
  getProjectModuleList,
  requirementCompleted,
  deleteProjectRequirement,
  requirementDesigned,
  getProjectRequirement,
  getProjectRequirementList,
  requirementReview
} from "@/service/api"
import { useProjectStore } from "@/store/project"
import { ProjectModule, ProjectRequirement, ProjectRequirementCondition, ProjectTask } from "@/types/project"
import { useMessage, NButton, NButtonGroup, NPopconfirm, PaginationProps, DataTableFilterState, DataTableBaseColumn, NSpace, NIcon, NTooltip } from "naive-ui"
import { storeToRefs } from "pinia"
import { ArrowsSplit, Refresh } from "@vicons/tabler"
import Drawer from "./drawer/index.vue"
import TaskDrawer from "../task/drawer/index.vue"
import { AiStatusComplete, AiStatusFailed, BrushFreehand, Delete, Edit } from "@vicons/carbon"
import { CodeOutlined } from "@vicons/antd"

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
        list.value = res.items || []
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

// 状态修改
const handleDesigned = (row: ProjectRequirement) => {
  requirementDesigned(row.id).then((res) => {
    if (res) {
      message.success(row.name + "设计完成")
      refresh()
    } else {
      message.error(row.name + "状态修改失败")
    }
  })
}
const handleReview = (row: ProjectRequirement, status: number) => {
  requirementReview(row.id, status).then((res) => {
    if (res) {
      message.success(row.name + "评审完成")
      refresh()
    } else {
      message.error(row.name + "状态修改失败")
    }
  })
}
const handleCompleted = (row: ProjectRequirement) => {
  requirementCompleted(row.id).then((res) => {
    if (res) {
      message.success(row.name + "已完成")
      refresh()
    } else {
      message.error(row.name + "状态修改失败")
    }
  })
}

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => projectStore.hasResourceCode("project:detail:requirement:instance"),
  renderExpand: (row: ProjectRequirement) => {
    if (!row.detail) {
      getProjectRequirement(row.id).then((res) => {
        if (res) {
          row.detail = res.detail
        }
      })
    }
    return h("div", {}, { default: () => row.detail })
  }
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
    // 状态 1-待设计 2-待评审 3-评审通过 9-已完成 -1-评审未通过
    switch (row.status) {
      case 1:
        // "待设计"
        return h(
          NSpace,
          { justify: "space-between" },
          {
            default: () => [
              h("span", {}, { default: () => "待设计" }),
              projectStore.hasResourceCode("project:detail:requirement:statusDesign")
                ? h(
                    NTooltip,
                    {},
                    {
                      default: () => "设计完成",
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: "small",
                            type: "info",
                            onClick: () => handleDesigned(row)
                          },
                          { default: () => h(NIcon, { component: BrushFreehand }) }
                        )
                    }
                  )
                : ""
            ]
          }
        )
      case 2:
        // "待评审"
        return h(
          NSpace,
          { justify: "space-between" },
          {
            default: () => [
              h("span", {}, { default: () => "待评审" }),
              projectStore.hasResourceCode("project:detail:requirement:statusReview")
                ? h(
                    NButtonGroup,
                    {},
                    {
                      default: () => [
                        h(
                          NTooltip,
                          {},
                          {
                            default: () => "评审通过",
                            trigger: () =>
                              h(
                                NButton,
                                {
                                  size: "small",
                                  type: "primary",
                                  onClick: () => handleReview(row, 3)
                                },
                                { default: () => h(NIcon, { component: AiStatusComplete }) }
                              )
                          }
                        ),
                        h(
                          NTooltip,
                          {},
                          {
                            default: () => "评审不通过",
                            trigger: () =>
                              h(
                                NButton,
                                {
                                  size: "small",
                                  type: "error",
                                  onClick: () => handleReview(row, -1)
                                },
                                { default: () => h(NIcon, { component: AiStatusFailed }) }
                              )
                          }
                        )
                      ]
                    }
                  )
                : ""
            ]
          }
        )
      case 3:
        // "评审通过"
        return h(
          NSpace,
          { justify: "space-between" },
          {
            default: () => [
              h("span", {}, { default: () => "评审通过" }),
              projectStore.hasResourceCode("project:detail:requirement:statusComplete")
                ? h(
                    NTooltip,
                    {},
                    {
                      default: () => "完成",
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: "small",
                            type: "primary",
                            secondary: true,
                            onClick: () => handleCompleted(row)
                          },
                          { default: () => h(NIcon, { component: CodeOutlined }) }
                        )
                    }
                  )
                : ""
            ]
          }
        )
      case 9:
        return "已完成"
      case -1:
        // "评审未通过", 可以再次设计完成
        return h(
          NSpace,
          { justify: "space-between" },
          {
            default: () => [
              h("span", {}, { default: () => "评审未通过" }),
              projectStore.hasResourceCode("project:detail:requirement:statusDesign")
                ? h(
                    NTooltip,
                    {},
                    {
                      default: () => "设计完成",
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: "small",
                            type: "info",
                            onClick: () => handleDesigned(row)
                          },
                          { default: () => h(NIcon, { component: BrushFreehand }) }
                        )
                    }
                  )
                : ""
            ]
          }
        )
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "待设计",
      value: 1
    },
    {
      label: "待评审",
      value: 2
    },
    {
      label: "评审通过",
      value: 3
    },
    {
      label: "已完成",
      value: 9
    },
    {
      label: "评审未通过",
      value: -1
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
      const btnGroup: VNode[] = []
      if (row.status && row.status > 2 && projectStore.hasResourceCode("project:detail:task:add")) {
        btnGroup.push(
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
                    disabled: !projectStore.hasResourceCode("project:detail:task:add"),
                    onClick: () => {
                      handleAddProjectTask(row)
                    }
                  },
                  {
                    default: () => h(NIcon, { component: ArrowsSplit })
                  }
                )
            }
          )
        )
      }
      if (projectStore.hasResourceCode("project:detail:requirement:edit") && (!row.status || row.status < 3)) {
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
                    disabled: !(projectStore.hasResourceCode("project:detail:requirement:edit") && (!row.status || row.status < 3)),
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
      if ((!row.status || row.status < 3) && projectStore.hasResourceCode("project:detail:requirement:delete")) {
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
                          disabled: !projectStore.hasResourceCode("project:detail:requirement:delete"),
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

      return h(NButtonGroup, {}, () => btnGroup)
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

const reload = () => {
  if (route.query.id) {
    condition.value = {
      id: route.query.id as string,
      projectId: project.value.id
    }
  }
  refresh()

  // 如果功能模块列表为空, 则加载
  if (!moduleTree.value.length) {
    loadModules()
  }
}

onMounted(() => {
  reload()
})

const route = useRoute()
// onBeforeUpdate(() => {
//   reload()
// })
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
        <n-gi v-if="projectStore.hasResourceCode('project:detail:module')">
          <n-tree :data="moduleTree" key-field="id" label-field="name" children-field="children" :on-update:selected-keys="treeSelected" selectable />
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi :span="5">
      <n-grid :cols="1" y-gap="16">
        <n-gi>
          <n-space justify="space-between">
            <span></span>
            <n-button type="primary" @click="handleAddProjectRequirement" v-if="projectStore.hasResourceCode('project:detail:requirement:add')"
              >新增需求</n-button
            >
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
