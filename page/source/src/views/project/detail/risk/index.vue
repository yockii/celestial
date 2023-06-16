<script setup lang="ts">
import { ProjectRisk, ProjectRiskCondition } from "@/types/project"
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
import { addProjectRisk, deleteProjectRisk, getProjectRisk, getProjectRiskList, updateProjectRisk } from "@/service/api"
import { useStageStore } from "@/store/stage"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import { useUserStore } from "@/store/user"
import { Delete, Edit } from "@vicons/carbon"

const { project } = storeToRefs(useProjectStore())
const userStore = useUserStore()
const stageStore = useStageStore()
const { stageListWithNone } = storeToRefs(stageStore)

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => userStore.hasResourceCode("project:detail:risk:instance"),
  renderExpand: (row: ProjectRisk) => {
    if (!row.riskDesc || !row.response || !row.result) {
      getProjectRisk(row.id).then((res) => {
        row.riskDesc = res.riskDesc
        row.result = res.result
        row.response = res.response
      })
    }
    return h(
      NGrid,
      {
        cols: 1,
        yGap: 8
      },
      [
        h(NGridItem, {}, { default: () => row.riskDesc }),
        h(NGridItem, {}, { default: () => "应对措施：" + (row.response || "暂无") }),
        h(NGridItem, {}, { default: () => "总结：" + (row.result || "暂无") })
      ]
    )
  }
})
const startTimeColumn = reactive({
  title: "开始时间",
  key: "startTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectRisk) => dayjs(row.startTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const endTimeColumn = reactive({
  title: "结束时间",
  key: "endTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectRisk) => dayjs(row.endTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const probabilityColumn = reactive({
  title: "概率",
  key: "probability",
  render: (row: ProjectRisk) => {
    switch (row.riskProbability) {
      case 1:
        return "低"
      case 2:
        return "中"
      case 3:
        return "高"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
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
const impactColumn = reactive({
  title: "影响",
  key: "impact",
  render: (row: ProjectRisk) => {
    switch (row.riskImpact) {
      case 1:
        return "低"
      case 2:
        return "中"
      case 3:
        return "高"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
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
const levelColumn = reactive({
  title: "级别",
  key: "level",
  render: (row: ProjectRisk) => {
    switch (row.riskLevel) {
      case 1:
        return "低"
      case 2:
        return "中"
      case 3:
        return "高"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
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
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectRisk) => {
    switch (row.status) {
      case 1:
        return "已识别"
      case 2:
        return "已应对"
      case 3:
        return "已发生"
      case 4:
        return "已解决"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "已识别",
      value: 1
    },
    {
      label: "已应对",
      value: 2
    },
    {
      label: "已发生",
      value: 3
    },
    {
      label: "已解决",
      value: 4
    }
  ]
})
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectRisk) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const operationColumn = reactive({
  title: "操作",
  key: "operation",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectRisk) => {
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
                disabled: !userStore.hasResourceCode("project:detail:risk:update"),
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
                      disabled: !userStore.hasResourceCode("project:detail:risk:delete")
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
    key: "riskName"
  },
  probabilityColumn,
  impactColumn,
  levelColumn,
  startTimeColumn,
  endTimeColumn,
  statusColumn,
  createTimeColumn,
  operationColumn
]
const condition = ref<ProjectRiskCondition>({
  projectId: project.value.id
})
const list = ref<ProjectRisk[]>([])
const loading = ref(false)

const refresh = () => {
  getProjectRiskList(condition.value).then((res) => {
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
const instance = ref<ProjectRisk>({ endTime: 0, id: "", riskName: "", projectId: "", startTime: 0, status: 1 })
const drawerActive = ref(false)
const isUpdate = computed(() => !!instance.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "编辑风险" : "新建风险"))
const resetInstance = (origin: ProjectRisk | undefined = undefined) => {
  if (origin) {
    instance.value = JSON.parse(JSON.stringify(origin))
  } else {
    instance.value = { endTime: 0, id: "", riskName: "", projectId: project.value.id, startTime: 0, status: 1 }
  }
}
const newInstance = () => {
  resetInstance()
  drawerActive.value = true
}
const riskRules = {
  riskName: [
    { required: true, message: "请输入风险名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  endTime: [
    {
      validator: (rule: FormItemRule, value: number): boolean => !(value && instance.value.startTime && value < instance.value.startTime),
      message: "结束时间不能小于开始时间",
      trigger: "blur"
    }
  ],
  status: {
    required: true,
    validator: (rule: FormItemRule, value: number) => value && value > 0 && value <= 4,
    message: "请选择状态",
    trigger: ["blur", "change"]
  }
}
const handleEditData = (row: ProjectRisk) => {
  instance.value = JSON.parse(JSON.stringify(row))
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectRisk(id).then((res) => {
    if (res) {
      useMessage().success("删除成功")
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
        updateProjectRisk(instance.value).then((res) => {
          if (res) {
            message.success("保存成功")
            drawerActive.value = false
            refresh()
          }
        })
      } else {
        addProjectRisk(instance.value).then((res) => {
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
const route = useRoute()
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
onBeforeUpdate(() => {
  reload()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-space justify="space-between">
        <span></span>
        <n-button type="primary" @click="newInstance" v-resource-code="'project:detail:risk:add'">新建风险</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: ProjectRisk) => row.id"
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
          v-resource-code="['project:detail:risk:add', 'project:detail:risk:update']"
          >提交</n-button
        >
      </template>
      <n-form ref="formRef" :model="instance" :rules="riskRules" label-width="120px" label-placement="left">
        <n-grid :cols="4" x-gap="4">
          <n-gi>
            <n-form-item label="风险名称：" path="riskName">
              <n-input v-model:value="instance.riskName" placeholder="请输入风险名称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="所处阶段：">
              <n-select v-model:value="instance.stageId" placeholder="请选择项目阶段" :options="stageListWithNone" label-field="name" value-field="id" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="开始时间：">
              <n-date-picker type="date" clearable v-model:value="instance.startTime" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="结束时间：" path="endTime">
              <n-date-picker type="date" clearable v-model:value="instance.endTime" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="概率：">
              <n-select
                v-model:value="instance.riskProbability"
                placeholder="请选择概率"
                :options="[
                  {
                    label: '低',
                    value: 1
                  },
                  {
                    label: '中',
                    value: 2
                  },
                  {
                    label: '高',
                    value: 3
                  }
                ]"
                label-field="label"
                value-field="value"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="影响：">
              <n-select
                v-model:value="instance.riskImpact"
                placeholder="请选择影响"
                :options="[
                  {
                    label: '低',
                    value: 1
                  },
                  {
                    label: '中',
                    value: 2
                  },
                  {
                    label: '高',
                    value: 3
                  }
                ]"
                label-field="label"
                value-field="value"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="等级：">
              <n-select
                v-model:value="instance.riskLevel"
                placeholder="请选择等级"
                :options="[
                  {
                    label: '低',
                    value: 1
                  },
                  {
                    label: '中',
                    value: 2
                  },
                  {
                    label: '高',
                    value: 3
                  }
                ]"
                label-field="label"
                value-field="value"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="状态：" path="status">
              <n-select
                v-model:value="instance.status"
                placeholder="请选择状态"
                :options="[
                  {
                    label: '已识别',
                    value: 1
                  },
                  {
                    label: '已应对',
                    value: 2
                  },
                  {
                    label: '已发生',
                    value: 3
                  },
                  {
                    label: '已解决',
                    value: 4
                  }
                ]"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="风险描述：">
              <n-input v-model:value="instance.riskDesc" placeholder="请输入风险描述" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="应对措施：">
              <n-input
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                v-model:value="instance.response"
                placeholder="请输入应对措施，即本风险的应对手段"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="结果总结：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.result" placeholder="请输入结果总结" />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
