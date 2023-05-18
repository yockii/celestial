<template>
  <n-menu v-model:value="appStore.activeSubMenuKey" :collapsed="collapsed" :collapsed-width="60" :collapsed-icon-size="20" :options="menuOptions" > </n-menu>
</template>

<script setup lang="ts">
import {computed} from "vue"
import {useMemStore} from "@/store/mem";
import {useAppStore} from "@/store/app";

const memStore = useMemStore()
const appStore = useAppStore()
const menuOptions = computed(() =>  memStore.sideMenus)
const collapsed = computed(() => appStore.collapsed)

</script>
