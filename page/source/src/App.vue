<template>
  <n-config-provider :theme="theme" :locale="zhCN" :date-locale="dateZhCN">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <router-view />
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
  <watermarker />
</template>

<script setup lang="ts">
import { darkTheme, zhCN, dateZhCN } from "naive-ui"
import { computed  /* , onMounted, onBeforeUnmount */ } from "vue"
import { useAppStore } from "./store/app"
import Watermarker from "./components/Watermark.vue"
import {getStageList } from "@/service/api/stage";
import {useProjectStore} from "@/store/project";
const projectStore = useProjectStore()

const appStore = useAppStore()
const theme = computed(() => {
  return appStore.theme === "dark" ? darkTheme : null
})

onMounted(() => {
    // memStore.startTicker()
  getStageList({limit: 100, name: "", offset: 0, status: 0}).then(res => {
        if (res) {
            projectStore.stageList = res.items
        }
    })
})
// onBeforeUnmount(() => {
//     memStore.haltTicker()
// })
</script>

<style scoped></style>
