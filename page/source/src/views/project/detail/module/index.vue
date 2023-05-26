<script setup lang="ts">
import { Project, ProjectModule, ProjectModuleCondition, ProjectRequirement } from "@/types/project"
import { ComputedRef, computed, ref, RendererElement, RendererNode } from "vue"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage } from "naive-ui"
import { addProjectModule, deleteProjectModule, getProjectModuleList, updateProjectModule } from "@/service/api/projectModule"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import Drawer from "../requirement/drawer/index.vue"

const message = useMessage()
const projectStore = useProjectStore()
const props = defineProps<{
  project: Project
}>()

type CombinedModule = {
  lv1Column: string
  lv2Column: string
  lv3Column: string
  lv4Column: string
  lv1Module: ProjectModule
  lv2Module?: ProjectModule
  lv3Module?: ProjectModule
  lv4Module?: ProjectModule
  [key: string]: unknown
}
const { modules: list } = storeToRefs(projectStore)
// const list: Ref<ProjectModule[]> = ref([])
const data: ComputedRef<CombinedModule[]> = computed(() => {
  const combinedModules: CombinedModule[] = []
  const lv1Modules = list.value.filter((item) => !item.parentId || item.parentId === "")

  for (let i = 0; i < lv1Modules.length; i++) {
    let module = lv1Modules[i]
    if (module.parentId && module.parentId !== "") {
      continue
    }
    // 查找所有该模块的子模块
    const subModules = list.value.filter((m) => m.parentId === module.id)
    if (subModules.length > 0) {
      for (let j = 0; j < subModules.length; j++) {
        let subModule = subModules[j]
        // 查找所有该模块的子模块
        const subSubModules = list.value.filter((m) => m.parentId === subModule.id)
        if (subSubModules.length > 0) {
          for (let k = 0; k < subSubModules.length; k++) {
            let subSubModule = subSubModules[k]
            // 查找所有该模块的子模块
            const subSubSubModules = list.value.filter((m) => m.parentId === subSubModule.id)
            if (subSubSubModules.length > 0) {
              for (let l = 0; l < subSubSubModules.length; l++) {
                let subSubSubModule = subSubSubModules[l]
                combinedModules.push({
                  lv1Column: module.name,
                  lv2Column: subModule.name,
                  lv3Column: subSubModule.name,
                  lv4Column: subSubSubModule.name,
                  lv1Module: module,
                  lv2Module: subModule,
                  lv3Module: subSubModule,
                  lv4Module: subSubSubModule
                })
              }
            } else {
              combinedModules.push({
                lv1Column: module.name,
                lv2Column: subModule.name,
                lv3Column: subSubModule.name,
                lv4Column: "",
                lv1Module: module,
                lv2Module: subModule,
                lv3Module: subSubModule
              })
            }
          }
        } else {
          combinedModules.push({
            lv1Column: module.name,
            lv2Column: subModule.name,
            lv3Column: "",
            lv4Column: "",
            lv1Module: module,
            lv2Module: subModule
          })
        }
      }
    } else {
      combinedModules.push({
        lv1Column: module.name,
        lv2Column: "",
        lv3Column: "",
        lv4Column: "",
        lv1Module: module
      })
    }
  }

  return combinedModules
})
type ModuleOption = {
  value: string
  label: string
  children?: ModuleOption[]
}
const moduleOptions = computed(() => {
  const options: ModuleOption[] = []
  const lv1Modules = list.value.filter((item) => !item.parentId || item.parentId === "")
  for (let i = 0; i < lv1Modules.length; i++) {
    const lv1Module = lv1Modules[i]
    const option: ModuleOption = {
      value: lv1Module.id,
      label: lv1Module.name,
      children: []
    }
    const lv2Modules = list.value.filter((item) => item.parentId === lv1Module.id)
    for (let j = 0; j < lv2Modules.length; j++) {
      const lv2Module = lv2Modules[j]
      const lv2Option: ModuleOption = {
        value: lv2Module.id,
        label: lv2Module.name,
        children: []
      }
      const lv3Modules = list.value.filter((item) => item.parentId === lv2Module.id)
      const lv3Options = lv3Modules.map((item) => {
        return {
          value: item.id,
          label: item.name
        }
      })
      lv2Option.children = lv3Options
      option.children?.push(lv2Option)
    }
    options.push(option)
  }
  return options
})
const countModules = (lv: number, cm: CombinedModule) => {
  const columnName = `lv${lv}Column`
  const moduleName = `lv${lv}Module`
  if (cm[columnName] === "") {
    return 1
  }
  const l = data.value.filter(
    (item) =>
      item[moduleName] &&
      cm[moduleName] &&
      (item[moduleName] as ProjectModule).parentId === (cm[moduleName] as ProjectModule).parentId &&
      item[columnName] === cm[columnName]
  ).length
  return l
}
const columns = [
  {
    title: "一级模块",
    key: "lv1Column",
    rowSpan: (rowData: CombinedModule) => countModules(1, rowData)
  },
  {
    title: "二级模块",
    key: "lv2Column",
    rowSpan: (rowData: CombinedModule) => countModules(2, rowData)
  },
  {
    title: "三级模块",
    key: "lv3Column",
    rowSpan: (rowData: CombinedModule) => countModules(3, rowData)
  },
  {
    title: "四级模块",
    key: "lv4Column",
    rowSpan: (rowData: CombinedModule) => countModules(4, rowData)
  },
  {
    title: "状态",
    key: "status",
    render: (rowData: CombinedModule) => {
      let module = rowData.lv4Module ? rowData.lv4Module : rowData.lv3Module ? rowData.lv3Module : rowData.lv2Module ? rowData.lv2Module : rowData.lv1Module
      switch (module.status) {
        case -1:
          return "废弃"
        case 1:
          return "待评审"
        case 2:
          return "评审通过"
        case 3:
          return "评审不通过"
        case 9:
          return "已完成"
      }
      return "未知"
    }
  },
  {
    title: "操作",
    key: "operation",
    render: (rowData: CombinedModule) => {
      let module = rowData.lv4Module ? rowData.lv4Module : rowData.lv3Module ? rowData.lv3Module : rowData.lv2Module ? rowData.lv2Module : rowData.lv1Module
      const btnGroup: globalThis.VNode<RendererNode, RendererElement, { [key: string]: unknown }>[] = []
      if (module.status === 2) {
        btnGroup.push(
          h(
            NButton,
            {
              size: "small",
              onClick: () => handleNewRequirement(module)
            },
            {
              default: () => "新增需求"
            }
          )
        )
      }
      btnGroup.push(
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            type: "primary",
            onClick: () => handleEditData(module)
          },
          {
            default: () => "编辑"
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteData(module.id)
          },
          {
            default: () => "确认删除",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  type: "error"
                },
                {
                  default: () => "删除"
                }
              )
          }
        )
      )

      return h(NButtonGroup, {}, () => btnGroup)
    }
  }
]
const handleNewModule = () => {
  resetCurrentData()
  drawerActive.value = true
}
const handleEditData = (module: ProjectModule) => {
  resetCurrentData(module)
  drawerActive.value = true
}
const handleDeleteData = (id: string) => {
  deleteProjectModule(id).then(() => {
    refresh()
  })
}
const loading = ref(false)
const condition = ref<ProjectModuleCondition>({
  offset: -1,
  limit: -1,
  projectId: props.project.id
})
const refresh = () => {
  loading.value = true
  getProjectModuleList(condition.value)
    .then((res) => {
      list.value = res.items
    })
    .finally(() => {
      loading.value = false
    })
}

