<script setup lang="ts">
import { ProjectIssue, ProjectIssueCondition } from "@/types/project"
import { computed, h, reactive, ref } from "vue"
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
import { addProjectIssue, deleteProjectIssue, getProjectIssue, getProjectIssueList, updateProjectIssue } from "@/service/api/project/projectIssue"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import { useUserStore } from "@/store/user"
import { Delete, Edit } from "@vicons/carbon"

const { project } = storeToRefs(useProjectStore())
const userStore = useUserStore()

const expandColumn = reactive({
  key: "expand",
  type: "expand",
  expandable: () => userStore.hasResourceCode("project:detail:plan:instance"),
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
      [
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
  render: (row: ProjectIssue) => dayjs(row.startTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const endTimeColumn = reactive({
  title: "解决完成时间",
  key: "endTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => dayjs(row.endTime).format("YYYY-MM-DD"),
  // 排序
  sorter: true,
  sortOrder: false
})
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectIssue) => {
    switch (row.status) {
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
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const operationColumn = reactive({
  title: "操作",
  key: "operation",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: ProjectIssue) => {
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
                disabled: !userStore.hasResourceCode("project:detail:issue:update"),
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
                      disabled: !userStore.hasResourceCode("project:detail:issue:delete")
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
    title: "标题",
    key: "title"
  },
  startTimeColumn,
  endTimeColumn,
  statusColumn,
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
  ]
}
const handleEditData = (row: ProjectIssue) => {
  instance.value = row
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectIssue(id).then((res) => {
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

// 加载完毕
onMounted(() => {
  refresh()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-space justify="space-between">
        <span></span>
        <n-button type="primary" @click="newInstance" v-resource-code="'project:detail:issue:add'">新建缺陷</n-button>
      </n-space>
    </n-gi>
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

  <n-drawer v-model:show="drawerActive" :default-height="600" resizable placement="bottom">
    <n-drawer-content>
      <template #header>
        <n-text>{{ drawerTitle }}</n-text>
        <n-button
          class="absolute right-8px mt--4px"
          type="primary"
          size="small"
          @click="submit"
          v-resource-code="['project:detail:issue:add', 'project:detail:issue:update']"
          >提交</n-button
        >
      </template>
      <n-form ref="formRef" :model="instance" :rules="planRules" label-width="120px" label-placement="left">
        <n-grid :cols="4" x-gap="4">
          <n-gi>
            <n-form-item label="缺陷名称：" path="title">
              <n-input v-model:value="instance.title" placeholder="请输入缺陷名称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="类型：">
              <n-select
                v-model:value="instance.type"
                placeholder="请选择类型"
                :options="[
                  { label: '代码错误', value: 1 },
                  { label: '功能异常', value: 2 },
                  { label: '界面优化', value: 3 },
                  { label: '配置相关', value: 4 },
                  { label: '安全相关', value: 5 },
                  { label: '性能相关', value: 6 },
                  { label: '其他问题', value: 9 }
                ]"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="缺陷描述：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.content" placeholder="请输入缺陷描述" />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="原因：">
              <n-input
                type="textarea"
                :autosize="{ minRows: 2, maxRows: 5 }"
                v-model:value="instance.issueCause"
                placeholder="请输入问题原因，即为什么会出现这个缺陷"
              />
            </n-form-item>
          </n-gi>
          <n-gi :span="4">
            <n-form-item label="解决方法：">
              <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.solveMethod" placeholder="请输入解决方法" />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
