<script setup lang="ts">
import { getProjectModuleList, deleteProjectTask, getProjectTask } from "@/service/api"
import { useProjectStore } from "@/store/project"
import { ProjectModule, ProjectTask, ProjectTaskCondition } from "@/types/project"
import { useMessage, NButton, NIcon } from "naive-ui"
import { storeToRefs } from "pinia"
import { Refresh } from "@vicons/tabler"
import Drawer from "./drawer/index.vue"
import ChildDrawer from "./childDrawer/index.vue"
import List from "./list/index.vue"

const message = useMessage()
const projectStore = useProjectStore()
const { project } = storeToRefs(useProjectStore())
const condition = ref<ProjectTaskCondition>({
  projectId: project.value.id,
  onlyParent: false
})

const listComp = ref<typeof List>()

const { modules, moduleTree } = storeToRefs(projectStore)
const treeSelected = (keys: string[], modules: ProjectModule[]) => {
  if (modules.length > 0) {
    condition.value.fullPath = modules[0].fullPath
    listComp.value?.refresh()
  }
}

const drawerActive = ref(false)
const currentData = ref<ProjectTask>({
  id: "",
  name: "",
  projectId: project.value.id
})

const handleEditData = (row: ProjectTask) => {
  currentData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectTask(id).then((res) => {
    if (res) {
      message.success("删除成功")
      listComp.value?.refresh()
    }
  })
}
const handleAddProjectTask = (
  row: ProjectTask = {
    id: "",
    name: "",
    projectId: project.value.id
  }
) => {
  if (row.id === "") {
    currentData.value = row
  } else {
    currentData.value = {
      id: "",
      name: "",
      projectId: project.value.id,
      parentId: row.id,
      moduleId: row.moduleId,
      requirementId: row.requirementId,
      priority: row.priority,
      taskDesc: row.taskDesc
    }
  }
  drawerActive.value = true
}

// 子任务抽屉
const childDrawer = ref<typeof ChildDrawer>()
const childDrawerActive = ref(false)
const handleShowChild = (row: ProjectTask) => {
  // 先获取任务详情，再打开抽屉
  getProjectTask(row.id).then((res) => {
    if (res) {
      currentData.value = res
      childDrawerActive.value = true
    }
  })
}

// 加载页面
const loadModules = () => {
  getProjectModuleList({
    projectId: project.value.id,
    offset: -1,
    limit: -1
  }).then((res) => {
    if (res) {
      modules.value = res.items || []
    }
  })
}

const reload = () => {
  // 如果功能模块列表为空, 则加载
  if (!moduleTree.value.length) {
    loadModules()
  }
  if (route.query.id) {
    condition.value = {
      id: route.query.id as string,
      projectId: project.value.id
    }
  } else {
    condition.value = {
      projectId: project.value.id
    }
  }
  nextTick(() => {
    listComp.value?.refresh()
    condition.value.id = ""
  })
  // listComp.value?.refresh()
}

onMounted(() => {
  reload()
})
const route = useRoute()
// onBeforeUpdate(() => {
//   reload()
// })
</script>

<template>
  <n-grid :cols="6" x-gap="16">
    <n-gi>
      <n-grid :cols="1" y-gap="16">
        <n-gi class="flex justify-between">
          <n-text class="font-bold">功能模块</n-text>
          <n-button circle size="small" @click="loadModules">
            <template #icon>
              <n-icon>
                <Refresh />
              </n-icon>
            </template>
          </n-button>
        </n-gi>
        <n-gi>
          <n-tree :data="moduleTree" key-field="id" label-field="name" children-field="children" :on-update:selected-keys="treeSelected" selectable />
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi :span="5">
      <n-grid :cols="1" y-gap="16">
        <n-gi>
          <n-space justify="space-between">
            <n-switch v-model:value="condition.onlyParent" :round="false" @update:value="listComp?.refresh" size="small">
              <template #checked>仅显示主任务</template>
              <template #unchecked>显示所有任务</template>
            </n-switch>
            <n-button type="primary" @click="handleAddProjectTask()" v-project-resource-code="'project:detail:task:add'">新增任务</n-button>
          </n-space>
        </n-gi>
        <n-gi v-project-resource-code="'project:detail:task:list'">
          <list
            ref="listComp"
            :condition="condition"
            @edit="handleEditData"
            @delete="handleDeleteData"
            @new-child="handleAddProjectTask"
            @show-child="handleShowChild"
          />
          <!-- <n-data-table
            size="small"
            remote
            :data="list"
            :loading="loading"
            :columns="columns"
            :pagination="paginationReactive"
            :row-key="(row: ProjectTask) => row.id"
            :on-update:page="handlePageChange"
            :on-update:page-size="handlePageSizeChange"
            :on-update:filters="handleFiltersChange"
          /> -->
        </n-gi>
      </n-grid>
    </n-gi>
  </n-grid>

  <drawer :project-id="project.id" v-model:drawer-active="drawerActive" v-model:data="currentData" @refresh="listComp?.refresh" />
  <child-drawer ref="childDrawer" v-model:drawer-active="childDrawerActive" :data="currentData" />
</template>
