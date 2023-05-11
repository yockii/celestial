<script setup lang="ts">
import {ref, onMounted} from "vue";
import {getProjectList, addProject} from "../../../service/api/project";
import type {Project,ProjectCondition} from "../../../service/api/project";
import {Search} from "@vicons/carbon"
import {getStageList, type Stage} from "../../../service/api/stage";
import {NButton} from "naive-ui";
import moment from "moment";

const condition = ref<ProjectCondition>({
    offset: 0,
    limit: 10,
    name: "",
    stageId: "",
})
const stageList = ref<Stage[]>([])
const total = ref(0)
const projectList = ref<Project[]>([]);

const refresh = () => {
    getProjectList(condition.value).then(res => {
        projectList.value = res.items;
        total.value = res.total;
    })
}

onMounted(() => {
    getStageList().then(res => {
        stageList.value = res.items;
    })
    refresh()
})

// 新建项目
const drawerActive = ref(false)
const newProject = ref<Project>({
    id: "",
    name: "",
    code: "",
    description: "",
    stageId: "",
    createTime: 0,
})
const projectRules = {
    name: [
        { required: true, message: '请输入项目名称', trigger: 'blur' },
        { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
    ],
    code: [
        { required: true, message: '请输入项目代码', trigger: 'blur' },
        { pattern: /^\D\w{2,19}$/, message: '长度在 3 到 20 个英文/数字/下划线，且不能以数字开头', trigger: 'blur' }
    ]
}
const resetNewProject = () => {
    newProject.value = {
        id: "",
        name: "",
        code: "",
        description: "",
        stageId: "",
        createTime: 0,
    }
}
const handleCommitNewProject = () => {
    addProject(newProject.value).then(res => {
        drawerActive.value = false;
        refresh();
    })
}
</script>

<template>
  <n-grid x-gap="12" :cols="7">
      <n-gi :span="5">
        <n-grid :cols="1" y-gap="12">
            <n-gi class="bg-gray-100 pa-16px !pb-0">
              <n-grid :cols="2">
                  <n-gi :span="1">
                      <n-input v-model="condition.name" placeholder="输入项目名称进行搜索" @keydown.enter.prevent="refresh">
                          <template #suffix>
                              <n-icon :component="Search" class="cursor-pointer" @click="refresh"></n-icon>
                          </template>
                      </n-input>
                  </n-gi>
                  <n-gi :span="1" class="flex flex-justify-end">
                      <n-button type="primary" @click="drawerActive = true">新建项目</n-button>
                  </n-gi>
              </n-grid>
              <n-grid :cols="1">
                  <n-gi>
                      <n-tabs type="line">
                          <n-tab name="all">全部</n-tab>
                          <n-tab name="stage.id" v-for="stage in stageList" :key="stage.id">
                              {{stage.name}}
                          </n-tab>
                      </n-tabs>
                  </n-gi>
              </n-grid>
            </n-gi>
            <n-gi class="bg-gray-100 pa-8px" v-for="project in projectList" :key="project.id">
                <n-grid :cols="20" x-gap="8" y-gap="8">
                    <n-gi :span="12" class="font-bold text-lg">
                        {{project.name}}
                    </n-gi>
                    <n-gi :span="8" class="flex flex-justify-end">
                        <n-text depth="3">
                            已用工时: 0小时
                        </n-text>
                    </n-gi>
                    <n-gi :span="20">
                        <n-ellipsis>{{project.description}}</n-ellipsis>
                    </n-gi>
                    <n-gi :span="4">
                        <n-text depth="3">
                            创建时间: {{moment(project.createTime).fromNow()}}
                        </n-text>
                    </n-gi>
                    <n-gi :span="4">
                        <n-text depth="3">
                            当前阶段: {{project.stageId}}
                        </n-text>
                    </n-gi>
                </n-grid>
            </n-gi>
        </n-grid>
      </n-gi>
      <n-gi :span="2" class="bg-red h-200px">

      </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
      <n-drawer-content title="新建项目" closable>
          <n-form ref="formRef" :model="newProject" :rules="projectRules" label-width="100px" label-placement="left">
              <n-form-item label="项目名称" path="name">
                  <n-input v-model:value="newProject.name" placeholder="请输入项目名称" />
              </n-form-item>
              <n-form-item label="项目代码" path="code">
                  <n-input v-model:value="newProject.code" placeholder="请输入项目代码" />
              </n-form-item>
              <n-form-item label="项目描述" path="description">
                  <n-input type="textarea" v-model:value="newProject.description" placeholder="请输入项目描述" />
              </n-form-item>
          </n-form>
          <template #footer>
              <n-button class="mr-a" @click="resetNewProject">重置</n-button>
              <n-button type="primary" @click="handleCommitNewProject">提交</n-button>
          </template>
      </n-drawer-content>
  </n-drawer>
</template>

<style scoped>

</style>