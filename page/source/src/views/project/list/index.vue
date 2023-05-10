<script setup lang="ts">
import {ref, onMounted} from "vue";
import {getProjectList} from "../../../service/api/project";
import type {Project,ProjectCondition} from "../../../service/api/project";
import {Search} from "@vicons/carbon"
import {getStageList, type Stage} from "../../../service/api/stage";

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
                      <n-button type="primary" >新建项目</n-button>
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
                {{project.name}}
            </n-gi>
        </n-grid>
      </n-gi>
      <n-gi :span="2" class="bg-red h-200px">

      </n-gi>
  </n-grid>
</template>

<style scoped>

</style>