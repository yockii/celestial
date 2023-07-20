<script setup lang="ts">
import { ProjectIssue } from "@/types/project"
import Vditor from "vditor"
import "vditor/dist/index.css"
import { onMounted } from "vue"

const props = defineProps<{
  issue: ProjectIssue
  drawerActive: boolean
}>()
const emit = defineEmits(["update:drawerActive"])

onMounted(() => {
  const div = document.getElementById("content") as HTMLDivElement
  if (div) {
    Vditor.preview(div, props.issue.content || "")
  }
})
</script>

<template>
  <n-drawer :show="drawerActive" :default-height="600" resizable placement="top" @update:show="(show:boolean) => emit('update:drawerActive', show)">
    <n-drawer-content closable :title="issue.title">
      <n-space vertical>
        <div id="content"></div>
        <div>问题原因: {{ issue.issueCause }}</div>
        <div>解决方法：{{ issue.solveMethod }}</div>
      </n-space>
    </n-drawer-content>
  </n-drawer>
</template>
