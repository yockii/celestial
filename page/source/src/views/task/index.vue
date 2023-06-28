<script setup lang="ts">
import MyProjectList from "@/components/project/MyProjectList.vue"
import TaskList from "./list/index.vue"

const TaskListInst = ref<typeof TaskList>()
const selectedProjectId = ref("")
const handleSelectedProjectIdUpdate = (id: string | undefined) => {
  selectedProjectId.value = id || ""
  TaskListInst.value?.projectSelected(id)
}
</script>

<template>
  <n-grid :cols="8" x-gap="16">
    <n-gi :span="1">
      <MyProjectList :selected-project-id="selectedProjectId" @update:selected-project-id="handleSelectedProjectIdUpdate" />
    </n-gi>
    <n-gi :span="7">
      <TaskList ref="TaskListInst" :selected-project-id="selectedProjectId" />
    </n-gi>
  </n-grid>
</template>
