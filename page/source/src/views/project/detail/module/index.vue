<script setup lang="ts">
import { ProjectModule, ProjectModuleCondition, ProjectRequirement } from "@/types/project"
import { ComputedRef, computed, ref, RendererElement, RendererNode } from "vue"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, NTooltip, NSpace, NIcon, NText } from "naive-ui"
import { addProjectModule, deleteProjectModule, getProjectModuleList, moduleReview, updateProjectModule } from "@/service/api"
import { storeToRefs } from "pinia"
import { useProjectStore } from "@/store/project"
import Drawer from "../requirement/drawer/index.vue"
import { useUserStore } from "@/store/user"
import { AiStatusComplete, AiStatusFailed, Delete, Edit } from "@vicons/carbon"
import { PlaylistAdd } from "@vicons/tabler"

const message = useMessage()
const userStore = useUserStore()
const projectStore = useProjectStore()
const { project } = storeToRefs(useProjectStore())

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

const handleReviewData = (module: ProjectModule, status: number) => {
  moduleReview(module.id, status).then((res) => {
    if (res) {
      message.success(module.name + "评审成功")
      refresh()
    } else {
      message.error(module.name + "状态修改")
    }
  })
}
const columns = [
  {
    title: "一级模块",
    key: "lv1Column",
    rowSpan: (rowData: CombinedModule) => countModules(1, rowData),
    render: (row: CombinedModule) =>
      h(
        NSpace,
        {
          justify: "space-between"
        },
        {
          default: () => [
            h(
              NTooltip,
              {},
              {
                default: () => row.lv1Module.remark || row.lv1Module.name,
                trigger: () => (markedModuleId.value !== row.lv1Module?.id ? row.lv1Column : h(NText, { type: "error" }, { default: () => row.lv1Column }))
              }
            ),
            h(
              NTooltip,
              {},
              {
                default: () => "编辑",
                trigger: () =>
                  h(
                    NButton,
                    {
                      size: "small",
                      type: "primary",
                      secondary: true,
                      disabled: !userStore.hasResourceCode("project:detail:module:update") || row.lv1Module.status > 1,
                      onClick: () => handleEditData(row.lv1Module)
                    },
                    {
                      default: () =>
                        h(NIcon, {
                          component: Edit
                        })
                    }
                  )
              }
            )
          ]
        }
      )
  },
  {
    title: "二级模块",
    key: "lv2Column",
    rowSpan: (rowData: CombinedModule) => countModules(2, rowData),
    render: (row: CombinedModule) =>
      h(
        NSpace,
        {
          justify: "space-between"
        },
        {
          default: () => [
            h(
              NTooltip,
              {},
              {
                default: () => row.lv2Module?.remark || row.lv2Module?.name,
                trigger: () => (markedModuleId.value !== row.lv2Module?.id ? row.lv2Column : h(NText, { type: "error" }, { default: () => row.lv2Column }))
              }
            ),
            row.lv2Module
              ? h(
                  NTooltip,
                  {},
                  {
                    default: () => "编辑",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "primary",
                          secondary: true,
                          disabled: !userStore.hasResourceCode("project:detail:module:update") || (row.lv2Module && row.lv2Module.status > 1),
                          onClick: () => (row.lv2Module ? handleEditData(row.lv2Module) : "")
                        },
                        {
                          default: () =>
                            h(NIcon, {
                              component: Edit
                            })
                        }
                      )
                  }
                )
              : ""
          ]
        }
      )
  },
  {
    title: "三级模块",
    key: "lv3Column",
    rowSpan: (rowData: CombinedModule) => countModules(3, rowData),
    render: (row: CombinedModule) =>
      h(
        NSpace,
        {
          justify: "space-between"
        },
        {
          default: () => [
            h(
              NTooltip,
              {},
              {
                default: () => row.lv3Module?.remark || row.lv3Module?.name,
                trigger: () => (markedModuleId.value !== row.lv3Module?.id ? row.lv3Column : h(NText, { type: "error" }, { default: () => row.lv3Column }))
              }
            ),
            row.lv3Module
              ? h(
                  NTooltip,
                  {},
                  {
                    default: () => "编辑",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "primary",
                          secondary: true,
                          disabled: !userStore.hasResourceCode("project:detail:module:update") || (row.lv3Module && row.lv3Module.status > 1),
                          onClick: () => (row.lv3Module ? handleEditData(row.lv3Module) : "")
                        },
                        {
                          default: () =>
                            h(NIcon, {
                              component: Edit
                            })
                        }
                      )
                  }
                )
              : ""
          ]
        }
      )
  },
  {
    title: "四级模块",
    key: "lv4Column",
    rowSpan: (rowData: CombinedModule) => countModules(4, rowData),
    render: (row: CombinedModule) =>
      h(
        NSpace,
        {
          justify: "space-between"
        },
        {
          default: () => [
            h(
              NTooltip,
              {},
              {
                default: () => row.lv4Module?.remark || row.lv4Module?.name,
                trigger: () => (markedModuleId.value !== row.lv4Module?.id ? row.lv4Column : h(NText, { type: "error" }, { default: () => row.lv4Column }))
              }
            ),
            row.lv4Module
              ? h(
                  NTooltip,
                  {},
                  {
                    default: () => "编辑",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "primary",
                          secondary: true,
                          disabled: !userStore.hasResourceCode("project:detail:module:update") || (row.lv4Module && row.lv4Module.status > 1),
                          onClick: () => (row.lv4Module ? handleEditData(row.lv4Module) : "")
                        },
                        {
                          default: () =>
                            h(NIcon, {
                              component: Edit
                            })
                        }
                      )
                  }
                )
              : ""
          ]
        }
      )
  },
  {
    title: "状态",
    key: "status",
    render: (rowData: CombinedModule) => {
      let module = rowData.lv4Module ? rowData.lv4Module : rowData.lv3Module ? rowData.lv3Module : rowData.lv2Module ? rowData.lv2Module : rowData.lv1Module
      switch (module.status) {
        case 1:
          // "待评审"
          return h(
            NSpace,
            {
              justify: "space-between"
            },
            {
              default: () => [
                "待评审",
                h(
                  NButtonGroup,
                  {},
                  {
                    default: () => [
                      h(
                        NTooltip,
                        {},
                        {
                          default: () => "评审通过",
                          trigger: () =>
                            h(
                              NButton,
                              {
                                size: "small",
                                type: "primary",
                                disabled: !userStore.hasResourceCode("project:detail:module:review"),
                                onClick: () => handleReviewData(module, 2)
                              },
                              {
                                default: () => h(NIcon, { component: AiStatusComplete })
                              }
                            )
                        }
                      ),
                      h(
                        NTooltip,
                        {},
                        {
                          default: () => "评审不通过",
                          trigger: () =>
                            h(
                              NButton,
                              {
                                size: "small",
                                type: "error",
                                disabled: !userStore.hasResourceCode("project:detail:module:review"),
                                onClick: () => handleReviewData(module, -1)
                              },
                              {
                                default: () => h(NIcon, { component: AiStatusFailed })
                              }
                            )
                        }
                      )
                    ]
                  }
                )
              ]
            }
          )
        case 2:
          return "评审通过"
        case -1:
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
            NTooltip,
            {},
            {
              default: () => "新增需求",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    disabled: !userStore.hasResourceCode("project:detail:requirement:add"),
                    onClick: () => handleNewRequirement(module)
                  },
                  {
                    default: () => h(NIcon, { component: PlaylistAdd })
                  }
                )
            }
          )
        )
      }
      btnGroup.push(
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteData(module.id)
          },
          {
            default: () => "确认删除",
            trigger: () =>
              h(
                NTooltip,
                {},
                {
                  default: () => "删除",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        disabled: !userStore.hasResourceCode("project:detail:module:delete"),
                        type: "error"
                      },
                      {
                        default: () => h(NIcon, { component: Delete })
                      }
                    )
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
  projectId: project.value.id
})
const refresh = () => {
  loading.value = true
  getProjectModuleList(condition.value)
    .then((res) => {
      list.value = res.items || []
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
      projectId: project.value.id,
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
  projectId: project.value.id
})
const handleNewRequirement = (module: ProjectModule) => {
  requirementData.value.moduleId = module.id
  showNewRequirement.value = true
}

const markedModuleId = computed(() => {
  const id = route.query.id as string
  return id
})

// 加载动作
const reload = () => {
  if (list.value.length === 0) {
    refresh()
  }
}
onMounted(() => {
  reload()
})
const route = useRoute()
onBeforeUpdate(() => {
  reload()
})
</script>

<template>
  <n-grid :cols="1" y-gap="12">
    <n-gi class="flex flex-justify-end">
      <n-button type="primary" @click="handleNewModule" v-resource-code="'project:detail:module:add'">新增模块</n-button>
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
