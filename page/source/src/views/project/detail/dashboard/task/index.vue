<script setup lang="ts">
import { Project } from "@/types/project"
defineProps<{
  project: Project
}>()
</script>

<template>
  <n-text tag="div" class="text-1.2em op-90 font-500 mb-10px">工作任务概览</n-text>
</template>

<style scoped></style>
