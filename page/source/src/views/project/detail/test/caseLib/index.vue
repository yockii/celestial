<script setup lang="ts">
import { getAllCommonTestCaseListOnlyShow } from "@/service/api"
import { useProjectStore } from "@/store/project"
import { CommonTestCase } from "@/types/asset"
import { ProjectTestCase, ProjectTestCaseItem } from "@/types/project"
import { TreeOption } from "naive-ui"
import { storeToRefs } from "pinia"

const props = defineProps<{
  selectedTestCase?: ProjectTestCase[]
  selectedTestCaseItem?: ProjectTestCaseItem[]
}>()
const emit = defineEmits(["update:selectedTestCase", "update:selectedTestCaseItem"])

const { project } = storeToRefs(useProjectStore())

const commonTestCaseList = ref<CommonTestCase[]>([])
const data = computed(() => {
  // commonTestCaseList 转化为 TreeOption[]
  return commonTestCaseList.value.map((item) => {
    const treeItem: TreeOption = {
      label: item.name,
      key: item.id,
      children: []
    }
    if (item.items) {
      for (let i = 0; i < item.items.length; i++) {
        const child = item.items[i]
        const childTreeItem: TreeOption = {
          label: child.content,
          key: child.id
        }
        treeItem.children?.push(childTreeItem)
      }
    }
    return treeItem
  })
})
const caseChecked = (keys: Array<string>) => {
  if (props.selectedTestCase) {
    const selectedTestCase: ProjectTestCase[] = []
    for (let i = 0; i < keys.length; i++) {
      const key = keys[i]
      const commonTestCase = commonTestCaseList.value.find((item) => item.id === key)
      if (commonTestCase) {
        const projectTestCase: ProjectTestCase = {
          id: commonTestCase.id,
          projectId: project.value.id,
          name: commonTestCase.name,
          items: []
        }
        if (commonTestCase.items) {
          const selectedItems = commonTestCase.items.filter((item) => keys.includes(item.id))
          const transferedItems = selectedItems.map((item) => {
            const projectTestCaseItem: ProjectTestCaseItem = {
              id: item.id,
              name: item.content,
              content: item.content
            }
            return projectTestCaseItem
          })
          projectTestCase.items && projectTestCase.items.push(...transferedItems)
        }
        selectedTestCase.push(projectTestCase)
      }
    }
    emit("update:selectedTestCase", selectedTestCase)
  }
  if (props.selectedTestCaseItem) {
    const selectedTestCaseItem: ProjectTestCaseItem[] = []
    for (let i = 0; i < commonTestCaseList.value.length; i++) {
      const commonTestCase = commonTestCaseList.value[i]
      if (commonTestCase.items) {
        const selectedItems = commonTestCase.items.filter((item) => keys.includes(item.id))
        const transferedItems = selectedItems.map((item) => {
          const projectTestCaseItem: ProjectTestCaseItem = {
            id: item.id,
            name: item.content,
            content: item.content
          }
          return projectTestCaseItem
        })
        selectedTestCaseItem.push(...transferedItems)
      }
    }
    emit("update:selectedTestCaseItem", selectedTestCaseItem)
  }
}

const refresh = () => {
  // 获取 commonTestCaseList
  getAllCommonTestCaseListOnlyShow().then((res) => {
    commonTestCaseList.value = res
  })
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-tree block-line block-node checkable @update:checked-keys="caseChecked" :data="data" animated />
</template>
