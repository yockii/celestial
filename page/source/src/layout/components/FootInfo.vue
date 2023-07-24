<template>
  <div class="h-12px line-height-12px b-t">
    <n-space class="text-1" justify="space-between">
      <div class="scale-60" @click="readyForEruda">版本: v0.2.0.0724</div>
      <div class="scale-60">© 2023-2028</div>
      <div class="scale-60">请求耗时: {{ lastRequestTime }}ms</div>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from "@/store/app"
import { storeToRefs } from "pinia"

const { lastRequestTime } = storeToRefs(useAppStore())

const initing = ref(false)
const clickCount = ref(0)
const readyForEruda = () => {
  if (initing.value) {
    return
  }
  if (clickCount.value >= 10) {
    initing.value = true
    const s = document.createElement("script")
    s.type = "text/javascript"
    s.src = "https://cdn.jsdelivr.net/npm/eruda"
    document.body.appendChild(s)
    s.onload = () => {
      const sInit = document.createElement("script")
      sInit.type = "text/javascript"
      sInit.innerHTML = "eruda.init();"
      document.body.appendChild(sInit)
      initing.value = false
    }
    clickCount.value = 0
  } else {
    clickCount.value++
  }
}
</script>
