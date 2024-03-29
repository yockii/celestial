<script setup lang="ts">
import { useRoute, useRouter } from "vue-router"
import { ref, onMounted } from "vue"
import { Project } from "@/types/project"
import { deleteProject, getProjectDetail, getProjectResourceCode, updateProject, getTopProjects } from "@/service/api"
import { KeyboardBackspaceOutlined } from "@vicons/material"
import { SettingsServices, Delete } from "@vicons/carbon"
import { FormInst, NButton } from "naive-ui"
import { useProjectStore } from "@/store/project"

const router = useRouter()
const route = useRoute()
const id = route.params.id as string
const tab = computed(() => route.meta.title as string)
const projectStore = useProjectStore()

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
  copiedProject.value = JSON.parse(JSON.stringify(projectStore.project))
}
const formRef = ref<FormInst>()
const handleCommitProject = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (copiedProject.value) {
        updateProject(copiedProject.value as Project).then((res) => {
          if (res) {
            projectStore.addProject(copiedProject.value)
            showSettings.value = false
          }
        })
      }
    }
  })
}

// 顶级项目筛选
const loadingTopProjects = ref<boolean>(false)
const topProjects = ref<Project[]>([])
const handleSearchTopProject = (query: string) => {
  loadingTopProjects.value = true
  getTopProjects(query)
    .then((res) => {
      if (res) {
        topProjects.value = res || []
      }
    })
    .finally(() => {
      loadingTopProjects.value = false
    })
}

// 删除项目 ////////////////////////////
const deleteProjectName = ref<string>("")
const doDeleteProject = () => {
  if (projectStore.project && projectStore.project?.name === deleteProjectName.value) {
    deleteProjectName.value = ""
    deleteProject(projectStore.project?.id as string).then((res) => {
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
  // project.value = { id: "", name: "", code: "", description: "", stageId: "" }
  // modules.value = []

  getProjectDetail(id).then((res) => {
    // project.value = res
    if (res) {
      projectStore.addProject(res)
    }
  })
  getProjectResourceCode(id).then((res) => {
    // resourceCodes.value = res || []
    if (res) {
      projectStore.setProjectResourceCodes(id, res)
    }
  })
})
</script>

<template>
  <n-layout class="ma--16px">
    <n-layout-header bordered class="h-48px py-4px px-20px">
      <n-grid :cols="24">
        <n-gi :span="2" class="flex flex-items-center">
          <router-link class="text-1.5em line-height-0.5em" :to="{ name: 'ProjectList' }">
            <n-icon>
              <KeyboardBackspaceOutlined />
            </n-icon>
          </router-link>
          <n-button text icon-placement="right" class="text-1.2em ml-16px">
            {{ projectStore.project?.name }}
          </n-button>
        </n-gi>
        <template v-if="projectStore.project?.id">
          <n-gi :span="16" :offset="2">
            <n-tabs id="project-tabs" :value="tab" type="line" justify-content="space-between" @update:value="handleChangeTab">
              <n-tab name="项目总览" v-if="projectStore.hasResourceCode('project:detail')"></n-tab>
              <n-tab name="项目计划" v-if="projectStore.hasResourceCode('project:detail:plan')"></n-tab>
              <n-tab name="功能模块" v-if="projectStore.hasResourceCode('project:detail:module')"></n-tab>
              <n-tab name="项目需求" v-if="projectStore.hasResourceCode('project:detail:requirement')"></n-tab>
              <n-tab name="工作任务" v-if="projectStore.hasResourceCode('project:detail:task')"></n-tab>
              <n-tab name="项目测试" v-if="projectStore.hasResourceCode('project:detail:test')"></n-tab>
              <n-tab name="项目缺陷" v-if="projectStore.hasResourceCode('project:detail:issue')"></n-tab>
              <n-tab name="项目变更" v-if="projectStore.hasResourceCode('project:detail:change')"></n-tab>
              <n-tab name="项目风险" v-if="projectStore.hasResourceCode('project:detail:risk')"></n-tab>
              <n-tab name="项目资产" v-if="projectStore.hasResourceCode('project:detail:asset')"></n-tab>
            </n-tabs>
          </n-gi>
          <n-gi :span="2" :offset="2" class="flex flex-justify-end flex-items-center">
            <n-tooltip v-if="tab == '项目总览'">
              <template #trigger>
                <div>
                  <n-button size="small" type="primary" v-if="projectStore.hasResourceCode('project:add')" @click="showProjectSettings">
                    <n-icon :component="SettingsServices" />
                  </n-button>
                </div>
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
      <template v-if="projectStore.project?.id">
        <router-view v-slot="{ Component, route }">
          <keep-alive>
            <component :is="Component" :key="route.fullPath" />
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
        <n-form-item label="主项目">
          <n-select
            v-model:value="copiedProject!.parentId"
            placeholder="请选择主项目"
            filterable
            clearable
            remote
            label-field="name"
            value-field="id"
            :loading="loadingTopProjects"
            :options="topProjects"
            @search="handleSearchTopProject"
          />
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
        <n-button type="primary" @click="handleCommitProject" v-if="projectStore.hasResourceCode('project:update')">提交</n-button>
      </template>
      <template #header>
        <div class="w-350px flex flex-justify-between">
          <n-text class="mt-4px">项目设置</n-text>
          <n-popconfirm @positive-click="doDeleteProject" :show-icon="false" v-if="projectStore.hasResourceCode('project:delete')">
            <template #trigger>
              <n-tooltip>
                <template #trigger>
                  <n-button type="error" size="tiny">
                    <n-icon :component="Delete" />
                  </n-button>
                </template>
                删除项目
              </n-tooltip>
            </template>
            <n-grid :cols="1" y-gap="16">
              <n-gi>请输入项目名称&lt;{{ projectStore.project?.name }}&gt;以确认删除</n-gi>
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
