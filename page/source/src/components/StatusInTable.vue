<script setup lang="ts">
import { Checkmark, Edit } from "@vicons/carbon"
import { CancelOutlined } from "@vicons/material"
const props = defineProps<{
  data: number | undefined
  options: {
    label: string
    value: number
  }[]
}>()
const emit = defineEmits(["update:data"])
const isEdit = ref(false)
const handleDataUpdate = (value: number | null) => {
  emit("update:data", value ? value : 0)
  isEdit.value = false
}
const editableData = ref(0)
// watch(
//   () => props.data,
//   (value) => {
//     editableData.value = value
//   }
// )
onMounted(() => {
  editableData.value = props.data || 0
})
</script>

<template>
  <n-space justify="space-between">
    <template v-if="isEdit">
      <n-select v-model:value="editableData" :options="options" :consistent-menu-width="false" />
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
      <span>
        {{ options.find((option) => option.value === data)?.label || "未知" }}
      </span>
      <n-tooltip>
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
