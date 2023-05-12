<script setup lang="ts">
import {useRoute} from 'vue-router'
const id = useRoute().params.id
</script>

<template>
<div>
    {{id}}
</div>
</template>

<style scoped>

</style>