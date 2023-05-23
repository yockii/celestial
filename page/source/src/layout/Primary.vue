<template>
  <n-layout position="absolute">
    <n-layout-header bordered style="height: 48px" class="nav">
      <HeadNav />
    </n-layout-header>
    <n-layout :has-sider="needSidebar" position="absolute" style="top: 48px; /*bottom: 24px*/">
      <n-layout-sider
        v-if="needSidebar"
        bordered
        :collapsed-width="60"
        :width="180"
        collapse-mode="width"
        :collapsed="collapsed"
        show-trigger
        @collapse="appStore.collapsed = true"
        @expand="appStore.collapsed = false"
      >
        <Sider />
      </n-layout-sider>
      <n-layout-content content-style="padding: 16px;">
        <RouterView />
      </n-layout-content>
    </n-layout>
    <!--    <n-layout-footer position="absolute" style="height: 24px">-->
    <!--      <FootInfo />-->
    <!--    </n-layout-footer>-->
  </n-layout>
</template>

<script setup lang="ts">
import HeadNav from "./components/HeadNav.vue"
import Sider from "./components/Sider.vue"
// import FootInfo from "./components/FootInfo.vue"
import { useMemStore } from "@/store/mem"
import { computed } from "vue"
import { useAppStore } from "@/store/app"
import { getStageList } from "@/service/api/stage"
import { useStageStore } from "@/store/stage"

const memStore = useMemStore()
const appStore = useAppStore()
const collapsed = computed(() => appStore.collapsed)

const needSidebar = computed(() => memStore.sideMenus.length > 0)

const stageStore = useStageStore()
onMounted(() => {
  // memStore.startTicker()
  getStageList({ limit: 100, name: "", offset: 0, status: 0 }).then((res) => {
    if (res) {
      stageStore.stageList = res.items
    }
  })
})
</script>

<style scoped>
.nav {
  align-items: center;
  display: flex;
  justify-content: space-between;
  padding: 4px 20px;
}
</style>
