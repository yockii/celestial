<template>
  <n-text tag="div" class="ui-logo">
    <svg-icon name="logo" :style="{ height: '26px', width: '40px', color: logoColor }" />
    <span style="margin-left: 10px"></span>
  </n-text>
  <n-menu v-model:value="appStore.activeMenuKey" mode="horizontal" :options="menuOptions" />
  <div class="nav-end">
    <n-space>
      <n-dropdown trigger="hover" :options="historyList" key-field="url" show-arrow @select="gotoHistory" :render-option="renderHistory">
        <n-text class="cursor-pointer">
          <n-icon>
            <Footsteps />
          </n-icon>
        </n-text>
      </n-dropdown>
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
      <span> {{ welcome }}, <user-dropdown :name="realName || '未登录'" />!</span>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { Sunny, Moon, Footsteps } from "@vicons/ionicons5"
import { useAppStore } from "@/store/app"
import { useUserStore } from "@/store/user"
import { useMemStore } from "@/store/mem"
import { storeToRefs } from "pinia"
import UserDropdown from "@/components/user/UserDropdown.vue"
import { RouteHistory } from "@/types/app"
import { NTooltip } from "naive-ui"
import dayjs from "dayjs"
const appStore = useAppStore()
const userStore = useUserStore()
const memStore = useMemStore()
const router = useRouter()

const themeMode = computed(() => {
  return appStore.theme === "dark"
})
const logoColor = computed(() => {
  return appStore.theme === "dark" ? "#fff" : "#0582EE"
})
const welcome = ref<string>("欢迎")
const { realName, history } = storeToRefs(userStore)

const changeTheme = (value: boolean) => {
  appStore.setTheme(value ? "dark" : "light")
}

const { mainMenus: menuOptions } = storeToRefs(memStore)

const historyList = computed(() => [...history.value].reverse())
const renderHistory = ({ node, option }: { node: VNode; option: RouteHistory }) => {
  return h(
    NTooltip,
    { keepAliveOnHover: false, style: { width: "max-content" } },
    {
      trigger: () => [node],
      default: () => dayjs(option.time).fromNow() + " " + option.url
    }
  )
}
const gotoHistory = (key: string | number) => {
  router.push(key as string)
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
