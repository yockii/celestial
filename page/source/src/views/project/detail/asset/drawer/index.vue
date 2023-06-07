<script setup lang="ts">
import { addProjectAsset, updateProjectAsset } from "@/service/api/project/projectAsset"
import { ProjectAsset } from "@/types/project"
import { FormInst, FormItemRule } from "naive-ui"
const message = useMessage()

const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])
const props = defineProps<{
  drawerActive: boolean
  data: ProjectAsset
}>()

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectAsset) => {
    emit("update:data", val)
  }
})
const isUpdate = computed(() => {
  return currentData.value && currentData.value.id
})
const drawerTitle = computed(() => {
  return isUpdate.value ? "编辑资产" : "新增资产"
})

const formRef = ref<FormInst | undefined>()

const resetCurrentData = () => {
  currentData.value = props.data
}
const handleCommitData = () => {
  formRef.value?.validate().then(() => {
    if (isUpdate.value) {
      updateProjectAsset(currentData.value).then((res) => {
        if (res) {
          message.success("更新成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    } else {
      addProjectAsset(currentData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    }
  })
}

// 规则定义
const rules = {
  name: [
    { required: true, message: "请输入资产名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: {
    validator(rule: FormItemRule, value: number) {
      return value > 0 && (value <= 4 || value === 9)
    },
    message: "请选择资产类型",
    trigger: "blur"
  }
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" @update:show="(show:boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="类型" path="type">
          <n-select
            v-model:value="currentData.type"
            placeholder="请选择类型"
            :options="[
              { label: '需求', value: 1 },
              { label: '设计', value: 2 },
              { label: '代码', value: 3 },
              { label: '测试', value: 4 },
              { label: '其他', value: 9 }
            ]"
          />
        </n-form-item>
        <n-form-item label="备注" path="remark">
          <n-input type="textarea" v-model:value="currentData.remark" placeholder="请输入备注" />
        </n-form-item>
        <n-form-item label="附件" path="attachment">
          <!-- <n-upload
            v-model:value="currentData.attachment"
            action="/api/asset/upload"
            multiple
            drag
            placeholder="点击或拖拽上传文件"
            :headers="{ Authorization: userStore.token }"
            :data="{ projectId: currentData.projectId }"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
          /> -->
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
