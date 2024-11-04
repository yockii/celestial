<template>
  <n-grid :cols="8" x-gap="16">
    <n-gi :span="1">
      <MyProjectList :selected-project-id="selectedProjectId"
        @update:selected-project-id="handleSelectedProjectIdUpdate" />
    </n-gi>
    <n-gi :span="7">
      <n-grid :cols="4" x-gap="16" y-gap="8">
        <n-gi :span="4">
          <n-space justify="space-between">
            <n-button type="primary" size="small" @click="openAddDrawer()">新增</n-button>
            <div class="flex gap-10">
              开始时间范围：<n-date-picker v-model:value="conditionStartDateCondition" type="daterange"
                 />
              <n-button type="primary" size="small" @click="getMyTotalWorktime()">我的总工时(该时间范围)</n-button>
            </div>
            <n-button type="info" v-if="userStore.hasResourceCode('workTime:statistics')" size="small"
              @click="showStatistics()">查看统计</n-button>
          </n-space>
        </n-gi>
        <n-gi v-for="workTime in workTimeList" :key="workTime.id">
          <n-card size="small"
            :title="dayjs(workTime.startDate).format('YYYY-MM-DD HH:mm:ss') + ' → ' + dayjs(workTime.endDate).format('YYYY-MM-DD HH:mm:ss')">
            <n-tooltip max-width="400">
              <template #trigger>
                <div class="cursor-pointer">工作时长： {{ workTime.workTime.toFixed(2) }} 小时</div>
              </template>
              <div>{{ workTime.workContent }}</div>
            </n-tooltip>
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
          <n-date-picker v-model:value="instanceDateRange" type="datetimerange" :is-date-disabled="isRangeDateDisabled"
            :disabled="instance.id !== ''" />
        </n-form-item>
        <n-form-item label="工作时长" required>
          <n-input-number v-model:value="instance.workTime" />
        </n-form-item>
        <n-form-item label="工作内容简述" required>
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
import { addWorkTime, deleteWorkTime, getWorkTimeList, updateWorkTime, getMyWorkTimeStatistics } from "@/service"
import dayjs from "dayjs"
import { useUserStore } from "@/store/user"

const message = useMessage()
const userStore = useUserStore()

// 我的总工时
const getMyTotalWorktime = () => {
  getMyWorkTimeStatistics(condition.value).then((res) => {
    if (res) {
      message.success(`我的总工时：${(res / 3600).toFixed(2)}小时`)
    } else if (res === 0) {
      // 工时为null，没有工时记录
      message.info("没有工时记录")
    }
  })
}

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
  startDateCondition: {
    // 本周一
    start: dayjs().startOf("week").format("YYYY-MM-DD HH:mm:ss"),
    // 当天
    end: dayjs().endOf("day").format("YYYY-MM-DD HH:mm:ss")
  }
})

const conditionStartDateCondition = computed({
  get: () => {
    return [dayjs(condition.value.startDateCondition?.start).valueOf(), dayjs(condition.value.startDateCondition?.end).valueOf()]
  },
  set: (val) => {
    condition.value.startDateCondition = {
      start: dayjs(val[0]).format("YYYY-MM-DD HH:mm:ss"),
      end: dayjs(val[1]).format("YYYY-MM-DD HH:mm:ss")
    }
    refreshWorkTimeList()
  }
})

const workTimeList = ref<WorkTime[]>([])

const refreshWorkTimeList = () => {
  if (condition.value.projectId === "") {
    return
  }
  workTimeList.value = []
  // 要修改一下condition中的 startDateCondition.start (-1)，因为查询时使用了between，不包含当天0点
  // 将condition.value 赋值给一个新的变量，否则会导致condition.value.startDateCondition.start被修改
  const conditionCopy = Object.assign({}, condition.value)
  if (conditionCopy.startDateCondition && conditionCopy.startDateCondition.start) {
    // 减去1毫秒再重新赋值为format("YYYY-MM-DD HH:mm:ss")
    conditionCopy.startDateCondition.start = dayjs(conditionCopy.startDateCondition.start).add(-1, "millisecond").format("YYYY-MM-DD HH:mm:ss")
  }
  getWorkTimeList(conditionCopy).then((res) => {
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
  // 周一-周日，只能填报当周的。
  // 当前时间
  const now = dayjs()
  let earliestPermitTime = 0
  // 如果当前是周一9点前
  if (now.hour() < 9 && now.day() === 1) {
    // 获取上周一0点
    earliestPermitTime = dayjs().subtract(1, "week").startOf("week").valueOf()
  } else { // 其他时间，只能填报本周一0点开始
    earliestPermitTime = dayjs().startOf("week").valueOf()
  }

  if (type === "start") {
    return ts > Date.now() || ts < earliestPermitTime
  }
  if (type === "end" && range) {
    const result = ts > Date.now() || ts <= range[0]
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
      // 当天0点
      startDate: dayjs().startOf("date").valueOf(),
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
  // 工作内容描述必填
  if (submitData.workContent === "") {
    message.error("工作内容描述必填")
    return
  }
  // 结束时间不能小于等于开始时间
  if (submitData.endDate <= submitData.startDate) {
    message.error("结束时间不能小于等于开始时间")
    return
  }

  submitData.workTime = submitData.workTime * 3600
  // if (submitData.endDate % 1000 === 0) {
  //   submitData.endDate = submitData.endDate + 86400000 - 1
  // }
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
