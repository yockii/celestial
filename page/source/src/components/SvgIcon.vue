<script setup lang="ts">
const props = defineProps({
  name: {
    type: String,
    required: true
  },
  color: {
    type: String,
    default: ""
  },
  prefix: {
    type: String,
    default: "icon"
  }
})
const symbolId = computed(() => `#${props.prefix}-${props.name}`)
</script>

<template>
  <svg aria-hidden="true" v-bind="$attrs" class="svg-icon">
    <use :xlink:href="symbolId" :fill="props.color"></use>
  </svg>
</template>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  fill: currentColor;
  vertical-align: middle;
}
</style>
