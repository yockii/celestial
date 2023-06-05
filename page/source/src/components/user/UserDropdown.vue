<script setup lang="ts">
import { useUserStore } from "@/store/user"

defineProps<{
  name: string
}>()
const options = ref([
  {
    label: "退出登录",
    key: "logout"
  }
])
const router = useRouter()
const handleSelect = (key: string) => {
  if (key === "logout") {
    useUserStore().logout()
    router.push("/login")
  }
}
</script>

<template>
  <n-dropdown trigger="click" :options="options" @select="handleSelect">
    <n-text>{{ name }}</n-text>
  </n-dropdown>
</template>
