<script setup lang="ts">
import { ProjectPlan, ProjectPlanCondition } from "@/types/project"
import { computed, h, reactive, ref } from "vue"
import dayjs from "dayjs"
import {
  DataTableFilterState,
  DataTableSortState,
  FormInst,
  FormItemRule,
  NButton,
  NButtonGroup,
  NGrid,
  NGridItem,
  NIcon,
  NPopconfirm,
  NTooltip,
  PaginationProps
} from "naive-ui"
import { addProjectPlan, deleteProjectPlan, getProjectPlan, getProjectPlanList, updateProjectPlan } from "@/service/api"
import { useStageStore } from "@/store/stage"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import { Delete, Edit } from "@vicons/carbon"
import DatePickerInTable from "@/components/DatePickerInTable.vue"
import StatusInTable from "@/components/StatusInTable.vue"

const message = useMessage()
const projectStore = useProjectStore()
const { project } = storeToRefs(projectStore)
const stageStore = useStageStore()
const { stageListWithNone } = storeToRefs(stageStore)

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => projectStore.hasResourceCode("project:detail:plan:instance"),
  renderExpand: (row: ProjectPlan) => {
    if (!row.planDesc || !row.target || !row.scope || !row.resource || !row.schedule) {
      getProjectPlan(row.id).then((res) => {
        row.planDesc = res.planDesc
        row.target = res.target
        row.scope = res.scope
        row.resource = res.resource
        row.schedule = res.schedule
      })
    }
    return h(
      NGrid,
      {
        cols: 1,
        yGap: 8
      },
      [
        h(NGridItem, {}, { default: () => row.planDesc }),
        h(NGridItem, {}, { default: () => "目标：" + row.target }),
        h(NGridItem, {}, { default: () => "范围：" + row.scope }),
        h(NGridItem, {}, { default: () => "资源：" + row.resource }),
        h(NGridItem, {}, { default: () => "进展：" + row.schedule })
      ]
    )
  }
})
const startTimeColumn = reactive({
  title: "计划开始时间",
  key: "startTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => dayjs(row.startTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const endTimeColumn = reactive({
  title: "计划结束时间",
  key: "endTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => dayjs(row.endTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const actualStartTimeColumn = reactive({
  title: "实际开始时间",
  key: "actualStartTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => {
    return h(DatePickerInTable, {
      data: row.actualStartTime,
      "onUpdate:data": (val: number) => {
        updateProjectPlan({
          id: row.id,
          actualStartTime: val
        }).then((res) => {
          if (res) {
            message.success("更新成功")
            row.actualStartTime = val
          }
        })
      }
    })

    // (row.actualStartTime === 0 ? "无" : dayjs(row.actualStartTime).format("YYYY-MM-DD"))
  },
  // 排序
  sorter: true,
  sortOrder: false
})
const actualEndTimeColumn = reactive({
  title: "实际结束时间",
  key: "actualEndTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => {
    return h(DatePickerInTable, {
      data: row.actualEndTime,
      "onUpdate:data": (val: number) => {
        updateProjectPlan({
          id: row.id,
          actualEndTime: val
        }).then((res) => {
          if (res) {
            message.success("更新成功")
            row.actualEndTime = val
          }
        })
      }
    })
  },
  // 排序
  sorter: true,
  sortOrder: false
})
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectPlan) => {
    return h(StatusInTable, {
      data: row.status,
      options: [
        {
          label: "废弃",
          value: -1
        },
        {
          label: "未开始",
          value: 1
        },
        {
          label: "执行中",
          value: 2
        },
        {
          label: "完成",
          value: 3
        }
      ],
      "onUpdate:data": (val: number) => {
        updateProjectPlan({
          id: row.id,
          status: val
        }).then((res) => {
          if (res) {
            message.success("更新成功")
            row.status = val
          }
        })
      }
    })

    // switch (row.status) {
    //   case -1:
    //     return "废弃"
    //   case 1:
    //     return "未开始"
    //   case 2:
    //     return "进行中"
    //   case 3:
    //     return "完成"
    // }
    // return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "废弃",
      value: -1
    },
    {
      label: "未开始",
      value: 1
    },
    {
      label: "执行中",
      value: 2
    },
    {
      label: "完成",
      value: 3
    }
  ]
})
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const operationColumn = reactive({
  title: "操作",
  key: "operation",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectPlan) => {
    return h(NButtonGroup, {}, () => [
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
                disabled: !projectStore.hasResourceCode("project:detail:plan:update"),
                onClick: () => handleEditData(row)
              },
              {
                default: () => h(NIcon, { component: Edit })
              }
            )
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
                      disabled: !projectStore.hasResourceCode("project:detail:plan:delete")
                    },
                    {
                      default: () => h(NIcon, { component: Delete })
                    }
                  )
              }
            )
        }
      )
    ])
  }
})
const columns = [
  expandColumn,
  {
    title: "名称",
    key: "planName"
  },
  startTimeColumn,
  endTimeColumn,
  actualStartTimeColumn,
  actualEndTimeColumn,
  statusColumn,
  createTimeColumn,
  operationColumn
]
const condition = ref<ProjectPlanCondition>({
  projectId: project.value.id
})
const list = ref<ProjectPlan[]>([])
const loading = ref(false)

