<template>
  <n-grid :cols="8" x-gap="16">
    <n-gi :span="1">
      <MyProjectList :selected-project-id="selectedProjectId" @update:selected-project-id="handleSelectedProjectIdUpdate" />
    </n-gi>
    <n-gi :span="7">
      <n-grid :cols="4" x-gap="16" y-gap="8">
        <n-gi :span="4">
          <n-space justify="space-between">
            <n-button type="primary" size="small" @click="openAddDrawer()">新增</n-button>
            <n-button type="info" v-if="userStore.hasResourceCode('workTime:statistics')" size="small" @click="showStatistics()">查看统计</n-button>
          </n-space>
        </n-gi>
        <n-gi v-for="workTime in workTimeList" :key="workTime.id">
          <n-card size="small" :title="dayjs(workTime.startDate).format('YYYY-MM-DD') + ' → ' + dayjs(workTime.endDate).format('YYYY-MM-DD')">
            <div>工作时长： {{ workTime.workTime.toFixed(2) }} 小时</div>
            <template #action>
              <div class="flex justify-end">
                <n-button-group>
                  <n-button type="primary" size="tiny" @click="openAddDrawer(workTime)">编辑</n-button>
                  <n-popconfirm @positive-click="handleDeleteWorkTime(workTime.id)">
                    <template #trigger>
                      <n-button type="error" size="tiny">删除</n-button>
                    </template>
                    确定要删除该工时记录吗？
                  </n-popconfirm>
                </n-button-group>
              </div>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="addDrawer" :width="401" placement="right">
    <n-drawer-content :title="instance.id === '' ? '新增工时' : '修改工时'">
      <n-form label-placement="top" :model="instance">
        <n-form-item label="记录周期" required>
          <n-date-picker v-model:value="instanceDateRange" type="daterange" :is-date-disabled="isRangeDateDisabled" />
        </n-form-item>
        <n-form-item label="工作时长" required>
          <n-input-number v-model:value="instance.workTime" />
        </n-form-item>
        <n-form-item label="工作内容简述">
          <n-input type="textarea" v-model:value="instance.workContent" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="space-between">
          <n-button type="error" @click="addDrawer = false">取消</n-button>
          <n-button type="primary" @click="submitWorkTime()">确定</n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import MyProjectList from "@/components/project/MyProjectList.vue"
import { WorkTime, WorkTimeCondition } from "@/types/log"
import { addWorkTime, deleteWorkTime, getWorkTimeList, updateWorkTime } from "@/service"
import dayjs from "dayjs"
import { useUserStore } from "@/store/user"

const message = useMessage()
const userStore = useUserStore()

const selectedProjectId = ref("")
const handleSelectedProjectIdUpdate = (id: string | undefined) => {
  selectedProjectId.value = id || ""
  condition.value.projectId = selectedProjectId.value
  refreshWorkTimeList()
}

const condition = ref<WorkTimeCondition>({
  offset: -1,
  limit: -1,
  projectId: "",
  startDateCondition: {}
})

const workTimeList = ref<WorkTime[]>([])

const refreshWorkTimeList = () => {
  if (condition.value.projectId === "") {
    return
  }
  workTimeList.value = []
  getWorkTimeList(condition.value).then((res) => {
    if (res) {
      workTimeList.value = res.items.map((item) => {
        item.workTime = item.workTime / 3600
        return item
      })
    }
  })
}

// 新增
const instance = ref<WorkTime>({
  id: "",
  projectId: "",
  startDate: 0,
  endDate: 0,
  workTime: 0
})
const instanceDateRange = computed({
  get: () => {
    return [instance.value.startDate, instance.value.endDate]
  },
  set: (val) => {
    instance.value.startDate = val[0]
    instance.value.endDate = val[1]
  }
})
const isRangeDateDisabled = (ts: number, type: "start" | "end", range: [number, number] | null) => {
  if (type === "start") {
    return ts > Date.now()
  }
  if (type === "end" && range) {
    const result = ts > Date.now()
    return result
  }
  return false
}
const addDrawer = ref<boolean>(false)
const openAddDrawer = (data: WorkTime | undefined = undefined) => {
  if (data) {
    instance.value = Object.assign({}, data)
  } else {
    instance.value = {
      id: "",
      projectId: selectedProjectId.value,
      startDate: Date.now() - 86400000,
      endDate: Date.now(),
      workTime: 0
    }
  }
  if (instance.value.projectId === "") {
    message.error("请选择项目")
    return
  }
  addDrawer.value = true
}

const submitWorkTime = () => {
  const submitData = Object.assign({}, instance.value)

  if (submitData.startDate > Date.now()) {
    message.error("开始时间不能大于当前时间")
    return
  }
  if (submitData.workTime <= 0) {
    message.error("工作时长必须大于0")
    return
  }

  submitData.workTime = submitData.workTime * 3600
  if (submitData.endDate % 1000 === 0) {
    submitData.endDate = submitData.endDate + 86400000 - 1
  }
  if (submitData.id === "") {
    // 新增
    addWorkTime(submitData).then((res) => {
      if (res) {
        message.success("新增成功")
        addDrawer.value = false
        refreshWorkTimeList()
      }
    })
  } else {
    // 修改
    updateWorkTime(submitData).then((res) => {
      if (res) {
        message.success("修改成功")
        addDrawer.value = false
        refreshWorkTimeList()
      }
    })
  }
}

// 删除工时记录
const handleDeleteWorkTime = (id: string) => {
  deleteWorkTime(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refreshWorkTimeList()
    }
  })
}

const router = useRouter()
// 查看工时统计
const showStatistics = () => {
  router.push({
    name: "WorkTimeStatistics"
  })
}
</script>
