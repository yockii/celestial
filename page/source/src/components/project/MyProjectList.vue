<script setup lang="ts">
import { storeToRefs } from "pinia"
import { useUserStore } from "@/store/user"
import { Project } from "@/types/project"
import { onMounted } from "vue"
import { getMyProjectList } from "@/service/api"

defineProps<{
  selectedProjectId: string
}>()
const emit = defineEmits(["update:selectedProjectId"])

const { myProjectList } = storeToRefs(useUserStore())
const projectTree = computed(() => {
  // 将 myProjectList 转换为树形结构, 使用 label/key/children
  type ProjectTree = {
    label: string
    key: string
    children: ProjectTree[]
    isLeaf?: boolean
  }
  const topProject: ProjectTree[] = myProjectList.value
    .filter((p: Project) => !p.parentId)
    .map((p: Project) => {
      return {
        label: p.name,
        key: p.id,
        children: []
      }
    })
  // 递归赋值子项目
  const setChildren = (project: ProjectTree) => {
    const children = myProjectList.value
      .filter((p: Project) => p.parentId === project.key)
      .map((p: Project) => {
        return {
          label: p.name,
          key: p.id,
          children: []
        }
      })
    project.children = children
    project.isLeaf = children.length === 0
    children.forEach((p) => setChildren(p))
  }

  topProject.forEach((p) => setChildren(p))
  return topProject
})

const handleSelect = (key: string[]) => {
  emit("update:selectedProjectId", key[0])
}

onMounted(() => {
  if (myProjectList.value.length === 0) {
    getMyProjectList({}).then((res) => {
      myProjectList.value = res
    })
  }
})
</script>

<template>
  <n-tree block-line :data="projectTree" selectable :on-update:selected-keys="handleSelect" />
</template>
