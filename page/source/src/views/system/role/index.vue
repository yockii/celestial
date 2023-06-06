<script setup lang="ts">
import { getRoleList, updateRole, addRole, deleteRole, setDefaultRole, getRoleResourceCodeList } from "@/service/api/role"
import { RoleCondition, Role } from "@/types/user"
import { DataSet, Delete, Edit, GroupResource, Search } from "@vicons/carbon"
import dayjs from "dayjs"
import {
  NTag,
  NButtonGroup,
  NButton,
  NPopconfirm,
  FormInst,
  useMessage,
  SelectOption,
  SelectGroupOption,
  DataTableFilterState,
  DataTableSortState,
  PaginationProps,
  NTooltip,
  NIcon
} from "naive-ui"
import ResourceDrawer from "./resourceDrawer/index.vue"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<RoleCondition>({
  name: "",
  dataPermission: 0,
  type: 0,
  status: 0,
  offset: 0,
  limit: 10
})
const list = ref<Role[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getRoleList(condition.value)
    .then((res) => {
      list.value = res.items
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
const nameColumn = {
  title: "角色名称",
  key: "name",
  render: (row: Role) => {
    return h(
      NTag,
      {
        color: row.style ? JSON.parse(row.style) : undefined
      },
      {
        default: () => row.name
      }
    )
  }
}
const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: Role) => {
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
const typeColumn = reactive({
  title: "类型",
  key: "type",
  render: (row: Role) => {
    switch (row.type) {
      case 1:
        return "普通"
      case 2:
        return "项目"
      case -1:
        return "超级管理员"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "普通",
      value: 1
    },
    {
      label: "项目",
      value: 2
    },
    {
      label: "超级管理员",
      value: -1
    }
  ]
})
const dataPermColumn = reactive({
  title: "数据权限",
  key: "dataPermission",
  render: (row: Role) => {
    switch (row.dataPermission) {
      case 1:
        return "所有数据"
      case 2:
        return "本级及子级数据"
      case 3:
        return "仅自己的数据"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValues: [0],
  filterOptions: [
    {
      label: "所有数据",
      value: 1
    },
    {
      label: "本级及子级数据",
      value: 2
    },
    {
      label: "仅自己的数据",
      value: 3
    }
  ]
})
// 排序字段
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: Role) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const columns = [
  nameColumn,
  typeColumn,
  dataPermColumn,
  statusColumn,
  createTimeColumn,
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: Role) => {
      return h(NButtonGroup, {}, () => [
        h(
          NTooltip,
          {},
          {
            default: () => "设为默认角色",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  type: "primary",
                  disabled: row.defaultRole === 1 || !userStore.hasResourceCode("system:role:update"),
                  onClick: () => handleSetDefault(row.id)
                },
                {
                  icon: () =>
                    h(
                      NIcon,
                      {},
                      {
                        default: () => h(DataSet)
                      }
                    )
                }
              )
          }
        ),
        h(
          NTooltip,
          {},
          {
            default: () => "分配资源",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  type: "warning",
                  disabled: !userStore.hasResourceCode("system:role:dispatchResources"),
                  onClick: () => handleAssignResource(row)
                },
                {
                  default: () =>
                    h(
                      NIcon,
                      {},
                      {
                        default: () => h(GroupResource)
                      }
                    )
                }
              )
          }
        ),
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
                  disabled: !userStore.hasResourceCode("system:role:update"),
                  onClick: () => handleEditData(row)
                },
                {
                  default: () =>
                    h(
                      NIcon,
                      {},
                      {
                        default: () => h(Edit)
                      }
                    )
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
                        disabled: row.defaultRole === 1 || !userStore.hasResourceCode("system:role:delete"),
                        type: "error"
                      },
                      {
                        default: () =>
                          h(
                            NIcon,
                            {},
                            {
                              default: () => h(Delete)
                            }
                          )
                      }
                    )
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
    const statusValues = filters.status || []
    if (statusValues instanceof Array) {
      condition.value.status = (statusValues[0] as number) || 0
    }
    const typeValues = filters.type || []
    if (typeValues instanceof Array) {
      condition.value.type = (typeValues[0] as number) || 0
    }
    const dpValues = filters.dataPermission || []
    if (dpValues instanceof Array) {
      condition.value.dataPermission = (dpValues[0] as number) || 0
    }
    refresh()
  }
}

