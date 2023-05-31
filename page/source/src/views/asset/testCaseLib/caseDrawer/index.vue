<script setup lang="ts">
import { CommonTestCase } from "@/types/asset"
import { useMessage, FormInst, FormItemRule } from "naive-ui"
import { addCommonTestCase, updateCommonTestCase } from "@/service/api/asset/commonTestCase"
const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  data: CommonTestCase
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: CommonTestCase) => {
    emit("update:data", val)
  }
})

const isUpdate = computed(() => {
  return currentData.value && currentData.value.id
})
const drawerTitle = computed(() => {
  return isUpdate.value ? "编辑用例" : "新增用例"
})
const formRef = ref<FormInst | undefined>()

const resetCurrentData = () => {
  currentData.value = props.data
}
const handleCommitData = () => {
  formRef.value?.validate().then(() => {
    if (isUpdate.value) {
      updateCommonTestCase(currentData.value).then((res) => {
        if (res) {
          message.success("更新成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    } else {
      addCommonTestCase(currentData.value).then((res) => {
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
    { required: true, message: "请输入用例名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ]
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show: boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="用例名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入用例名称" />
        </n-form-item>
        <n-form-item label="用例描述" path="remark">
          <n-input v-model:value="currentData.remark" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData">重置</n-button>
        <n-button size="small" type="primary" @click="handleCommitData">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
