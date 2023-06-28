<script setup lang="ts">
import MyProjectList from "@/components/project/MyProjectList.vue"
import TestList from "./list/index.vue"

const TestListInst = ref<typeof TestList>()

const selectedProjectId = ref("")
const handleSelectedProjectIdUpdate = (id: string | undefined) => {
  selectedProjectId.value = id || ""
  TestListInst.value?.projectSelected(id)
}
</script>

<template>
  <n-grid :cols="8" x-gap="16">
    <n-gi :span="1">
      <MyProjectList :selected-project-id="selectedProjectId" @update:selected-project-id="handleSelectedProjectIdUpdate" />
    </n-gi>
    <n-gi :span="7">
      <TestList />
    </n-gi>
  </n-grid>
</template>
