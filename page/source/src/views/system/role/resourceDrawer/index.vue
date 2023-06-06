<script setup lang="ts">
import { getResourceList } from "@/service/api/resource"
import { assignResource } from "@/service/api/role"
import { Resource } from "@/types/user"

const emit = defineEmits(["update:drawerActive"])
const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  roleId: string
  roleResourceCodeList: string[]
}>()
watch(
  () => props.roleResourceCodeList,
  () => {
    resetCurrentData()
  }
)
const resources = ref<Resource[]>([])
const resourceCodeList = ref<string[]>([])
const resourceTree = computed(() => {
  // 先获取顶层节点
  const topNodes = resources.value.filter((item) => !item.resourceCode.includes(":"))
  // 递归获取子节点
  const getChildren = (node: Resource) => {
    const children = resources.value.filter(
      (item) => item.resourceCode.includes(node.resourceCode + ":") && item.resourceCode.replace(node.resourceCode + ":", "").split(":").length === 1
    )
    if (children.length > 0) {
      node.children = children
      children.forEach((item) => {
        getChildren(item)
      })
    }
  }
  topNodes.forEach((item) => {
    getChildren(item)
  })
  return topNodes
})

onMounted(() => {
  getResourceList({
    offset: -1,
    limit: -1
  }).then((res) => {
    resources.value = res.items
  })
})

const resetCurrentData = () => {
  resourceCodeList.value = props.roleResourceCodeList
}
const handleCommitData = () => {
  assignResource(props.roleId, resourceCodeList.value).then((res) => {
    if (res) {
      message.success("更新成功")
      emit("update:drawerActive", false)
    }
  })
}
const updateCheckedKeys = (keys: string[]) => {
  resourceCodeList.value = keys
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show: boolean) => $emit('update:drawerActive', show)">
    <n-drawer-content title="分配资源">
      <n-tree
        block-line
        checkable
        check-on-click
        :data="resourceTree"
        :checked-keys="resourceCodeList"
        key-field="resourceCode"
        label-field="resourceName"
        children-field="children"
        @update:checked-keys="updateCheckedKeys"
      />

      <template #footer>
        <n-button class="mr-a" @click="resetCurrentData">重置</n-button>
        <n-button size="small" type="primary" @click="handleCommitData" v-resource-code="'system:role:dispatchResources'">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
