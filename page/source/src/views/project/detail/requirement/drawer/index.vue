<script setup lang="ts">
import { ProjectRequirement } from "@/types/project"
import { useProjectStore } from "@/store/project"
import { storeToRefs } from "pinia"
import { useMessage, FormInst, FormItemRule } from "naive-ui"
import { addProjectRequirement, updateProjectRequirement } from "@/service/api/project/projectRequirement"
const message = useMessage()
const projectStore = useProjectStore()
const { moduleTree } = storeToRefs(projectStore)
const props = defineProps<{
  projectId: string
  drawerActive: boolean
  data: ProjectRequirement
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectRequirement) => {
    emit("update:data", val)
  }
})

const isUpdate = computed(() => {
  return currentData.value && currentData.value.id
})
const drawerTitle = computed(() => {
  return isUpdate.value ? "编辑需求" : "新增需求"
})
const formRef = ref<FormInst | undefined>()

const resetCurrentData = () => {
  currentData.value = props.data
}
const handleCommitData = () => {
  formRef.value?.validate().then(() => {
    if (isUpdate.value) {
      updateProjectRequirement(currentData.value).then((res) => {
        if (res) {
          message.success("更新成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    } else {
      addProjectRequirement(currentData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    }
  })
}

onMounted(() => {
  resetCurrentData()
})

// 规则定义
const rules = {
  name: [
    { required: true, message: "请输入需求名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: {
    validator(rule: FormItemRule, value: number) {
      return value > 0 && value <= 7
    },
    message: "请选择需求类型",
    trigger: "blur"
  },
  priority: {
    validator(rule: FormItemRule, value: number) {
      return value > 0 && value <= 3
    },
    message: "请选择需求优先级",
    trigger: "blur"
  }
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show: boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="需求名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入需求名称" />
        </n-form-item>
        <n-form-item label="功能模块">
          <n-cascader
            v-model:value="currentData.moduleId"
            :options="moduleTree"
            value-field="id"
            label-field="name"
            children-field="children"
            placeholder="请选择功能模块"
            clearable
            show-path
          />
        </n-form-item>
        <n-form-item label="需求类型" path="type">
          <n-select
            v-model:value="currentData.type"
            placeholder="请选择需求类型"
            :options="[
              { label: '功能', value: 1 },
              { label: '接口', value: 2 },
              { label: '性能', value: 3 },
              { label: '安全', value: 4 },
              { label: '体验', value: 5 },
              { label: '改进', value: 6 },
              { label: '其他', value: 7 }
            ]"
          />
        </n-form-item>
        <n-form-item label="优先级" path="priority">
          <n-select
            v-model:value="currentData.priority"
            placeholder="请选择优先级"
            :options="[
              { label: '低', value: 1 },
              { label: '中', value: 2 },
              { label: '高', value: 3 }
            ]"
          />
        </n-form-item>
        <n-form-item label="来源" path="source">
          <n-select
            v-model:value="currentData.source"
            placeholder="请选择来源"
            :options="[
              { label: '客户', value: 1 },
              { label: '内部', value: 2 }
            ]"
          />
        </n-form-item>
        <n-form-item label="详情描述" path="detail">
          <n-input v-model:value="currentData.detail" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData">重置</n-button>
        <n-button
          size="small"
          type="primary"
          @click="handleCommitData"
          v-resource-code="['project:detail:requirement:add', 'project:detail:requirement:update']"
          >提交</n-button
        >
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
