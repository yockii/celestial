<script setup lang="ts">

import {Project} from "@/service/api/project";
import {getProjectRiskCoefficient, ProjectRiskCoefficient} from "@/service/api/projectRisk";
import {getStageDetail, Stage} from "@/service/api/stage";
import {getExecutingProjectPlanByProjectId, ProjectPlan} from "@/service/api/projectPlan";
import {getRoleList, Role} from "@/service/api/role";

const props = defineProps<{
    project: Project | null
}>()

// 健康度/风险分析内容 ///////////////////////
const riskCoefficient = ref<ProjectRiskCoefficient>({
    riskCoefficient: 0,
    riskCoefficientDetail: null
})
const soh = computed(() => {
    const r = (1 - riskCoefficient.value.riskCoefficient) * 100 || 0
    return r <= 0 ? 0 : r
})
const riskType = computed(() => {
    if (soh.value > 90) {
        return 'success'
    }
    if (soh.value > 60) {
        return 'warning'
    }
    return 'error'
})
// ///////////////////////////////////////////
// 项目计划信息
const projectPlan = ref<ProjectPlan>({})

// 阶段信息 //////////////////////////////////
const projectStage = ref<Stage>({})

// 项目角色信息 //////////////////////////////
const projectRoles = ref<Role[]>([])

onMounted(() => {
    // 加载所有项目角色
    getRoleList({
        type: 2, dataPermission: 0, name: "",
        status: 1,
        offset: 0,
        limit: 100
    }).then(res => {
        projectRoles.value = res.items
    })

    getProjectRiskCoefficient(props.project?.id as string).then(res => {
        riskCoefficient.value = res
    })
    getExecutingProjectPlanByProjectId(props.project?.id as string).then(res => {
        projectPlan.value = res
    })
    if (projectPlan.value?.stageId) {
        getStageDetail(projectPlan.value?.stageId as string).then(res => {
            projectStage.value = res
        })
    }

})
</script>

<template>
    <n-grid :cols="1" y-gap="16">
        <n-gi>
          <n-grid :cols="3" x-gap="1">
              <n-gi>
                  <n-card embedded size="small" :title="project?.name" class="h-120px">
                        <n-ellipsis :line-clamp="2">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{{project?.description}}</n-ellipsis>
                  </n-card>
              </n-gi>
              <n-gi>
                  <n-card embedded size="small" title="健康度" class="h-120px">
                      <n-text :type="riskType" class="text-2.5em">{{soh}} %</n-text>
                  </n-card>
              </n-gi>
              <n-gi>
                  <n-card embedded size="small" title="当前阶段" class="h-120px">
                    {{projectStage.name || "未建立项目计划"}}
                  </n-card>
              </n-gi>
          </n-grid>
        </n-gi>
        <n-gi>
            <n-grid :cols="2" x-gap="16">
                <n-gi>
                    <n-card embedded size="small" title="项目组成员" class="h-120px">

                    </n-card>
                </n-gi>
                <n-gi>
                    <n-card embedded size="small" title="工作填报" class="h-120px">

                    </n-card>
                </n-gi>
            </n-grid>
        </n-gi>
    </n-grid>
</template>

<style scoped>

</style>