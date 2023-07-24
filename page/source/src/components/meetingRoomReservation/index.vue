<script setup lang="ts">
import { MeetingRoom, MeetingRoomReservation } from "@/types/meetingRoom"
import { deleteMeetingRoomReservation, getMeetingRoomReservationList, reserveMeetingRoom, updateMeetingRoomReservation } from "@/service/api"
import dayjs from "dayjs"
import { Edit } from "@vicons/carbon"
import { FormRules, FormInst } from "naive-ui"
import { useUserStore } from "@/store/user"

const props = defineProps<{
  room: MeetingRoom
}>()
const message = useMessage()
const userStore = useUserStore()
const selectedDate = ref(Date.now())
const reservationList = ref<MeetingRoomReservation[]>([])

const today = dayjs().startOf("date").valueOf()

const isDateDisabled = (timestamp: number) => {
  if (today > timestamp) {
    return true
  }
  return false
}
const canBook = (y: number, m: number, d: number) => {
  let judgeDate = dayjs().startOf("date")
  judgeDate = judgeDate
    .year(y)
    .month(m - 1)
    .date(d)

  if (today > judgeDate.valueOf()) {
    return false
  }

  const now = dayjs().valueOf()
  if (today === judgeDate.valueOf()) {
    if (now > dayjs().startOf("date").hour(18).valueOf()) {
      return false
    }
  }

  // 找出当天的所有预约记录的总时间
  const total = reservationList.value
    .filter((item) => {
      return judgeDate.valueOf() === dayjs(item.startTime).startOf("date").valueOf()
    })
    .reduce((prev, cur) => {
      return prev + (cur.endTime - cur.startTime)
    }, 0)

  // 一天最多预约7小时
  if (total >= 7 * 60 * 60 * 1000) {
    return false
  }
  return true
}
const dateReservationList = (y: number, m: number, d: number) => {
  const judgeDate = dayjs()
    .startOf("date")
    .year(y)
    .month(m - 1)
    .date(d)
  return reservationList.value.filter((item) => {
    return judgeDate.valueOf() === dayjs(item.startTime).startOf("date").valueOf()
  })
}

