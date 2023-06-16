<script setup lang="ts">
import { ProjectChange, ProjectChangeCondition } from "@/types/project"
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
  NPopconfirm,
  NTooltip,
  PaginationProps
} from "naive-ui"
import { addProjectChange, deleteProjectChange, getProjectChange, getProjectChangeList, updateProjectChange } from "@/service/api"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import { useUserStore } from "@/store/user"
import { Delete, Edit } from "@vicons/carbon"

const message = useMessage()
const { project } = storeToRefs(useProjectStore())
const userStore = useUserStore()

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => userStore.hasResourceCode("project:detail:plan:instance"),
  renderExpand: (row: ProjectChange) => {
    if (!row.reason || !row.plan || !row.review || !row.risk || !row.result) {
      getProjectChange(row.id).then((res) => {
        row.reason = res.reason
        row.plan = res.plan
        row.review = res.review
        row.risk = res.risk
        row.result = res.result
      })
    }
    return h(
      NGrid,
      {
        cols: 1,
        yGap: 8
      },
      [
        h(NGridItem, {}, { default: () => row.reason }),
        h(NGridItem, {}, { default: () => "变更：" + row.plan }),
        h(NGridItem, {}, { default: () => "风险：" + row.risk }),
        h(NGridItem, {}, { default: () => "评审结果：" + row.review }),
        h(NGridItem, {}, { default: () => "结果说明：" + row.result })
      ]
    )
  }
})

const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectChange) => {
    switch (row.status) {
      case -1:
        return "已拒绝"
      case 1:
        return "待评审"
      case 2:
        return "已批准"
      case 9:
        return "关闭"
    }
    return "未知"
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "已拒绝",
      value: -1
    },
    {
      label: "待评审",
      value: 1
    },
    {
      label: "已批准",
      value: 2
    },
    {
      label: "关闭",
      value: 9
    }
  ]
})
const reviewTimeColumn = reactive({
  title: "评审时间",
  key: "reviewTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectChange) => (row.reviewTime ? dayjs(row.reviewTime).fromNow() : "未评审"),
  // 排序
  sorter: true,
  sortOrder: false
})
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectChange) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const operationColumn = reactive({
  title: "操作",
  key: "operation",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectChange) => {
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
                disabled: !userStore.hasResourceCode("project:detail:change:update"),
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
                      disabled: !userStore.hasResourceCode("project:detail:change:delete")
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
    key: "title"
  },
  {
    title: "类型",
    key: "type",
    render: (row: ProjectChange) => {
      switch (row.type) {
        case 1:
          return "时间节点调整"
        case 2:
          return "需求变更"
        case 3:
          return "资源变动"
        case 9:
          return "其他变更"
      }
      return "未知"
    }
  },
  {
    title: "级别",
    key: "level",
    render: (row: ProjectChange) => {
      switch (row.level) {
        case 1:
          return "一般"
        case 2:
          return "重大"
      }
      return "未知"
    }
  },
  statusColumn,
  reviewTimeColumn,
  createTimeColumn,
  operationColumn
]
const condition = ref<ProjectChangeCondition>({
  projectId: project.value.id
})
const list = ref<ProjectChange[]>([])
const loading = ref(false)

const refresh = () => {
  getProjectChangeList(condition.value).then((res) => {
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
    }
    condition.value.orderBy = field + (order === "ascend" ? " asc" : " desc")
    refresh()
  }
}

// 抽屉部分，新建、编辑内容
const instance = ref<ProjectChange>({
  id: "",
  projectId: project.value.id,
  title: ""
})
const drawerActive = ref(false)
const isUpdate = computed(() => !!instance.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "编辑变更" : "新建变更"))
const resetInstance = (origin: ProjectChange | undefined = undefined) => {
  if (origin) {
    instance.value = JSON.parse(JSON.stringify(origin))
  } else {
    instance.value = {
      id: "",
      projectId: project.value.id,
      title: ""
    }
  }
}
const newInstance = () => {
  resetInstance()
  drawerActive.value = true
}
const rules = {
  title: [
    { required: true, message: "请输入变更名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: {
    validator: (rule: unknown, value: number) => {
      if (value && value > 0 && value < 10) {
        return true
      }
      return false
    },
    required: true,
    message: "请选择变更类型",
    trigger: ["blur", "change"]
  }
}
const handleEditData = (row: ProjectChange) => {
  instance.value = JSON.parse(JSON.stringify(row))
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectChange(id).then((res) => {
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
        updateProjectChange(instance.value).then((res) => {
          if (res) {
            message.success("保存成功")
            drawerActive.value = false
            refresh()
          }
        })
      } else {
        addProjectChange(instance.value).then((res) => {
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
        <n-button type="primary" @click="newInstance" v-resource-code="'project:detail:change:add'">新建变更</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: ProjectChange) => row.id"
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
          v-resource-code="['project:detail:change:add', 'project:detail:change:update']"
          >提交</n-button
        >
      </template>
      <n-form ref="formRef" :model="instance" :rules="rules" label-width="120px" label-placement="left">
        <n-grid :cols="4" x-gap="4">
          <n-gi>
            <n-form-item label="变更名称：" path="title">
              <n-input v-model:value="instance.title" placeholder="请输入变更名称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="类型：" path="type">
              <n-select
                v-model:value="instance.type"
                placeholder="请选择变更类型"
                :options="[
                  { label: '时间节点调整', value: 1 },
                  { label: '需求变更', value: 2 },
                  { label: '资源变动', value: 3 },
                  { label: '其他', value: 9 }
                ]"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="级别">
              <n-select
                v-model:value="instance.level"
                placeholder="请选择变更级别"
                :options="[
                  { label: '一般', value: 1 },
                  { label: '重大', value: 2 }
                ]"
              />
            </n-form-item>
          </n-gi>

          <n-gi :span="4">
            <n-form-item label="变更原因：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.reason" placeholder="请输入变更原因" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="方案：">
              <n-input
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                v-model:value="instance.plan"
                placeholder="请输入方案，即本变更的目标和实施方案"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="评审结果：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.review" placeholder="请输入评审结果" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="风险：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.risk" placeholder="请输入风险，及时更新" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="结果说明：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.result" placeholder="请输入结果说明" />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
