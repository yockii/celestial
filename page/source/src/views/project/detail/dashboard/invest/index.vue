<script setup lang="ts">
import {Project, ProjectMember, ProjectTaskWorkTimeStatistics} from "@/types/project";
import {getProjectWorkTimeStatistics} from "@/service/api/projectTask";

const props = defineProps<{
    project: Project,
    members: ProjectMember[],
}>()

const memberList = computed(() => {
    const map = new Map()
    const tmp = []
    for (let i = 0; i < props.members.length; i++) {
        if (!map.has(props.members[i].userId)) {
            map.set(props.members[i].userId, true)
            tmp.push(props.members[i])
        }
    }
    return tmp
})
const workTimeStatistics = ref<ProjectTaskWorkTimeStatistics>({})

onMounted(() => {
    // 获取统计信息
    getProjectWorkTimeStatistics(props.project.id).then(res => {
        workTimeStatistics.value = res
    })
})
</script>

<template>
    <n-text tag="div" class="text-1.2em op-90 font-500 mb-10px">投入概览</n-text>
    <n-grid :cols="3">
        <n-gi>
            <n-text tag="div" class="mb-4px">参与总人数</n-text>
            <n-text class="text-2em">{{ memberList.length }}</n-text>
        </n-gi>
        <n-gi>
            <n-text tag="div" class="mb-4px">预计总工时</n-text>
            <n-text class="text-2em">{{ workTimeStatistics.estimateDuration || 0 }}</n-text>
        </n-gi>
        <n-gi>
            <n-text tag="div" class="mb-4px">实际消耗工时</n-text>
            <n-text class="text-2em">{{ workTimeStatistics.actualDuration || 0 }}</n-text>
        </n-gi>
    </n-grid>
</template>

<style scoped>

</style>