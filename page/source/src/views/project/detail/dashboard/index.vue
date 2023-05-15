<script setup lang="ts">
import {Project} from "@/service/api/project";

const props = defineProps<{
    project: Project | null
}>()
</script>

<template>
  <div>总览 {{project?.name}}</div>
</template>

<style scoped>

</style>