<script setup lang="ts">
import { ProjectTask, ProjectTaskMember } from "@/types/project"
import { developDoneProjectTask, confirmProjectTask } from "@/service/api"
import { useMessage } from "naive-ui"
const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  task: ProjectTask
  data: ProjectTaskMember
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const taskMemberInfo = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectTaskMember) => {
    emit("update:data", val)
  }
})

const realDuration = ref(0)

const estimateDuration = computed({
  get: () => {
    return realDuration.value / 3600
  },
  set: (val: number) => {
    realDuration.value = val * 3600
  }
})
const actualDuration = computed({
  get: () => {
    return realDuration.value / 3600
  },
  set: (val: number) => {
    realDuration.value = val * 3600
  }
})

const drawerTitle = computed(() => {
  return props.task.name + "工时记录"
})

const handleCommitData = () => {
  if (taskMemberInfo.value.status === 1) {
    if (estimateDuration.value === 0) {
      message.error("预计工时不能为0")
      return
    }
    confirmProjectTask(props.task.id, realDuration.value).then((res) => {
      if (res) {
        message.success("提交成功")
        emit("refresh")
        emit("update:drawerActive", false)
      }
    })
  } else if (taskMemberInfo.value.status === 3 || taskMemberInfo.value.status === 5) {
    if (actualDuration.value === 0) {
      message.error("使用工时不能为0")
      return
    }
    developDoneProjectTask(props.task.id, realDuration.value).then((res) => {
      if (res) {
        message.success("提交成功")
        emit("refresh")
        emit("update:drawerActive", false)
      }
    })
  }
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show:boolean) => $emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="taskMemberInfo" label-placement="left" label-width="100px">
        <n-form-item label="预计工时" v-if="taskMemberInfo.status === 1">
          <n-input-number v-model:value="estimateDuration" min="0">
            <template #suffix>小时</template>
          </n-input-number>
        </n-form-item>
        <n-form-item label="使用工时" v-if="taskMemberInfo.status === 3 || taskMemberInfo.status === 5">
          <n-input-number v-model:value="actualDuration" min="0">
            <template #suffix>小时</template>
          </n-input-number>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button type="primary" @click="handleCommitData">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
