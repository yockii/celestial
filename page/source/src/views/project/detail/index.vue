<script setup lang="ts">
import { useRoute, useRouter } from "vue-router"
import { ref, onMounted } from "vue"
import { Project } from "@/types/project"
import { deleteProject, getProjectDetail, getProjectResourceCode, updateProject } from "@/service/api"
import { KeyboardBackspaceOutlined } from "@vicons/material"
import { SettingsServices, Delete } from "@vicons/carbon"
// import Dashboard from "./dashboard/index.vue"
// import Plan from "./plan/index.vue"
// import Module from "./module/index.vue"
// import Requirement from "./requirement/index.vue"
// import Task from "./task/index.vue"
// import Test from "./test/index.vue"
// import Issue from "./issue/index.vue"
// import Risk from "./risk/index.vue"
// import Asset from "./asset/index.vue"
import { FormInst, NButton } from "naive-ui"
import { useProjectStore } from "@/store/project"
import { storeToRefs } from "pinia"

const router = useRouter()
const route = useRoute()
const id = route.params.id as string
const tab = computed(() => route.meta.title as string)
const projectStore = useProjectStore()
const { project, resourceCodes } = storeToRefs(projectStore)

// 项目设置 ////////////////
const showSettings = ref<boolean>(false)
const copiedProject = ref<Project>({
  id: "",
  name: "",
  code: "",
  description: "",
  stageId: ""
})
const showProjectSettings = () => {
  resetUpdateProject()
  showSettings.value = true
}
const projectRules = {
  name: [
    { required: true, message: "请输入项目名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  code: [
    { required: true, message: "请输入项目代码", trigger: "blur" },
    { pattern: /^\D\w{2,19}$/, message: "长度在 3 到 20 个英文/数字/下划线，且不能以数字开头", trigger: "blur" }
  ]
}
const resetUpdateProject = () => {
  copiedProject.value = JSON.parse(JSON.stringify(project.value))
}
const formRef = ref<FormInst>()
const handleCommitProject = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (copiedProject.value) {
        updateProject(copiedProject.value as Project).then((res) => {
          if (res) {
            project.value = copiedProject.value
            showSettings.value = false
          }
        })
      }
    }
  })
}

// 删除项目 ////////////////////////////
const deleteProjectName = ref<string>("")
const doDeleteProject = () => {
  if (project.value && project.value?.name === deleteProjectName.value) {
    deleteProjectName.value = ""
    deleteProject(project.value?.id as string).then((res) => {
      if (res) {
        useMessage().success("项目删除成功")
        router.back()
      }
    })
  }
}

const handleChangeTab = (value: string | number) => {
  switch (value) {
    case "项目总览":
      router.push(`/project/detail/${id}/dashboard`)
      break
    case "项目计划":
      router.push(`/project/detail/${id}/plan`)
      break
    case "功能模块":
      router.push(`/project/detail/${id}/module`)
      break
    case "项目需求":
      router.push(`/project/detail/${id}/requirement`)
      break
    case "工作任务":
      router.push(`/project/detail/${id}/task`)
      break
    case "项目测试":
      router.push(`/project/detail/${id}/test`)
      break
    case "项目缺陷":
      router.push(`/project/detail/${id}/issue`)
      break
    case "项目变更":
      router.push(`/project/detail/${id}/change`)
      break
    case "项目风险":
      router.push(`/project/detail/${id}/risk`)
      break
    case "项目资产":
      router.push(`/project/detail/${id}/asset`)
      break
  }
}

onMounted(() => {
  getProjectDetail(id).then((res) => {
    project.value = res
  })
  getProjectResourceCode(id).then((res) => {
    resourceCodes.value = res
  })
})
</script>

