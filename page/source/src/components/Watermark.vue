<script setup lang="ts">
import { computed } from "vue"
import { useUserStore } from "@/store/user"
import { storeToRefs } from "pinia"
const userStore = useUserStore()

const { realName } = storeToRefs(userStore)
const show = computed(() => realName.value && realName.value !== "")
</script>

<template>
  <n-watermark
    v-if="show"
    :content="realName"
    cross
    fullscreen
    :font-size="12"
    :line-height="12"
    :width="256"
    :height="256"
    :x-offset="12"
    :y-offset="60"
    :rotate="-17"
  />
</template>

<style scoped></style>
