<script setup lang="ts">
import dayjs from "dayjs"
import { Checkmark, Edit } from "@vicons/carbon"
import { CancelOutlined } from "@vicons/material"
const props = defineProps<{
  data: number | undefined
}>()
const emit = defineEmits(["update:data"])
const isEdit = ref(false)
const handleDataUpdate = (value: number | null) => {
  emit("update:data", value ? value : 0)
  isEdit.value = false
}
const editableData = ref(0)
onMounted(() => {
  editableData.value = props.data ? props.data : new Date().getTime()
})
</script>

<template>
  <n-space justify="space-between">
    <template v-if="isEdit">
      <n-date-picker v-model:value="editableData" type="datetime" clearable />
      <n-button-group>
        <n-tooltip>
          <template #trigger>
            <n-button size="small" circle type="primary" style="width: 60px" @click="handleDataUpdate(editableData)">
              <template #icon>
                <n-icon>
                  <Checkmark />
                </n-icon>
              </template>
            </n-button>
          </template>
          确认
        </n-tooltip>
        <n-tooltip>
          <template #trigger>
            <n-button size="small" circle style="width: 60px" @click="isEdit = false">
              <template #icon>
                <n-icon>
                  <CancelOutlined />
                </n-icon>
              </template>
            </n-button>
          </template>
          取消
        </n-tooltip>
      </n-button-group>
    </template>
    <template v-else>
      <span v-show="!isEdit">
        {{ data === 0 ? "无" : dayjs(data).format("YYYY-MM-DD") }}
      </span>
      <n-tooltip v-show="!isEdit">
        <template #trigger>
          <n-button type="small" text @click="isEdit = true">
            <template #icon>
              <n-icon>
                <Edit />
              </n-icon>
            </template>
          </n-button>
        </template>
        编辑
      </n-tooltip>
    </template>
  </n-space>
</template>
