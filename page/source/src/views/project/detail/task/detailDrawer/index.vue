<script setup lang="ts">
import { ProjectTask } from "@/types/project"
defineProps<{
  drawerActive: boolean
  data: ProjectTask
  height?: number
}>()
const emit = defineEmits(["update:drawerActive"])
</script>

<template>
  <n-drawer
    :show="drawerActive"
    placement="bottom"
    :default-height="height ? height : 600"
    resizable
    :on-update:show="(show: boolean) => emit('update:drawerActive', show)"
  >
    <n-drawer-content :title="data.name + '详情'" closable>
      <n-text>{{ data.taskDesc }}</n-text>
    </n-drawer-content>
  </n-drawer>
</template>