const refresh = () => {
  getProjectPlanList(condition.value).then((res) => {
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
    } else if (columnKey === "actualStartTime") {
      actualStartTimeColumn.sortOrder = order === "ascend"
      field = "actual_start_time"
    } else if (columnKey === "actualEndTime") {
      actualEndTimeColumn.sortOrder = order === "ascend"
      field = "actual_end_time"
    }
    condition.value.orderBy = field + (order === "ascend" ? " asc" : " desc")
    refresh()
  }
}

// 抽屉部分，新建、编辑内容
const instance = ref<ProjectPlan>({ endTime: 0, id: "", planName: "", projectId: "", startTime: 0, status: 0 })
const drawerActive = ref(false)
const isUpdate = computed(() => !!instance.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "编辑计划" : "新建计划"))
const resetInstance = (origin: ProjectPlan | undefined = undefined) => {
  if (origin) {
    instance.value = JSON.parse(JSON.stringify(origin))
  } else {
    const nowa = new Date().getTime()
    instance.value = { id: "", planName: "", projectId: project.value.id, startTime: nowa, endTime: nowa, status: 0 }
  }
}
const newInstance = () => {
  resetInstance()
  drawerActive.value = true
}
const planRules = {
  planName: [
    { required: true, message: "请输入计划名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  endTime: [
    {
      validator: (rule: FormItemRule, value: number): boolean => !(value && instance.value.startTime && value < instance.value.startTime),
      message: "结束时间不能小于开始时间",
      trigger: "blur"
    }
  ]
}
const handleEditData = (row: ProjectPlan) => {
  getProjectPlan(row.id).then((res) => {
    if (res) {
      resetInstance(res)
      drawerActive.value = true
    }
  })
}
const handleDeleteData = (id: string) => {
  deleteProjectPlan(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
const formRef = ref<FormInst>()
const submit = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (isUpdate.value) {
        updateProjectPlan(instance.value).then((res) => {
          if (res) {
            message.success("保存成功")
            drawerActive.value = false
            refresh()
          }
        })
      } else {
        addProjectPlan(instance.value).then((res) => {
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

// 加载完毕
const reload = () => {
  if (route.query.id) {
    condition.value = {
      id: route.query.id as string,
      projectId: project.value.id
    }
  }
  refresh()
}
onMounted(() => {
  reload()
})
const route = useRoute()
onBeforeUpdate(() => {
  reload()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-space justify="space-between">
        <span></span>
        <n-button type="primary" @click="newInstance" v-if="projectStore.hasResourceCode('project:detail:plan:add')">新建计划</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: ProjectPlan) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :on-update:sorter="handleSorterChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :default-height="600" resizable placement="bottom">
    <n-drawer-content>
      <template #header>
        <n-text>{{ drawerTitle }}</n-text>
        <n-button
          class="absolute right-8px mt--4px"
          type="primary"
          size="small"
          @click="submit"
          v-if="projectStore.hasResourceCodes(['project:detail:plan:add', 'project:detail:plan:update'])"
          >提交</n-button
        >
      </template>
      <n-form ref="formRef" :model="instance" :rules="planRules" label-width="120px" label-placement="left">
        <n-grid :cols="4" x-gap="4">
          <n-gi>
            <n-form-item label="计划名称：" path="planName">
              <n-input v-model:value="instance.planName" placeholder="请输入计划名称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="所处阶段：">
              <n-select v-model:value="instance.stageId" placeholder="请选择项目阶段" :options="stageListWithNone" label-field="name" value-field="id" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="计划开始时间：">
              <n-date-picker type="date" v-model:value="instance.startTime" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="计划结束时间：" path="endTime">
              <n-date-picker type="date" v-model:value="instance.endTime" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="计划描述：">
              <n-input v-model:value="instance.planDesc" placeholder="请输入计划描述" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="目标：">
              <n-input
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                v-model:value="instance.target"
                placeholder="请输入目标，应包含要做的事情及相应的输出物"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="范围：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.scope" placeholder="请输入范围，即本计划的边界确定" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="资源：">
              <n-input
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                v-model:value="instance.resource"
                placeholder="请输入资源，即要调用的人力、物力等资源信息"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="进展：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.schedule" placeholder="请输入进展，及时更新" />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
