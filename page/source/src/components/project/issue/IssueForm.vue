<script setup lang="ts">
import { ProjectIssue } from "@/types/project"
import RichableTextArea from "@/components/RichableTextArea.vue"
import { FormInst } from "naive-ui"
defineProps<{
  value: ProjectIssue
}>()
const emit = defineEmits(["update:value"])

const rules = {
  title: [
    { required: true, message: "请输入缺陷名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: [{ type: "number", required: true, message: "请选择缺陷类型", trigger: ["blur", "change"] }]
}

const formRef = ref<FormInst>()
defineExpose({
  formRef: formRef
})
</script>

<template>
  <n-form ref="formRef" :model="value" :rules="rules" label-width="120px" label-placement="left">
    <n-grid :cols="4" x-gap="4">
      <n-gi>
        <n-form-item label="缺陷名称：" path="title">
          <n-input :value="value.title" @update:value="(v:string) => emit('update:value', Object.assign(value, {title: v}))" placeholder="请输入缺陷名称" />
        </n-form-item>
      </n-gi>
      <n-gi>
        <n-form-item label="类型：" path="type">
          <n-select
            :value="value.type"
            @update:value="(v:number) => emit('update:value', Object.assign(value, {type: v}))"
            placeholder="请选择类型"
            :options="[
              { label: '代码错误', value: 1 },
              { label: '功能异常', value: 2 },
              { label: '界面优化', value: 3 },
              { label: '配置相关', value: 4 },
              { label: '安全相关', value: 5 },
              { label: '性能相关', value: 6 },
              { label: '其他问题', value: 9 }
            ]"
          />
        </n-form-item>
      </n-gi>
      <n-gi :span="4">
        <n-form-item label="缺陷描述：">
          <!-- <n-input type="textarea" :autosize="{ minRows: 2, maxRows: 5 }" v-model:value="instance.content" placeholder="请输入缺陷描述" /> -->
          <richable-text-area
            :value="value.content"
            @update:value="(v:string) => emit('update:value', Object.assign(value, {content: v}))"
            label="缺陷描述"
            placeholder="请输入缺陷描述"
          />
        </n-form-item>
      </n-gi>
      <n-gi :span="4">
        <n-form-item label="原因：">
          <n-input
            type="textarea"
            :autosize="{ minRows: 2, maxRows: 5 }"
            :value="value.issueCause"
            @update:value="(v:string) => emit('update:value', Object.assign(value, {issueCause: v}))"
            placeholder="请输入问题原因，即为什么会出现这个缺陷"
          />
        </n-form-item>
      </n-gi>
      <n-gi :span="4">
        <n-form-item label="解决方法：">
          <n-input
            type="textarea"
            :autosize="{ minRows: 2, maxRows: 5 }"
            :value="value.solveMethod"
            @update:value="(v:string) => emit('update:value', Object.assign(value, {solveMethod: v}))"
            placeholder="请输入解决方法"
          />
        </n-form-item>
      </n-gi>
    </n-grid>
  </n-form>
</template>
