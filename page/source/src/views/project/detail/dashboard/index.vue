<script setup lang="ts">
import { Project, Stage, ProjectPlan } from "@/types/project"
import { Role } from "@/types/user"
import { getProjectMembers } from "@/service/api/project"
import { getProjectRiskCoefficient, ProjectRiskCoefficient } from "@/service/api/projectRisk"
import { getStageDetail } from "@/service/api/stage"
import { getExecutingProjectPlanByProjectId } from "@/service/api/projectPlan"
import { getRoleList } from "@/service/api/role"
import Members from "./members/index.vue"
import Invest from "./invest/index.vue"
import Requirement from "./requirement/index.vue"
import Task from "./task/index.vue"
import Plan from "./plan/index.vue"
import Risk from "./risk/index.vue"

const props = defineProps<{
  project: Project
}>()

// 健康度/风险分析内容 ///////////////////////
const riskCoefficient = ref<ProjectRiskCoefficient>({
  riskCoefficient: 0
})
const soh = computed(() => {
  const r = (1 - riskCoefficient.value.riskCoefficient) * 100 || 0
  return r <= 0 ? 0 : r
})
const riskType = computed(() => {
  if (soh.value > 90) {
    return "success"
  }
  if (soh.value > 60) {
    return "warning"
  }
  return "error"
})
// ///////////////////////////////////////////
// 项目计划信息
const projectPlan = ref<ProjectPlan>({
  id: "",
  projectId: "",
  planName: "",
  startTime: 0,
  endTime: 0,
  status: 0
})

// 阶段信息 //////////////////////////////////
const projectStage = ref<Stage>({
  id: "",
  name: "",
  orderNum: 0,
  status: 0
})

// 项目角色信息 //////////////////////////////
const projectRoles = ref<Array<Role>>([])

// 项目成员信息 //////////////////////////////
const refreshProjectMembers = () => {
  if (props.project) {
    getProjectMembers(props.project.id as string).then((res) => {
      props.project!.members = res
    })
  }
}

onMounted(() => {
  // 加载所有项目角色
  getRoleList({
    type: 2,
    dataPermission: 0,
    name: "",
    status: 1,
    offset: 0,
    limit: 100
  }).then((res) => {
    projectRoles.value = res.items
  })

  getProjectRiskCoefficient(props.project?.id as string).then((res) => {
    riskCoefficient.value = res
  })
  getExecutingProjectPlanByProjectId(props.project?.id as string).then((res) => {
    projectPlan.value = res
  })
  if (projectPlan.value?.stageId) {
    getStageDetail(projectPlan.value?.stageId as string).then((res) => {
      projectStage.value = res
    })
  }
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-grid :cols="3">
        <n-gi>
          <n-card embedded size="small" class="h-120px">
            <n-text tag="div" class="text-1.2em op-90 font-500 mb-10px">{{ project?.name }}</n-text>
            <n-ellipsis :line-clamp="2">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{{ project?.description }}</n-ellipsis>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card embedded size="small" class="h-120px">
            <n-text tag="div" class="text-1.2em op-90 font-500 mb-10px">健康度</n-text>
            <n-text :type="riskType" class="text-2.5em">{{ soh }} %</n-text>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card embedded size="small" class="h-120px">
            <n-text tag="div" class="text-1.2em op-90 font-500 mb-10px">当前阶段</n-text>
            {{ projectStage.name || "未建立项目计划" }}
          </n-card>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="2" x-gap="16">
        <n-gi>
          <n-card embedded size="small" class="h-160px">
            <members :members="project?.members || []" :roles="projectRoles" :project="project" @project-member-changed="refreshProjectMembers" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card embedded size="small" class="h-160px">
            <invest :members="project?.members || []" :project="project" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="2" x-gap="16">
        <n-gi>
          <n-card embedded size="small" class="h-180px">
            <plan :project="project" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card embedded size="small" class="h-180px">
            <requirement :project="project" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="2" x-gap="16">
        <n-gi>
          <n-card embedded size="small" class="h-180px">
            <task :project="project" />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card embedded size="small" class="h-180px">
            <risk :project="project" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-gi>
  </n-grid>
</template>

<style scoped></style>
