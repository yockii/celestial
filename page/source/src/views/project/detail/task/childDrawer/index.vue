<script setup lang="ts">
import { deleteProjectTask } from "@/service/api/project/projectTask"
import { ProjectTask, ProjectTaskCondition } from "@/types/project"
import List from "../list/index.vue"
import Drawer from "../drawer/index.vue"

const message = useMessage()
const listComp = ref<typeof List>()

const props = defineProps<{
  drawerActive: boolean
  data: ProjectTask
  height?: number
}>()
const emit = defineEmits(["update:drawerActive"])

const condition = ref<ProjectTaskCondition>({
  projectId: props.data.projectId,
  parentId: props.data.id,
  offset: 0,
  limit: 10
})
watch(
  () => props.data,
  (val) => {
    condition.value.projectId = val.projectId
    condition.value.parentId = val.id
    listComp.value?.refresh()
  }
)

const editDrawerActive = ref(false)
const currentData = ref<ProjectTask>({
  id: "",
  name: "",
  projectId: ""
})

const handleEditData = (row: ProjectTask) => {
  currentData.value = Object.assign({}, row)
  editDrawerActive.value = true
}
const handleAddProjectTask = (
  row: ProjectTask = {
    id: "",
    name: "",
    projectId: props.data.projectId,
    parentId: props.data.id,
    moduleId: props.data.moduleId,
    requirementId: props.data.requirementId,
    priority: props.data.priority,
    taskDesc: props.data.taskDesc
  }
) => {
  if (row.id === "") {
    currentData.value = row
  } else {
    currentData.value = {
      id: "",
      name: "",
      projectId: row.projectId,
      parentId: row.id,
      moduleId: row.moduleId,
      requirementId: row.requirementId,
      priority: row.priority,
      taskDesc: row.taskDesc
    }
  }
  editDrawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectTask(id).then((res) => {
    if (res) {
      message.success("删除成功")
      listComp.value?.refresh()
    }
  })
}

// 子任务抽屉
const childDrawerActive = ref(false)
const handleShowChild = (row: ProjectTask) => {
  currentData.value = row
  childDrawerActive.value = true
}

const loadData = () => {
  // listComp.value?.refresh()
}
defineExpose({
  refresh: loadData
})
</script>

<template>
  <n-drawer
    :show="drawerActive"
    placement="top"
    :default-height="height ? height : 600"
    resizable
    :on-update:show="(show: boolean) => emit('update:drawerActive', show)"
  >
    <n-drawer-content :title="data.name + '的子任务'" closable>
      <n-grid :cols="1" y-gap="8">
        <n-gi>
          <n-space justify="space-between">
            <span></span>
            <n-button type="primary" @click="handleAddProjectTask()">新建子任务</n-button>
          </n-space>
        </n-gi>
        <n-gi>
          <list
            ref="listComp"
            :condition="condition"
            useTree
            @edit="handleEditData"
            @delete="handleDeleteData"
            @new-child="handleAddProjectTask"
            @showChild="handleShowChild"
          />
        </n-gi>
      </n-grid>
    </n-drawer-content>
  </n-drawer>

  <drawer v-model:drawer-active="editDrawerActive" v-model:data="currentData" @refresh="listComp?.refresh" />
</template>
