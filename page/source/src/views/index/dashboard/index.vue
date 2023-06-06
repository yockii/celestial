<template>
  <div>abcdefghijklmn</div>
  <div>12345678901234</div>
</template>

<script setup lang="ts">
import { useAppStore } from "@/store/app"

const appStore = useAppStore()
onMounted(() => {
  appStore.activeMenuKey = "Home"
  appStore.activeSubMenuKey = "Dashboard"
})
</script>

<style scoped></style>
