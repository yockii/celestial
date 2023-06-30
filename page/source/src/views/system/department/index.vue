<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { getDepartmentList, updateDepartment, addDepartment, deleteDepartment } from "@/service/api"
import { Department, DepartmentCondition } from "@/types/user"
import { Search } from "@vicons/carbon"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps } from "naive-ui"
import { useUserStore } from "@/store/user"

const message = useMessage()
const userStore = useUserStore()
const condition = ref<DepartmentCondition>({
  name: "",
  onlyParent: true,
  offset: 0,
  limit: 10
})
const list = ref<Department[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getDepartmentList(condition.value)
    .then((res) => {
      list.value = res.items.map((item) => {
        item.isLeaf = !item.childCount || item.childCount === 0
        return item
      })
      if (condition.value.onlyParent) {
        parentDepartmentList.value = list.value
      }
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
const columns = [
  {
    title: "部门名称",
    key: "name"
  },
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: Department) => dayjs(row.createTime).fromNow()
  },
  {
    title: "序号",
    key: "orderNum"
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: Department) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            type: "primary",
            disabled: !userStore.hasResourceCode("system:department:update"),
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
                  disabled: !userStore.hasResourceCode("system:department:delete"),
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
const handleEditData = (row: Department) => {
  checkedData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteDepartment(id).then((res) => {
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
  checkedData.value = { id: "", name: "" }
}
const handleAddDepartment = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改部门" : "新增部门"))
const checkedData = ref<Department>({ id: "", name: "" })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateDepartment(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addDepartment(checkedData.value).then((res) => {
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
    message: "请输入部门名称",
    min: 2,
    max: 10,
    trigger: "blur"
  }
}

// 加载子级
const getChildren = (parentId: string) => {
  return getDepartmentList({ parentId, offset: -1, limit: -1, onlyParent: !parentId }).then((res) => {
    const list = res.items || []
    return list.map((item) => {
      item.isLeaf = !item.childCount || item.childCount === 0
      return item
    })
  })
}
const loadChildren = (row: Department) => {
  return new Promise<void>((resolve) => {
    getChildren(row.id).then((res) => {
      row.children = res
      resolve()
    })
  })
}

// 父级目录加载
const parentDepartmentList = ref<Department[]>([])
const loadParentDepartment = (option: Department) => {
  return new Promise<void>((resolve) => {
    getChildren(option.id || "").then((res) => {
      option.children = res.map((item) => {
        item.isLeaf = !item.childCount || item.childCount === 0
        return item
      })
      resolve()
    })
  })
}
const parentDepartmentChanged = (parentId: string) => {
  checkedData.value.parentId = parentId
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>部门管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddDepartment" v-resource-code="'system:department:add'">新增部门</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="部门名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
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
        :row-key="(row: Department) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :columns="columns"
        @load="loadChildren"
        :cascade="false"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="部门名称" path="name">
          <n-input v-model:value="checkedData.name" placeholder="请输入部门名称" />
        </n-form-item>
        <n-form-item label="父级部门" path="parentId">
          <n-cascader
            v-model:value="checkedData.parentId"
            placeholder="请选择父级部门"
            :options="parentDepartmentList"
            label-field="name"
            value-field="id"
            clearable
            remote
            :on-load="loadParentDepartment"
            :on-update:value="parentDepartmentChanged"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:department:add', 'system:department:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