const bookStartLimit = ref<number>(0)
const bookEndLimit = ref<number>(0)
const bookDrawerActive = ref(false)
const book = (y: number, m: number, d: number) => {
  bookStartLimit.value = dayjs()
    .startOf("date")
    .year(y)
    .month(m - 1)
    .date(d)
    .hour(9)
    .valueOf()
  // const nearest = dayjs().add(5, "minute").valueOf()
  // if (nearest > bookStartLimit.value) {
  //   bookStartLimit.value = nearest
  // }
  bookEndLimit.value = dayjs()
    .startOf("date")
    .year(y)
    .month(m - 1)
    .date(d)
    .hour(18)
    .valueOf()

  // 计算一下当前已有的预约记录中，最早可用的时间
  const records = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }
    if (bookInfo.value.id === item.id) {
      return false
    }
    const itemStartDayjs = dayjs(item.startTime)
    return itemStartDayjs.hour() === 9 && itemStartDayjs.minute() === 0
  })
  if (records.length > 0) {
    // 如果有9点整的记录，那么最早可用时间需要后延
    // 获取当天所有预约记录，顺序排序后，检查结束时间是否为下一个开始时间，若不是，则该结束时间作为最早可用时间
    const sortedRecords = reservationList.value
      .filter((item) => {
        const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
        if (!isSameDay) {
          return false
        }

        if (bookInfo.value.id === item.id) {
          return false
        }
        return item.startTime >= bookStartLimit.value
      })
      .sort((a, b) => {
        return a.startTime - b.startTime
      })
    let calculated = false
    for (let i = 0; i < sortedRecords.length - 1; i++) {
      if (sortedRecords[i].endTime !== sortedRecords[i + 1].startTime) {
        bookStartLimit.value = sortedRecords[i].endTime
        calculated = true
        break
      }
    }
    if (!calculated) {
      bookStartLimit.value = sortedRecords[sortedRecords.length - 1].endTime
    }
  }

  // 根据最早可用时间，计算当前已有预约记录中，最晚可用的时间
  const records2 = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }

    if (bookInfo.value.id === item.id) {
      return false
    }
    return dayjs(item.startTime).hour() >= dayjs(bookStartLimit.value).hour()
  })
  if (records2.length > 0) {
    // 获取最大的结束时间
    const maxEndTime = records2.reduce((prev, cur) => {
      if (prev.startTime < cur.startTime) {
        return prev
      }
      return cur
    }).startTime
    bookEndLimit.value = dayjs(maxEndTime).valueOf()
  }

  bookInfo.value = {
    id: "",
    meetingRoomId: props.room.id,
    startTime: bookStartLimit.value,
    endTime: bookEndLimit.value,
    subject: ""
  }
  bookDrawerActive.value = true
}
const showEdit = (item: MeetingRoomReservation) => {
  bookStartLimit.value = item.startTime
  bookEndLimit.value = item.endTime
  bookInfo.value = {
    id: item.id,
    meetingRoomId: item.meetingRoomId,
    startTime: item.startTime,
    endTime: item.endTime,
    subject: item.subject,
    participants: item.participants
  }
  bookDrawerActive.value = true
}
const bookInfo = ref<MeetingRoomReservation>({
  id: "",
  meetingRoomId: props.room.id,
  startTime: 0,
  endTime: 0,
  subject: ""
})
const bookRules: FormRules = {
  startTime: [
    {
      type: "number",
      required: true,
      message: "请选择合法的开始时间"
    }
  ],
  endTime: [
    {
      type: "number",
      required: true,
      message: "请选择正确的结束时间"
    }
  ],
  subject: [
    {
      required: true,
      message: "请输入主题"
    }
  ]
}
const bookFormRef = ref<FormInst | null>(null)
const submitBook = () => {
  bookFormRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (bookInfo.value.id === "") {
      reserveMeetingRoom(bookInfo.value).then((res) => {
        if (res) {
          message.success("预约成功")
          bookDrawerActive.value = false
          refresh()
        }
      })
    } else {
      updateMeetingRoomReservation(bookInfo.value).then((res) => {
        if (res) {
          message.success("修改成功")
          bookDrawerActive.value = false
          refresh()
        }
      })
    }
  })
}
// 可选的开始时间列表
const findNextConnectedEndTime = (records: MeetingRoomReservation[], endTime: number): number => {
  const next = records.find((item) => {
    return item.startTime === endTime
  })
  if (next) {
    return findNextConnectedEndTime(records, next.endTime)
  }
  return endTime
}
const startHourList = computed(() => {
  const list = []
  // 找到所有当天的记录（非本次修改）
  const records = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }

    if (bookInfo.value.id === item.id) {
      return false
    }
    return true
  })
  if (records.length === 0) {
    for (let i = 9; i <= 18; i++) {
      list.push(i)
    }
  } else {
    // 将所有记录按照开始时间顺序排序
    records.sort((a, b) => {
      return a.startTime - b.startTime
    })
    // 循环查找可用的小时
    for (let i = 9; i <= 18; i++) {
      // 查找records中是否存在startTime.hour<i, endTime.hour>i的记录，存在则该小时不可用
      let found = records.find((item) => {
        const startHour = dayjs(item.startTime).hour()
        const startMinute = dayjs(item.startTime).minute()
        const endHour = dayjs(item.endTime).hour()
        const endMinute = dayjs(item.endTime).minute()
        return (startHour < i || (startHour === i && startMinute === 0)) && (endHour > i || (endHour === i && endMinute === 59))
      })
      if (found) {
        continue
      }
      // 另一种可能，多个记录拼起来，无空闲时间
      found = records.find((item) => {
        const startHour = dayjs(item.startTime).hour()
        const startMinute = dayjs(item.startTime).minute()
        if (startHour < i || (startHour === i && startMinute === 0)) {
          const nextEndTime = findNextConnectedEndTime(records, item.endTime)
          const nextEndHour = dayjs(nextEndTime).hour()
          const nextEndMinute = dayjs(nextEndTime).minute()
          return nextEndHour > i || (nextEndHour === i && nextEndMinute === 59)
        }
        return false
      })
      if (found) {
        continue
      }
      list.push(i)
    }
  }

  return list
})
// 可选的结束时间列表
const endHourList = computed(() => {
  const startHour = dayjs(bookInfo.value.startTime).hour()
  const list = []
  // 根据开始时间，查找在开始时间之后，最近的预约记录
  const records = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }

    if (bookInfo.value.id === item.id) {
      return false
    }
    return item.startTime > bookInfo.value.startTime
  })
  if (records.length === 0) {
    for (let i = startHour; i <= 18; i++) {
      list.push(i)
    }
  } else {
    // 获取最小的开始时间
    const minStartTime = records.reduce((prev, cur) => {
      if (prev.startTime < cur.startTime) {
        return prev
      }
      return cur
    }).startTime
    for (let i = startHour; i <= dayjs(minStartTime).hour(); i++) {
      list.push(i)
    }
  }
  return list
})
const startMinuteList = computed(() => {
  const list = []
  const selectedStartTimeHour = dayjs(bookInfo.value.startTime).hour()

  // 找到同一天的预约记录，并且开始/结束时间与当前所选的小时相同的记录
  const records = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }

    if (bookInfo.value.id === item.id) {
      return false
    }

    return dayjs(item.startTime).hour() === selectedStartTimeHour || dayjs(item.endTime).hour() === selectedStartTimeHour
  })
  if (records.length === 0) {
    for (let i = 0; i < 60; i++) {
      list.push(i)
    }
  } else {
    // 循环记录，查找不可用的分钟段
    const unavaliableMinutes: { start: number; end: number }[] = []
    records.forEach((item) => {
      const startHour = dayjs(item.startTime).hour()
      const startMinute = dayjs(item.startTime).minute()
      const endHour = dayjs(item.endTime).hour()
      const endMinute = dayjs(item.endTime).minute()

      const itemSitTime = {
        start: -1,
        end: -1
      }
      if (startHour === selectedStartTimeHour && startMinute !== 0) {
        // 开始时间在当前所选的小时
        itemSitTime.start = startMinute
      }

      if (endHour > selectedStartTimeHour) {
        itemSitTime.end = 59
      }
      if (endHour === dayjs(bookInfo.value.startTime).hour()) {
        // 结束时间在当前所选的小时之前
        itemSitTime.end = endMinute
      }
      unavaliableMinutes.push(itemSitTime)
    })
    // 将不可用的分钟段按照开始时间排序
    unavaliableMinutes.sort((a, b) => {
      return a.start - b.start
    })
    // 将排序后的记录进行合并，即如果有endTime===startTime，合并为一个记录
    const mergedUnavaliableMinutes: { start: number; end: number }[] = []
    unavaliableMinutes.forEach((item) => {
      if (mergedUnavaliableMinutes.length === 0) {
        mergedUnavaliableMinutes.push(item)
      } else {
        const last = mergedUnavaliableMinutes[mergedUnavaliableMinutes.length - 1]
        if (last.end === item.start) {
          last.end = item.end
        } else {
          mergedUnavaliableMinutes.push(item)
        }
      }
    })

    // 循环分钟，查找可用的分钟
    for (let i = 0; i < 60; i++) {
      const found = unavaliableMinutes.find((item) => {
        return item.start < i && item.end > i
      })
      if (found) {
        continue
      }
      list.push(i)
    }
  }

  return list
})
const endMinuteList = computed(() => {
  const list = []

  // 如果选择的结束时间小时和开始时间小时相同，则需要过滤掉开始时间之前的分钟
  let min = 0
  let max = 59
  if (dayjs(bookInfo.value.startTime).hour() === dayjs(bookInfo.value.endTime).hour()) {
    min = dayjs(bookInfo.value.startTime).minute()
  }

  // 找到同一天的预约记录，且开始时间在当前所选的小时之后的记录
  const records = reservationList.value.filter((item) => {
    const isSameDay = dayjs(item.startTime).startOf("date").valueOf() === dayjs(bookStartLimit.value).startOf("date").valueOf()
    if (!isSameDay) {
      return false
    }

    if (bookInfo.value.id === item.id) {
      return false
    }
    return item.startTime > bookInfo.value.startTime && dayjs(item.startTime).hour() === dayjs(bookInfo.value.endTime).hour()
  })
  if (records.length > 0) {
    // 获取最小的开始时间
    const minStartTime = records.reduce((prev, cur) => {
      if (prev.startTime < cur.startTime) {
        return prev
      }
      return cur
    }).startTime
    max = dayjs(minStartTime).minute()
  }

  for (let i = min; i <= max; i++) {
    list.push(i)
  }

  return list
})

