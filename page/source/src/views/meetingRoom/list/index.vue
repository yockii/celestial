<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { Search } from "@vicons/carbon"
import { MeetingRoom, MeetingRoomCondition } from "@/types/meetingRoom"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, NSwitch } from "naive-ui"
import { addMeetingRoom, deleteMeetingRoom, getMeetingRoomDetail, getMeetingRoomList, updateMeetingRoom } from "@/service/api"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<MeetingRoomCondition>({
  offset: 0,
  limit: 10
})
const list = ref<MeetingRoom[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getMeetingRoomList(condition.value)
    .then((res) => {
      list.value = res.items || []
      pagination.itemCount = res.total
      pagination.pageCount = Math.ceil(res.total / res.limit)
      pagination.page = Math.ceil(res.offset / res.limit) + 1
    })
    .finally(() => {
      loading.value = false
    })
}

const pagination = reactive({
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
const statusColumn = {
  title: "状态",
  key: "status",
  render: (row: MeetingRoom) => {
    // return row.status === 1 ? "启用" : "禁用"
    return h(
      NSwitch,
      {
        value: row.status === 1,
        onUpdateValue: (value: boolean) => {
          updateMeetingRoom({ id: row.id, status: value ? 1 : -1 }).then((res) => {
            if (res) {
              message.success("修改成功")
              refresh()
            }
          })
        }
      },
      {
        checked: () => "启用",
        unchecked: () => "禁用"
      }
    )
  }
}
const columns = [
  {
    title: "会议室",
    key: "name"
  },
  {
    title: "位置",
    key: "position"
  },
  {
    title: "最大人数",
    key: "capacity"
  },
  statusColumn,
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: MeetingRoom) => dayjs(row.createTime).fromNow()
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: MeetingRoom) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            disabled: !userStore.hasResourceCode("meetingRoom:update"),
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
                  disabled: !userStore.hasResourceCode("meetingRoom:delete"),
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
const handleEditData = (row: MeetingRoom) => {
  getMeetingRoomDetail(row.id).then((res) => {
    if (res) {
      checkedData.value = res
      drawerActive.value = true
    }
  })
}
const handleDeleteData = (id: string) => {
  deleteMeetingRoom(id).then((res) => {
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
  checkedData.value = { id: "", name: "", status: 0 }
}
const handleAddMeetingRoom = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改会议室" : "新增会议室"))
const checkedData = ref<MeetingRoom>({ id: "", name: "", status: 0 })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateMeetingRoom(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addMeetingRoom(checkedData.value).then((res) => {
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
    message: "请输入名称",
    min: 2,
    max: 20,
    trigger: "blur"
  },
  position: [
    {
      required: true,
      message: "请输入位置",
      trigger: "blur"
    }
  ],
  capacity: [
    {
      type: "number",
      required: true,
      message: "请输入最大人数",
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
          <n-h3>会议室管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddMeetingRoom" v-resource-code="'meetingRoom:add'">新增</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi>
          <n-input size="small" placeholder="位置" v-model:value="condition.position" @keydown.enter.prevent="refresh" />
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
        :row-key="(row: MeetingRoom) => row.id"
        :pagination="pagination"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form size="small" ref="formRef" :model="checkedData" :rules="rules" label-width="90px">
        <n-form-item label="名称" path="name" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="位置" path="type" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.position" placeholder="请输入位置" />
        </n-form-item>
        <n-form-item label="最大人数" path="capacity" label-placement="left" label-align="left">
          <n-input-number v-model:value="checkedData.capacity" placeholder="请输入最大人数" :min="2" />
        </n-form-item>
        <n-form-item label="设备列表" path="devices" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.devices" placeholder="请输入设备信息" />
        </n-form-item>
        <n-form-item label="备注" path="devices" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.remark" placeholder="请输入备注" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['meetingRoom:add', 'meetingRoom:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