const handleSorterChange = (sorter: DataTableSortState) => {
  if (!loading.value) {
    const { columnKey, order } = sorter
    if (columnKey === "createTime") {
      createTimeColumn.sortOrder = order === "ascend"
      condition.value.orderBy = "create_time " + (order === "ascend" ? "asc" : "desc")
      refresh()
    }
  }
}
const handleSetDefault = (id: string) => {
  setDefaultRole(id).then((res) => {
    if (res) {
      message.success("设置成功")
    }
    refresh()
  })
}
const handleEditData = (row: Role) => {
  checkedData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteRole(id).then((res) => {
    if (res) {
      message.success("删除成功")
    }
    refresh()
  })
}
const resetRoleData = () => {
  checkedData.value = {
    dataPermission: 3,
    desc: "",
    type: 1,
    id: "",
    name: "",
    status: 1
  }
}
const handleAddRole = () => {
  resetRoleData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改角色" : "新增角色"))
const checkedData = ref<Role>({ dataPermission: 0, desc: "", type: 0, id: "", status: 1, name: "" })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  console.log(checkedData.value)
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateRole(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addRole(checkedData.value).then((res) => {
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
    message: "请输入角色名称",
    min: 2,
    max: 10,
    trigger: "blur"
  }
}
// 角色样式调色
const presetBgColors = ref<string[]>([
  "",
  "rgba(245, 63, 63, .16)",
  "rgba(246, 114, 52, .16)",
  "rgba(247, 186, 31, .16)",
  "rgba(159, 219, 29, .16)",
  "rgba(53, 180, 42, .16)",
  "rgba(65, 201, 201, .16)",
  "rgba(53, 145, 250, .16)",
  "rgba(153, 80, 255, .16)",
  "rgba(217, 27, 217, .16)",
  "rgba(245, 49, 157, .16)",
  "rgba(134, 144, 156, .16)"
])
const presetTextColors = ref<string[]>([
  "",
  "rgba(245, 63, 63, 1)",
  "rgba(246, 114, 52, 1)",
  "rgba(247, 186, 31, 1)",
  "rgba(159, 219, 29, 1)",
  "rgba(53, 180, 42, 1)",
  "rgba(65, 201, 201, 1)",
  "rgba(53, 145, 250, 1)",
  "rgba(153, 80, 255, 1)",
  "rgba(217, 27, 217, 1)",
  "rgba(245, 49, 157, 1)",
  "rgba(134, 144, 156, 1)"
])
const presetBorderColors = ref<string[]>([
  "",
  "rgba(245, 63, 63, .6)",
  "rgba(246, 114, 52, .6)",
  "rgba(247, 186, 31, .6)",
  "rgba(159, 219, 29, .6)",
  "rgba(53, 180, 42, .6)",
  "rgba(65, 201, 201, .6)",
  "rgba(53, 145, 250, .6)",
  "rgba(153, 80, 255, .6)",
  "rgba(217, 27, 217, .6)",
  "rgba(245, 49, 157, .6)",
  "rgba(134, 144, 156, .6)"
])
const bgColor = computed({
  get: () => {
    if (checkedData.value?.style) {
      return JSON.parse(checkedData.value.style).color
    }
    return ""
  },
  set: (val) => {
    if (checkedData.value?.style) {
      const sJson = JSON.parse(checkedData.value.style)
      sJson.color = val
      checkedData.value.style = JSON.stringify(sJson)
    } else {
      checkedData.value.style = JSON.stringify({ color: val })
    }
  }
})
const textColor = computed({
  get: () => {
    if (checkedData.value?.style) {
      return JSON.parse(checkedData.value.style).textColor
    }
    return ""
  },
  set: (val) => {
    if (checkedData.value?.style) {
      const sJson = JSON.parse(checkedData.value.style)
      sJson.textColor = val
      checkedData.value.style = JSON.stringify(sJson)
    } else {
      checkedData.value.style = JSON.stringify({ textColor: val })
    }
  }
})
const borderColor = computed({
  get: () => {
    if (checkedData.value?.style) {
      return JSON.parse(checkedData.value.style).borderColor
    }
    return ""
  },
  set: (val) => {
    if (checkedData.value?.style) {
      const sJson = JSON.parse(checkedData.value.style)
      sJson.borderColor = val
      checkedData.value.style = JSON.stringify(sJson)
    } else {
      checkedData.value.style = JSON.stringify({ borderColor: val })
    }
  }
})
const selectedThemeValue = computed(() => {
  // 返回bgColor的值在presetBgColors中的索引
  return presetBgColors.value.findIndex((item) => item === bgColor.value)
})
const themeOptions = ref<Array<SelectOption | SelectGroupOption>>([])

const themeRenderer = (option: SelectOption | SelectGroupOption) => {
  return [
    h(
      NTag,
      {
        color: {
          color: presetBgColors.value[option.value as number],
          textColor: presetTextColors.value[option.value as number],
          borderColor: presetBorderColors.value[option.value as number]
        }
      },

      { default: () => "示例" }
    )
  ]
}
const changeTheme = (value: number) => {
  bgColor.value = presetBgColors.value[value]
  textColor.value = presetTextColors.value[value]
  borderColor.value = presetBorderColors.value[value]
}

onMounted(() => {
  // 给themeOptions赋值, 将presetBgColors转换为SelectOption数组
  themeOptions.value = presetBgColors.value.map((item, index) => {
    return {
      label: item,
      value: index
    }
  })

  refresh()
})

// 资源分配抽屉
const resourceDrawerActive = ref(false)
const resourceRoleId = ref("")
const roleResourceCodeList = ref<string[]>([])
const handleAssignResource = (row: Role) => {
  resourceRoleId.value = row.id
  getRoleResourceCodeList(row.id).then((res) => {
    roleResourceCodeList.value = res
    resourceDrawerActive.value = true
  })
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>角色管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddRole" v-resource-code="'system:role:add'">新增角色</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="角色名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
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
                label: '普通',
                value: 1
              },
              {
                label: '项目',
                value: 2
              },
              {
                label: '超级管理员',
                value: -1
              }
            ]"
          >
          </n-select>
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

        <n-gi :offset="3" class="flex flex-justify-end">
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
        :row-key="(row: Role) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :columns="columns"
        :on-update:sorter="handleSorterChange"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="角色名称" path="name">
          <n-input v-model:value="checkedData.name" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input type="textarea" v-model:value="checkedData.desc" placeholder="请输入描述信息" />
        </n-form-item>
        <n-form-item label="类型" required>
          <n-select
            v-model:value="checkedData.type"
            :options="[
              {
                label: '普通',
                value: 1
              },
              {
                label: '项目',
                value: 2
              },
              {
                label: '超级管理员',
                value: -1
              }
            ]"
          />
        </n-form-item>
        <n-form-item label="数据权限" required>
          <n-radio-group v-model:value="checkedData.dataPermission">
            <n-radio :value="1">所有数据</n-radio>
            <n-radio :value="2">本部门及以下</n-radio>
            <n-radio :value="3">仅本人</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="isUpdate" label="状态" required>
          <n-radio-group v-model:value="checkedData.status">
            <n-radio :value="1">正常</n-radio>
            <n-radio :value="2">禁用</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="主题样式">
          <n-select :value="selectedThemeValue" :options="themeOptions" :render-label="themeRenderer" @update:value="changeTheme" />
        </n-form-item>
        <n-form-item label="背景色">
          <n-color-picker v-model:value="bgColor" :swatches="presetBgColors" />
        </n-form-item>
        <n-form-item label="文字色">
          <n-color-picker v-model:value="textColor" :swatches="presetTextColors" />
        </n-form-item>
        <n-form-item label="边框色">
          <n-color-picker v-model:value="borderColor" :swatches="presetBorderColors" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetRoleData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:role:add', 'system:role:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>

  <resource-drawer v-model:drawer-active="resourceDrawerActive" :role-id="resourceRoleId" :role-resource-code-list="roleResourceCodeList" />
</template>

<style scoped></style>
