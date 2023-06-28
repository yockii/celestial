<script setup lang="ts">
import MyProjectList from "@/components/project/MyProjectList.vue"
import TaskList from "./list/index.vue"

const selectedProjectId = ref("")
</script>

<template>
  <n-grid :cols="8" x-gap="16">
    <n-gi :span="1">
      <MyProjectList @update:selected-project-id="selectedProjectId = $event" />
    </n-gi>
    <n-gi :span="7">
      <TaskList />
    </n-gi>
  </n-grid>
</template>