<template>
  <n-layout class="ma--16px">
    <n-layout-header bordered class="h-48px py-4px px-20px">
      <n-grid :cols="24">
        <n-gi :span="2" class="flex flex-items-center">
          <n-icon class="text-1.5em cursor-pointer" @click="router.back()">
            <KeyboardBackspaceOutlined />
          </n-icon>
          <n-button text icon-placement="right" class="text-1.2em ml-16px">
            {{ project?.name }}
          </n-button>
        </n-gi>
        <template v-if="project?.id">
          <n-gi :span="16" :offset="2">
            <n-tabs id="project-tabs" :value="tab" type="line" justify-content="space-between" @update:value="handleChangeTab">
              <n-tab name="项目总览" v-project-resource-code="'project:detail'"></n-tab>
              <n-tab name="项目计划" v-project-resource-code="'project:detail:plan'"></n-tab>
              <n-tab name="功能模块" v-project-resource-code="'project:detail:module'"></n-tab>
              <n-tab name="项目需求" v-project-resource-code="'project:detail:requirement'"></n-tab>
              <n-tab name="工作任务" v-project-resource-code="'project:detail:task'"></n-tab>
              <n-tab name="项目测试" v-project-resource-code="'project:detail:test'"></n-tab>
              <n-tab name="项目缺陷" v-project-resource-code="'project:detail:issue'"></n-tab>
              <n-tab name="项目变更" v-project-resource-code="'project:detail:change'"></n-tab>
              <n-tab name="项目风险" v-project-resource-code="'project:detail:risk'"></n-tab>
              <n-tab name="项目资产" v-project-resource-code="'project:detail:asset'"></n-tab>
            </n-tabs>
          </n-gi>
          <n-gi :span="2" :offset="2" class="flex flex-justify-end flex-items-center">
            <n-tooltip v-if="tab == '项目总览'">
              <template #trigger>
                <n-button size="small" type="primary" v-project-resource-code="'project:add'" @click="showProjectSettings">
                  <n-icon :component="SettingsServices" />
                </n-button>
              </template>
              项目设置
            </n-tooltip>
          </n-gi>
        </template>
        <n-gi v-else :span="18" class="flex flex-justify-center flex-items-center h-full">
          <n-text type="error" class="text-1.2em">项目不存在，请检查数据</n-text>
        </n-gi>
      </n-grid>
    </n-layout-header>
    <n-layout-content content-style="margin: 16px;">
      <template v-if="project?.id">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </template>
    </n-layout-content>
  </n-layout>

  <n-drawer v-model:show="showSettings" :width="401">
    <n-drawer-content>
      <n-form ref="formRef" :model="copiedProject" :rules="projectRules" label-width="100px" label-placement="left">
        <n-form-item label="项目名称" path="name">
          <n-input v-model:value="copiedProject!.name" placeholder="请输入项目名称" />
        </n-form-item>
        <n-form-item label="项目代码" path="code">
          <n-input v-model:value="copiedProject!.code" placeholder="请输入项目代码" />
        </n-form-item>
        <n-form-item label="项目描述" path="description">
          <n-input type="textarea" v-model:value="copiedProject!.description" placeholder="请输入项目描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" @click="resetUpdateProject">重置</n-button>
        <n-button type="primary" @click="handleCommitProject" v-project-resource-code="'project:update'">提交</n-button>
      </template>
      <template #header>
        <div class="w-350px flex flex-justify-between">
          <n-text class="mt-4px">项目设置</n-text>
          <n-popconfirm @positive-click="doDeleteProject" :show-icon="false">
            <template #trigger>
              <n-tooltip>
                <template #trigger>
                  <n-button type="error" size="tiny" v-project-resource-code="'project:delete'">
                    <n-icon :component="Delete" />
                  </n-button>
                </template>
                删除项目
              </n-tooltip>
            </template>
            <n-grid :cols="1" y-gap="16">
              <n-gi>请输入项目名称&lt;{{ project?.name }}&gt;以确认删除</n-gi>
              <n-gi><n-input size="small" v-model:value="deleteProjectName" /></n-gi>
            </n-grid>
          </n-popconfirm>
        </div>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style lang="scss" scoped>
#project-tabs {
  :deep(.n-tabs-wrapper) {
    margin-bottom: 2px;
  }
}
</style>
