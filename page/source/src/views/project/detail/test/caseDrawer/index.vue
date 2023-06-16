<script setup lang="ts">
import { addProjectTestCase, batchSubmitProjectTestCase, updateProjectTestCase } from "@/service/api"
import { ProjectTestCase } from "@/types/project"
import { useMessage, FormInst } from "naive-ui"
import CaseLib from "../caseLib/index.vue"
const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  data: ProjectTestCase
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectTestCase) => {
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

const currentTab = ref<string>("directly")

const resetCurrentData = () => {
  currentData.value = props.data
}

// 批量选择的情况
const selectedTestCase = ref<ProjectTestCase[]>([])

const handleCommitData = () => {
  if (currentTab.value === "batchSelect") {
    if (selectedTestCase.value.length === 0) {
      message.warning("必须至少选择一个测试用例再进行提交")
      return
    }
    // 批量提交已选用例
    batchSubmitProjectTestCase(selectedTestCase.value).then((res) => {
      if (res) {
        message.success("批量提交成功")
        emit("refresh")
        emit("update:drawerActive", false)
      }
    })
  } else {
    formRef.value?.validate().then(() => {
      if (isUpdate.value) {
        updateProjectTestCase(currentData.value).then((res) => {
          if (res) {
            message.success("更新成功")
            emit("refresh")
            emit("update:drawerActive", false)
          }
        })
      } else {
        addProjectTestCase(currentData.value).then((res) => {
          if (res) {
            message.success("新增成功")
            emit("refresh")
            emit("update:drawerActive", false)
          }
        })
      }
    })
  }
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
      <n-tabs v-if="!isUpdate" type="card" animated v-model:value="currentTab">
        <n-tab-pane name="directly" tab="直接添加">
          <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
            <n-form-item label="用例名称" path="name">
              <n-input v-model:value="currentData.name" placeholder="请输入用例名称" />
            </n-form-item>
            <n-form-item label="用例描述" path="remark">
              <n-input v-model:value="currentData.remark" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入描述" />
            </n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="batchSelect" tab="用例库批量选择">
          <case-lib v-model:selected-test-case="selectedTestCase" />
        </n-tab-pane>
      </n-tabs>
      <n-form v-else ref="formRef" :model="currentData" label-width="100px" :rules="rules">
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
