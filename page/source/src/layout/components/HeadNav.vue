<template>
  <n-text tag="div" class="ui-logo">
    <svg-icon name="logo" :style="{height: '26px', width: '40px', color: logoColor}"/>
    <span style="margin-left: 10px;">项目管理</span>
  </n-text>
  <n-menu :value="activeKey" mode="horizontal" :options="menuOptions" :on-update:value="changeMainMenuKey"/>
  <div class="nav-end">
    <n-switch size="small" :value="themeMode" :on-update:value="changeTheme">
      <template #checked>
        <n-icon>
          <Sunny />
        </n-icon>
      </template>
      <template #unchecked>
        <n-icon>
          <Moon />
        </n-icon>
      </template>
    </n-switch>
    <span style="margin-left: 16px">  {{ welcome }}, {{ username }}!</span>
  </div>
</template>

<script setup lang="ts">
import { MenuOption } from "naive-ui"
import { ref, computed, onMounted } from "vue"
import { Sunny, Moon } from "@vicons/ionicons5"
import { useAppStore } from "../../store/app"
import {useUserStore} from "../../store/user";
import {useMemStore} from "../../store/mem";
const appStore = useAppStore()
const userStore = useUserStore()
const memStore = useMemStore()

const themeMode = computed(() => {
  return appStore.theme === "dark"
})
const logoColor = computed(() => {
  return appStore.theme === "dark" ? "#fff" : "#0582EE"
})
const welcome = ref<string>("欢迎")
const username = computed(() => userStore.user.username || "未登录")

const changeTheme = (value: boolean) => {
  appStore.setTheme(value ? "dark" : "light")
}

const activeKey = computed(() => appStore.activeMenuKey)
const menuOptions: MenuOption[] = memStore.mainMenus
const changeMainMenuKey = (key: string) => {
    appStore.activeMenuKey = key
}

onMounted(() => {
  const hour = new Date().getHours()
  if (hour >= 6 && hour < 12) {
    welcome.value = "早上好"
  } else if (hour >= 12 && hour < 14) {
    welcome.value = "中午好"
  } else if (hour >= 14 && hour < 18) {
    welcome.value = "下午好"
  } else {
    welcome.value = "晚上好"
  }
})
</script>

<style scoped>
.ui-logo {
  cursor: pointer;
  display: flex;
  font-size: 18px;
}
.logo {
    height: 26px;
    margin-right: 10px;
    fill: red !important;
}
</style>
