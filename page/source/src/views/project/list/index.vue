<script setup lang="ts">
import {ref, onMounted, computed} from "vue";
import {
    getProjectList,
    addProject,
    getProjectStageStatistics,
} from "@/service/api/project";
import {Project,ProjectCondition,ProjectMember, ProjectStageStatistics} from "@/types/project";
import {Search} from "@vicons/carbon"
import {NButton, FormInst} from "naive-ui";
import dayjs from "dayjs";
import NameAvatar from "@/components/NameAvatar.vue";
import {useProjectStore} from "@/store/project";
import {storeToRefs} from "pinia";

const projectStore = useProjectStore()

const {stageList, stageListWithNone} = storeToRefs(projectStore)

const condition = ref<ProjectCondition>({
    offset: 0,
    limit: 10,
    name: "",
    stageId: "",
})
const total = ref(0)
const projectList = ref<Project[]>([]);

const refresh = () => {
    getProjectList(condition.value).then(res => {
        projectList.value = res.items;
    })
}
const getStageName = (stageId: string) => {
    return stageList.value.find(item => item.id === stageId)?.name || "无";
}
const projectStatistics = ref<ProjectStageStatistics[]>([])
const findStageProjectCount = (stageId: string) => {
    return projectStatistics.value.find(item => item.stageId === stageId)?.count || 0;
}
onMounted(() => {
    // 获取项目的阶段统计数据
    getProjectStageStatistics().then(res => {
        projectStatistics.value = res;
        total.value = res.reduce((total, item) => total + item.count, 0)
    })
    refresh()
})

// 切换stage
const selectedStageId = ref<string>('all')
const stageSelected = (stageId: string) => {
    selectedStageId.value = stageId
    condition.value.stageId = stageId === 'all' ? '' : stageId;
    refresh()
}

// 新建项目
const drawerActive = ref(false)
const newProject = ref<Project>({
    id: "",
    name: "",
    code: "",
    description: "",
    stageId: "",
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
const handleNewProject = () => {
    resetNewProject()
    drawerActive.value = true
}
const resetNewProject = () => {
    newProject.value = {
        id: "",
        name: "",
        code: "",
        description: "",
        stageId: stageList.value[0]?.id || "",
    }
}
const formRef = ref<FormInst>()
const handleCommitNewProject = () => {
    formRef.value?.validate(errors => {
        if (errors) {
            return
        }
        addProject(newProject.value).then((/* res */) => {
            drawerActive.value = false;
            refresh();
        })
    })
}
// 项目参与人
const projectMemberNames = (members: ProjectMember[] | undefined) : {src:string}[] => {
    if (members) {
        const map = new Map()
        const tmp = []
        for (let i = 0; i < members.length; i++) {
            if (!map.has(members[i].userId)) {
                map.set(members[i].userId, true)
                tmp.push({src: members[i].realName})
            }
        }
        return tmp
    }
    return   []
}
const createDropdownOptions = (options: Array<{ src:string }>) =>
    options.map((option) => ({
    key: option.src,
    label: option.src
}))

// 计算时间
const timeBefore = computed(() => (t) => dayjs(t).fromNow())
</script>

<template>
  <n-grid x-gap="12" :cols="7">
      <n-gi :span="5">
        <n-grid :cols="1" y-gap="12">
            <n-gi>
                <n-card embedded size="small">
                    <n-grid :cols="2">
                      <n-gi :span="1">
                          <n-input v-model:value="condition.name" placeholder="输入项目名称进行搜索" @keydown.enter.prevent="refresh">
                              <template #suffix>
                                  <n-icon class="cursor-pointer" @click="refresh">
                                      <Search />
                                  </n-icon>
                              </template>
                          </n-input>
                      </n-gi>
                      <n-gi :span="1" class="flex flex-justify-end">
                          <n-button type="primary" @click="handleNewProject">新建项目</n-button>
                      </n-gi>
                    </n-grid>
                    <n-grid :cols="1">
                  <n-gi style="margin-bottom: -12px;">
                      <n-tabs type="line" :on-update:value="stageSelected" :value="selectedStageId">
                          <n-tab name="all">全部</n-tab>
                          <n-tab v-for="stage in stageList" :name="stage.id" :key="stage.id">
                              {{stage.name}}
                          </n-tab>
                      </n-tabs>
                  </n-gi>
              </n-grid>
                </n-card>
            </n-gi>
            <n-gi class="" v-for="project in projectList" :key="project.id">
                <n-card embedded size="small">
                    <n-grid :cols="20" x-gap="8" y-gap="8">
                    <n-gi :span="12">
                        <router-link class="decoration-none" :to="{name: 'ProjectDetail', params: {id: project.id}}">
                            <n-text type="success" class="font-bold text-lg">
                            {{project.name}}
                            </n-text>
                        </router-link>
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
                            创建时间: {{timeBefore(project.createTime)}}
                        </n-text>
                    </n-gi>
                    <n-gi :span="4">
                        <n-text depth="3">
                            当前阶段: {{getStageName(project.stageId)}}
                        </n-text>
                    </n-gi>
                    <n-gi :span="6" :offset="6" class="flex flex-justify-end">
                        <n-avatar-group :options="projectMemberNames(project.members)" :size="24" :max="5">
                            <template #avatar="{ option: {src} }">
                                <n-tooltip>
                                    <template #trigger>
                                        <name-avatar :name="src"  />
                                    </template>
                                    {{ src }}
                                </n-tooltip>
                            </template>
                            <template #rest="{ options: restOptions, rest }">
                                <n-dropdown :options="createDropdownOptions(restOptions)" placement="top">
                                    <n-avatar>+{{ rest }}</n-avatar>
                                </n-dropdown>
                            </template>
                        </n-avatar-group>
                    </n-gi>
                </n-grid>
                </n-card>
            </n-gi>
        </n-grid>
      </n-gi>
      <n-gi :span="2">
        <n-grid :cols="1" y-gap="8">
            <n-gi class="">
                <n-card embedded size="small">
                    <n-text class="font-700">项目阶段统计</n-text>
                    <n-grid :cols="3" class="mt-8px">
                        <n-gi @click="stageSelected('all')">
                            <n-text class="cursor-pointer list-item ml-20px text-1em text-gray" :type="selectedStageId == 'all' ? 'success' : 'default'">所有项目</n-text>
                            <n-text tag="div" class="cursor-pointer font-500 text-2.5em w-full pl-20px" :type="selectedStageId == 'all' ? 'success' : 'default'">{{total}}</n-text>
                        </n-gi>
                        <n-gi v-for="stage in stageList" :key="stage.id" @click="stageSelected(stage.id)">
                            <n-text class="cursor-pointer list-item ml-20px text-1em text-gray" :type="selectedStageId == stage.id ? 'success' : 'default'">{{ stage.name }}</n-text>
                            <n-text tag="div" class="cursor-pointer font-500 text-2.5em w-full pl-20px" :type="selectedStageId == stage.id ? 'success' : 'default'">
                                {{ findStageProjectCount(stage.id) }}
                            </n-text>
                        </n-gi>
                    </n-grid>
                </n-card>
            </n-gi>
        </n-grid>
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
              <n-form-item label="项目阶段" path="stageId">
                  <n-select
                      v-model:value="newProject.stageId"
                      placeholder="请选择项目阶段"
                      :options="stageListWithNone"
                      label-field="name"
                      value-field="id"
                  />
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