// 抽屉
const drawerActive = ref(false)
const currentData = ref<ProjectModule>({ id: "", name: "", projectId: "", status: 0 })
const resetCurrentData = (data: ProjectModule | null = null) => {
  if (data) {
    currentData.value = data
  } else {
    currentData.value = {
      alias: "",
      childrenCount: 0,
      id: "",
      name: "",
      parentId: "",
      projectId: props.project.id,
      remark: "",
      status: 1
    }
  }
}
const isUpdate = computed(() => !!currentData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改模块" : "新增模块"))
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateProjectModule(currentData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addProjectModule(currentData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          refresh()
          drawerActive.value = false
        }
      })
    }
  })
}
const rules = {
  name: {
    required: true,
    message: "请输入模块名称",
    min: 2,
    max: 20,
    trigger: "blur"
  }
}

// TODO 新增需求
const showNewRequirement = ref(false)
const requirementData = ref<ProjectRequirement>({
  id: "",
  name: "",
  projectId: props.project.id
})
const handleNewRequirement = (module: ProjectModule) => {
  requirementData.value.moduleId = module.id
  showNewRequirement.value = true
}

// 加载动作
onMounted(() => {
  if (list.value.length === 0) {
    refresh()
  }
})
</script>

<template>
  <n-grid :cols="1" y-gap="12">
    <n-gi class="flex flex-justify-end">
      <n-button type="primary" @click="handleNewModule">新增模块</n-button>
    </n-gi>
    <n-gi>
      <n-data-table :columns="columns" :data="data" :single-line="false" />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="currentData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="模块名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入模块名称" />
        </n-form-item>
        <n-form-item label="模块别名">
          <n-input v-model:value="currentData.alias" placeholder="请输入模块别名" />
        </n-form-item>
        <n-form-item label="模块描述">
          <n-input v-model:value="currentData.remark" placeholder="请输入模块描述" />
        </n-form-item>
        <n-form-item label="所属模块">
          <n-cascader
            v-model:value="currentData.parentId"
            placeholder="请选择所属模块"
            :options="moduleOptions"
            expand-trigger="click"
            check-strictly="all"
            show-path
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData()">重置</n-button>
        <n-button type="primary" @click="handleCommitData">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>

  <drawer :project-id="project.id" v-model:drawer-active="showNewRequirement" :data="requirementData" />
</template>

<style scoped></style>
