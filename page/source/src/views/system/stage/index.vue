<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { getStageList, updateStage, addStage, deleteStage } from "@/service/api"
import { Stage, StageCondition } from "@/types/project"
import { Search } from "@vicons/carbon"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, DataTableFilterState } from "naive-ui"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<StageCondition>({
  name: "",
  status: 0,
  offset: 0,
  limit: 10
})
const list = ref<Stage[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getStageList(condition.value)
    .then((res) => {
      list.value = res.items || []
      paginationReactive.itemCount = res.total
      paginationReactive.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
      paginationReactive.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
      statusColumn.filterOptionValues = [condition.value.status || 0]
    })
    .finally(() => {
      loading.value = false
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
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: Stage) => {
    switch (row.status) {
      case 1:
        return "正常"
      case 2:
        return "禁用"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "正常",
      value: 1
    },
    {
      label: "禁用",
      value: 2
    }
  ]
})
const columns = [
  {
    title: "阶段名称",
    key: "name"
  },
  {
    title: "排序号",
    key: "orderNum"
  },
  statusColumn,
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: Stage) => dayjs(row.createTime).fromNow()
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: Stage) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            disabled: !userStore.hasResourceCode("system:stage:update"),
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
                  disabled: !userStore.hasResourceCode("system:stage:delete"),
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
const handleFiltersChange = (filters: DataTableFilterState) => {
  if (!loading.value) {
    const filterValues = filters.status || []
    if (filterValues instanceof Array) {
      condition.value.status = (filterValues[0] as number) || 0
    }
    refresh()
  }
}
const handleEditData = (row: Stage) => {
  checkedData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteStage(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
onMounted(() => {
  refresh()
})
const resetCheckedData = () => {
  checkedData.value = { id: "", name: "", status: 1, orderNum: 0 }
}
const handleAddStage = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改用户" : "新增用户"))
const checkedData = ref<Stage>({ id: "", status: 1, name: "", orderNum: 0 })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateStage(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addStage(checkedData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          refresh()
          drawerActive.value = false
        }
      })
    }
  })
}
// 规则定义
const rules = {
  name: {
    required: true,
    message: "请输入阶段名称",
    min: 2,
    max: 10,
    trigger: "blur"
  },
  orderNum: [
    {
      validator(_: unknown, value: number) {
        return value > 0
      },
      type: "number",
      required: true,
      message: "请输入大于0的排序号",
      trigger: "blur"
    }
  ]
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>阶段管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddStage" v-resource-code="'system:stage:add'">新增阶段</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="阶段名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi>
          <!-- 状态条件 -->
          <n-select
            size="small"
            placeholder="状态"
            v-model:value="condition.status"
            @update:value="refresh"
            :options="[
              {
                label: '全部',
                value: 0
              },
              {
                label: '正常',
                value: 1
              },
              {
                label: '禁用',
                value: 2
              }
            ]"
          >
          </n-select>
        </n-gi>

        <n-gi :offset="4" class="flex flex-justify-end">
          <n-button size="small" @click="refresh">
            <template #icon>
              <n-icon><search /></n-icon>
            </template>
          </n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: Stage) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="阶段名称" path="name">
          <n-input v-model:value="checkedData.name" placeholder="请输入阶段名称" />
        </n-form-item>
        <n-form-item label="排序号" path="orderNum">
          <n-input-number v-model:value="checkedData.orderNum" placeholder="请输入排序号" />
        </n-form-item>
        <n-form-item v-if="isUpdate" label="状态" required>
          <n-radio-group v-model:value="checkedData.status">
            <n-radio :value="1">正常</n-radio>
            <n-radio :value="2">禁用</n-radio>
          </n-radio-group>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:stage:add', 'system:stage:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