// 删除
const deleteReservation = () => {
  if (bookInfo.value.id === "") {
    return
  }
  deleteMeetingRoomReservation(bookInfo.value.id).then((res) => {
    if (res) {
      message.success("删除成功")
      bookDrawerActive.value = false
      refresh()
    }
  })
}

const refresh = () => {
  getMeetingRoomReservationList(props.room.id).then((res) => {
    reservationList.value = res || []
  })
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-calendar v-model:value="selectedDate" :is-date-disabled="isDateDisabled" #="{ year, month, date }">
    <div class="relative">
      <n-button class="absolute right-0 top--28px" size="tiny" type="primary" v-if="canBook(year, month, date)" @click="book(year, month, date)">
        预约
      </n-button>
      <n-space vertical>
        <div v-for="item in dateReservationList(year, month, date)" :key="item.id">
          <n-tooltip>
            <template #trigger>
              <div class="text-truncate">
                {{ item.subject }}
              </div>
            </template>
            <div>主题：{{ item.subject }}</div>
            <div>与会人员：{{ item.participants || "未填写" }}</div>
            <div>会议时间：{{ dayjs(item.startTime).format("HH:mm") }} - {{ dayjs(item.endTime).format("HH:mm") }}</div>
            <div>
              <!-- 修改 -->
              <n-button v-if="item.bookerId === userStore.user.id" size="tiny" type="info" @click="showEdit(item)">
                <template #icon>
                  <n-icon>
                    <Edit />
                  </n-icon>
                </template>
              </n-button>
            </div>
          </n-tooltip>
        </div>
      </n-space>
    </div>
  </n-calendar>
  <n-drawer v-model:show="bookDrawerActive" :width="502">
    <n-drawer-content :title="dayjs(bookStartLimit).format('YYYY-MM-DD')">
      <n-form ref="bookFormRef" :model="bookInfo" :rules="bookRules" label-width="auto" label-placement="left">
        <n-grid :cols="2" :x-gap="16">
          <n-form-item-gi path="startTime" label="开始时间">
            <n-time-picker
              v-model:value="bookInfo.startTime"
              :min="bookStartLimit"
              :max="bookEndLimit"
              format="HH:mm"
              :hours="startHourList"
              :minutes="startMinuteList"
              @update:value="() => (bookInfo.endTime = bookInfo.startTime + 60 * 1000)"
            />
          </n-form-item-gi>
          <n-form-item-gi path="endTime" label="结束时间">
            <n-time-picker
              v-model:value="bookInfo.endTime"
              :min="bookStartLimit"
              :max="bookEndLimit"
              format="HH:mm"
              :hours="endHourList"
              :minutes="endMinuteList"
            />
          </n-form-item-gi>
          <n-form-item-gi :span="2" path="subject" label="主题">
            <n-input v-model:value="bookInfo.subject" placeholder="请输入主题" />
          </n-form-item-gi>
          <n-form-item-gi :span="2" label="与会人员">
            <n-input v-model:value="bookInfo.participants" placeholder="请输入与会人员" />
          </n-form-item-gi>
        </n-grid>
      </n-form>
      <template #footer>
        <n-space justify="space-between" class="w-full">
          <n-popconfirm @positive-click="deleteReservation">
            <template #trigger>
              <n-button size="small" type="error">删除</n-button>
            </template>
            确认要删除吗？
          </n-popconfirm>
          <n-button size="small" type="primary" @click="submitBook">提交</n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
