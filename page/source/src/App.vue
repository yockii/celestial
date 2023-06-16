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
import { useAppStore } from "./store/app"
import Watermarker from "./components/Watermark.vue"

const appStore = useAppStore()
const theme = computed(() => {
  return appStore.theme === "dark" ? darkTheme : null
})

// onBeforeUnmount(() => {
//     memStore.haltTicker()
// })
</script>

<style scoped></style>
