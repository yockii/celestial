<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { getAssetCategoryList, updateAssetCategory, addAssetCategory, deleteAssetCategory } from "@/service/api"
import { AssetCategory, AssetCategoryCondition } from "@/types/asset"
import { Search } from "@vicons/carbon"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, DataTableFilterState, DataTableBaseColumn } from "naive-ui"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<AssetCategoryCondition>({
  name: "",
  onlyParent: true,
  offset: 0,
  limit: 10
})
const list = ref<AssetCategory[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getAssetCategoryList(condition.value)
    .then((res) => {
      list.value = res.items.map((item) => {
        item.isLeaf = !item.childrenCount || item.childrenCount === 0
        return item
      })
      paginationReactive.itemCount = res.total
      paginationReactive.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
      paginationReactive.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
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
const typeColumn = {
  title: "类型",
  key: "type",
  filter: true,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "公共资产",
      value: 1
    },
    {
      label: "项目资产",
      value: 2
    },
    {
      label: "个人资产",
      value: 3
    }
  ],
  filterMultiple: false,
  render: (row: AssetCategory) => {
    switch (row.type) {
      case 1:
        return "公共资产"
      case 2:
        return "项目资产"
      case 3:
        return "个人资产"
      default:
        return "未知"
    }
  }
}
const columns = [
  {
    title: "资产目录名称",
    key: "name"
  },
  typeColumn,
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: AssetCategory) => dayjs(row.createTime).fromNow()
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: AssetCategory) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            type: "primary",
            disabled: !userStore.hasResourceCode("system:assetCategory:update"),
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
                  disabled: !userStore.hasResourceCode("system:assetCategory:delete"),
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
const handleEditData = (row: AssetCategory) => {
  checkedData.value = Object.assign({}, row)
  drawerActive.value = true
  typeChanged(checkedData.value.type)
}
const handleDeleteData = (id: string) => {
  deleteAssetCategory(id).then((res) => {
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
  checkedData.value = { id: "", name: "", type: 1 }
}
const handleAddAssetCategory = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改用户" : "新增用户"))
const checkedData = ref<AssetCategory>({ id: "", name: "", type: 1 })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateAssetCategory(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addAssetCategory(checkedData.value).then((res) => {
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
    message: "请输入资产目录名称",
    min: 2,
    max: 10,
    trigger: "blur"
  },
  type: [
    {
      validator(_: unknown, value: number) {
        return value === 1 || value === 2 || value === 9
      },
      type: "number",
      required: true,
      message: "请选择正确资产目录类型",
      trigger: "blur"
    }
  ]
}

// 加载子级
const getChildren = (parentId: string, type = 0) => {
  return getAssetCategoryList({ type, parentId, offset: -1, limit: -1, onlyParent: !parentId }).then((res) => {
    return res.items || []
  })
}
const loadChildren = (row: AssetCategory) => {
  return new Promise<void>((resolve) => {
    getChildren(row.id, row.type).then((res) => {
      row.children = res
      resolve()
    })
  })
}

// 父级目录加载
const typeChanged = (type: number) => {
  checkedData.value.type = type
  parentCategoryList.value = []
  getChildren("", type).then((res) => {
    parentCategoryList.value = res.map((item) => {
      item.isLeaf = !item.childrenCount || item.childrenCount === 0
      return item
    })
  })
}
const parentCategoryList = ref<AssetCategory[]>([])
const loadParentCategory = (option: AssetCategory) => {
  if (!checkedData.value.type) {
    return
  }
  return new Promise<void>((resolve) => {
    getChildren(option.id || "", checkedData.value.type).then((res) => {
      option.children = res.map((item) => {
        item.isLeaf = !item.childrenCount || item.childrenCount === 0
        return item
      })
      resolve()
    })
  })
}
const parentCategoryChanged = (parentId: string) => {
  checkedData.value.parentId = parentId
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>资产目录管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddAssetCategory" v-resource-code="'system:assetCategory:add'">新增资产目录</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="资产目录名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi>
          <!-- 类型条件 -->
          <n-select
            size="small"
            placeholder="类型"
            v-model:value="condition.type"
            @update:value="refresh"
            :options="[
              {
                label: '全部',
                value: 0
              },
              {
                label: '公共资产',
                value: 1
              },
              {
                label: '项目资产',
                value: 2
              },
              {
                label: '个人资产',
                value: 9
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
        :row-key="(row: AssetCategory) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :columns="columns"
        @load="loadChildren"
        :cascade="false"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="目录名称" path="name">
          <n-input v-model:value="checkedData.name" placeholder="请输入资产目录名称" />
        </n-form-item>
        <n-form-item label="类型" required>
          <n-select
            v-model:value="checkedData.type"
            placeholder="请选择目录类型"
            :options="[
              {
                label: '公共资产',
                value: 1
              },
              {
                label: '项目资产',
                value: 2
              },
              {
                label: '个人资产',
                value: 9
              }
            ]"
            :on-update:value="typeChanged"
          />
        </n-form-item>
        <n-form-item label="父级目录" path="parentId">
          <n-cascader
            v-model:value="checkedData.parentId"
            placeholder="请选择父级目录"
            :options="parentCategoryList"
            label-field="name"
            value-field="id"
            clearable
            remote
            :on-load="loadParentCategory"
            :on-update:value="parentCategoryChanged"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:assetCategory:add', 'system:assetCategory:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
