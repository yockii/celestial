<script setup lang="ts">
import {useRoute} from 'vue-router'
import {ref, onMounted} from "vue";
import {getProjectDetail, type Project} from '@/service/api/project'
import {KeyboardBackspaceOutlined} from '@vicons/material'
import Dashboard from './dashboard/index.vue'
const id = useRoute().params.id as string
const project = ref<Project | null>(null)
const projectTab = ref<string>('项目总览')

onMounted(() => {
    getProjectDetail(id).then(res => {
        project.value = res
    })
})
</script>

<template>
    <n-layout class="ma--16px">
        <n-layout-header bordered class="h-48px py-4px px-20px">
            <n-grid :cols="24">
                <n-gi :span="2" class="flex flex-items-center">
                    <n-icon class="text-1.5em cursor-pointer">
                        <KeyboardBackspaceOutlined />
                    </n-icon>
                    <n-button text icon-placement="right" class="text-1.2em ml-16px">
                        {{project?.name }}
                    </n-button>
                </n-gi>
                <n-gi :span="12" :offset="4">
                    <n-tabs id="project-tabs" v-model:value="projectTab" type="line" justify-content="space-between">
                        <n-tab name="项目总览"></n-tab>
                        <n-tab name="工作任务"></n-tab>
                        <n-tab name="工作填报"></n-tab>
                        <n-tab name="项目文档"></n-tab>
                        <n-tab name="里程碑"></n-tab>
                        <n-tab name="商务总览"></n-tab>
                    </n-tabs>
                </n-gi>
                <n-gi :span="2" :offset="4" class="flex flex-justify-end flex-items-center">
                    <n-button type="primary">项目设置</n-button>
                </n-gi>
            </n-grid>
        </n-layout-header>
        <n-layout-content content-style="margin: 16px;">
            <dashboard v-if="projectTab == '项目总览'" :project="project" />
        </n-layout-content>
    </n-layout>
</template>

<style lang="scss" scoped>
#project-tabs {
    :deep(.n-tabs-wrapper) {
        margin-bottom: 2px;
    }
}
</style>