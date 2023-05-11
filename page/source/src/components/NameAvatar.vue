<script setup lang="ts">
import {defineProps, computed} from 'vue'
const props = defineProps<{
    name: string;
}>()
const color = computed(() => {
    if (props.name) {
        const tmp: string[] = []
        tmp.push("#")
        for (let i = 0; i < props.name.length; i++) {
            tmp.push(parseInt(String(props.name[i].charCodeAt(0)), 10).toString(16))
        }
        return tmp.slice(0, 5).join('').slice(0, 4);
    }
    return "#000000"
})
const frontColor = computed(() => {
    // 计算灰度
    const gray = parseInt(color.value.slice(1, 3), 16) * 0.299 + parseInt(color.value.slice(3, 5), 16) * 0.587 + parseInt(color.value.slice(5, 7), 16) * 0.114
    // 根据灰度，给出前景色
    return gray > 192 ? "#000000" : "#ffffff"
})
</script>

<template>
  <n-avatar round :color="color" :style="{color: frontColor}">
      {{name[0]}}
  </n-avatar>
</template>

<style scoped>

</style>