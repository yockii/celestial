<script setup lang="ts">
import { Department } from "@/types/user"
import { getDepartmentList } from "@/service/api"

const emit = defineEmits(["update:department-selected"])

const topDepartments = ref<Department[]>([])

const handleLoadChildren = (node: Department) => {
  return new Promise<void>((resolve) => {
    getDepartmentList({ parentId: node.id, offset: -1, limit: -1 }).then((res) => {
      node.children = res.items.map((item) => {
        item.isLeaf = !item.childCount || item.childCount === 0
        return item
      })
      resolve()
    })
  })
}
const handleSelectedKeysChange = (keys: string[]) => {
  emit("update:department-selected", keys && keys.length > 0 ? keys[0] : "")
}

onMounted(() => {
  getDepartmentList({ onlyParent: true, offset: -1, limit: -1 }).then((res) => {
    topDepartments.value = res.items.map((item) => {
      item.isLeaf = !item.childCount || item.childCount === 0
      return item
    })
  })
})
</script>

<template>
  <n-tree
    blck-line
    block-node
    cancelable
    :data="topDepartments"
    selectable
    :on-load="handleLoadChildren"
    key-field="id"
    label-field="name"
    children-field="children"
    :on-update:selected-keys="handleSelectedKeysChange"
  />
</template>
