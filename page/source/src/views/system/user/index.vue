<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { getUserList, updateUser, addUser, deleteUser, getUserRoleIdList } from "@/service/api/settings/user"
import { User, UserCondition } from "@/types/user"
import { Delete, Edit, Search, UserRole } from "@vicons/carbon"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, DataTableFilterState, DataTableSortState, NIcon, NTooltip } from "naive-ui"
import RoleDrawer from "./roleDrawer/index.vue"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<UserCondition>({
  username: "",
  realName: "",
  email: "",
  mobile: "",
  status: 0,
  offset: 0,
  limit: 10
})
const list = ref<User[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getUserList(condition.value)
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
  render: (row: User) => {
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
const createTimeColumn = reactive({
  title: "创建时间",
  key: "createTime",
  // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
  render: (row: User) => dayjs(row.createTime).fromNow(),
  // 排序
  sorter: true,
  sortOrder: false
})
const columns = [
  {
    title: "用户名",
    key: "username"
  },
  {
    title: "姓名",
    key: "realName"
  },
  {
    title: "邮箱",
    key: "email"
  },
  {
    title: "手机号",
    key: "mobile",
    // 中间四位隐藏
    render: (row: User) => row.mobile?.replace(/(\d{3})\d{4}(\d{4})/, "$1****$2")
  },
  statusColumn,
  createTimeColumn,
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: User) => {
      return h(NButtonGroup, {}, () => [
        h(
          NTooltip,
          {},
          {
            default: () => "分配角色",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  disabled: !userStore.hasResourceCode("system:user:dispatchRoles"),
                  type: "warning",
                  onClick: () => handleAssignRole(row)
                },
                {
                  default: () =>
                    h(
                      NIcon,
                      {},
                      {
                        default: () => h(UserRole)
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
                  disabled: !userStore.hasResourceCode("system:user:edit"),
                  type: "primary",
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
                        disabled: row.username === "admin" || !userStore.hasResourceCode("system:user:delete"),
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
    if (columnKey === "createTime") {
      createTimeColumn.sortOrder = order === "ascend"
      condition.value.orderBy = "create_time " + (order === "ascend" ? "asc" : "desc")
      refresh()
    }
  }
}
const handleEditData = (row: User) => {
  checkedData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteUser(id).then((res) => {
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
  checkedData.value = { id: "", createTime: 0, email: "", mobile: "", realName: "", status: 1, username: "", password: "" }
}
const handleAddUser = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改用户" : "新增用户"))
const checkedData = ref<User>({ id: "", createTime: 0, email: "", mobile: "", realName: "", status: 1, username: "", password: "" })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateUser(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addUser(checkedData.value).then((res) => {
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
  username: {
    required: true,
    message: "请输入用户名",
    min: 3,
    max: 20,
    trigger: "blur"
  },
  realName: {
    required: true,
    message: "请输入姓名",
    min: 2,
    trigger: "blur"
  },
  email: [
    {
      required: true,
      message: "请输入邮箱",
      trigger: "blur"
    },
    {
      pattern: /^\w+@\w+(\.\w+)+$/,
      message: "请输入正确的邮箱",
      trigger: "blur"
    }
  ],
  mobile: [
    {
      required: true,
      message: "请输入手机号",
      trigger: "blur"
    },
    {
      pattern: /^1[3-9]\d{9}$/,
      message: "请输入正确的手机号",
      trigger: "blur"
    }
  ],
  password: [
    {
      required: true,
      message: "请输入密码",
      trigger: "blur"
    },
    {
      pattern: /^.*(?=.{5,20})(?=.*\d)(?=.*[A-Z]+)(?=.*[a-z]+)(?=.*[!@#$%^&*?]).*$/,
      message: "密码长度为5-20位，必须至少包含大小写字母、数字和特殊字符",
      trigger: "blur"
    }
  ]
}

// 分配角色抽屉
const roleDrawerActive = ref(false)
const roleUserId = ref<string>("")
const userRoleIdList = ref<string[]>([])
const handleAssignRole = (row: User) => {
  roleUserId.value = row.id
  getUserRoleIdList(row.id).then((res) => {
    userRoleIdList.value = res
    roleDrawerActive.value = true
  })
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>用户管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddUser" v-resource-code="'system:user:add'">新增用户</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="用户名" v-model:value="condition.username" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi :span="2">
          <n-input size="small" placeholder="姓名" v-model:value="condition.realName" @keydown.enter.prevent="refresh" />
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

        <n-gi :offset="2" class="flex flex-justify-end">
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
        :row-key="(row: User) => row.id"
        :pagination="paginationReactive"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :on-update:filters="handleFiltersChange"
        :on-update:sorter="handleSorterChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="用户名" path="username">
          <n-input :disabled="isUpdate" v-model:value="checkedData.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item label="姓名" path="realName">
          <n-input v-model:value="checkedData.realName" placeholder="请输入姓名" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="checkedData.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item label="手机" path="mobile">
          <n-input v-model:value="checkedData.mobile" placeholder="请输入手机" />
        </n-form-item>
        <n-form-item v-if="isUpdate" label="状态" required>
          <n-radio-group v-model:value="checkedData.status">
            <n-radio :value="1">正常</n-radio>
            <n-radio :value="2">禁用</n-radio>
          </n-radio-group>
        </n-form-item>
        <n-form-item v-if="!isUpdate" label="密码" path="password">
          <n-input type="password" show-password-on="mousedown" :minlength="5" v-model:value="checkedData.password" placeholder="请输入密码" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:user:add', 'system:user:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>

  <role-drawer v-model:drawer-active="roleDrawerActive" :user-id="roleUserId" :user-role-id-list="userRoleIdList" />
</template>

<style scoped></style>
