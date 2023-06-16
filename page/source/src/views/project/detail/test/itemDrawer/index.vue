<script setup lang="ts">
import { ProjectTestCaseItem } from "@/types/project"
import { useMessage, FormInst } from "naive-ui"
import { addProjectTestCaseItem, batchSubmitProjectTestCaseItem, updateProjectTestCaseItem } from "@/service/api"
import CaseLib from "../caseLib/index.vue"

const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  data: ProjectTestCaseItem
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectTestCaseItem) => {
    emit("update:data", val)
  }
})

const isUpdate = computed(() => {
  return currentData.value && currentData.value.id
})
const drawerTitle = computed(() => {
  return isUpdate.value ? "编辑用例项" : "新增用例项"
})
const formRef = ref<FormInst | undefined>()

const currentTab = ref<string>("directly")

const resetCurrentData = () => {
  currentData.value = props.data
}

// 批量选择的情况
const selectedTestCaseItem = ref<ProjectTestCaseItem[]>([])

const handleCommitData = () => {
  if (currentTab.value === "batchSelect") {
    if (selectedTestCaseItem.value.length === 0) {
      message.warning("必须至少选择一个测试用例项再进行提交")
      return
    }

    // 新建数组，不改变原来的数据，为每个用例项添加 projectId 和 testCaseId
    const newSelectedTestCaseItem = selectedTestCaseItem.value.map((item) => {
      const newItem = { ...item }
      newItem.projectId = props.data.projectId
      newItem.testCaseId = props.data.testCaseId
      return newItem
    })

    // 批量提交已选择用例项
    batchSubmitProjectTestCaseItem(newSelectedTestCaseItem).then((res) => {
      if (res) {
        message.success("批量提交成功")
        emit("refresh")
        emit("update:drawerActive", false)
      }
    })
  } else {
    formRef.value?.validate().then(() => {
      if (isUpdate.value) {
        updateProjectTestCaseItem(currentData.value).then((res) => {
          if (res) {
            message.success("更新成功")
            emit("refresh")
            emit("update:drawerActive", false)
          }
        })
      } else {
        addProjectTestCaseItem(currentData.value).then((res) => {
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
    { required: true, message: "请输入用例项内容", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 30 个字符", trigger: "blur" }
  ]
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show: boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-tabs v-if="!isUpdate" type="card" animated v-model:value="currentTab">
        <n-tab-pane name="directly" tab="直接添加">
          <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
            <n-form-item label="名称" path="name">
              <n-input v-model:value="currentData.name" placeholder="请输入名称" />
            </n-form-item>
            <n-form-item label="类型" path="type">
              <n-select
                v-model:value="currentData.type"
                :options="[
                  { label: '功能测试', value: 1 },
                  { label: '性能测试', value: 2 },
                  { label: '安全测试', value: 3 },
                  { label: '兼容性测试', value: 4 },
                  { label: '接口测试', value: 5 },
                  { label: '压力测试', value: 6 },
                  { label: '其他', value: 9 }
                ]"
                placeholder="请选择类型"
              />
            </n-form-item>
            <n-form-item label="内容" path="content">
              <n-input v-model:value="currentData.content" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入内容" />
            </n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="batchSelect" tab="用例库批量选择">
          <case-lib v-model:selected-test-case-item="selectedTestCaseItem" />
        </n-tab-pane>
      </n-tabs>
      <n-form v-else ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="类型" path="type">
          <n-select
            v-model:value="currentData.type"
            :options="[
              { label: '功能测试', value: 1 },
              { label: '性能测试', value: 2 },
              { label: '安全测试', value: 3 },
              { label: '兼容性测试', value: 4 },
              { label: '接口测试', value: 5 },
              { label: '压力测试', value: 6 },
              { label: '其他', value: 9 }
            ]"
            placeholder="请选择类型"
          />
        </n-form-item>
        <n-form-item label="内容" path="content">
          <n-input v-model:value="currentData.content" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入内容" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData">重置</n-button>
        <n-button size="small" type="primary" @click="handleCommitData">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
