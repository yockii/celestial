<script setup lang="ts">
import { ProjectRequirement, ProjectTask } from "@/types/project"
import { useProjectStore } from "@/store/project"
import { storeToRefs } from "pinia"
import { useMessage, FormInst } from "naive-ui"
import { addProjectTask, getProjectTaskList, updateProjectTask, getProjectRequirementList } from "@/service/api"
const message = useMessage()
const projectStore = useProjectStore()
const { project, modules, moduleTree, memberList } = storeToRefs(projectStore)
const props = defineProps<{
  drawerActive: boolean
  data: ProjectTask
}>()
const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: ProjectTask) => {
    emit("update:data", val)
  }
})

const isUpdate = computed(() => {
  return currentData.value && currentData.value.id
})
const drawerTitle = computed(() => {
  return isUpdate.value ? "编辑任务" : "新增任务"
})
const formRef = ref<FormInst | undefined>()

const resetCurrentData = () => {
  currentData.value = props.data
}
const handleCommitData = () => {
  formRef.value?.validate().then(() => {
    currentData.value.members = memberList.value
      .filter((item) => memberIdList.value.includes(item.userId))
      .map((item) => {
        return {
          id: "",
          projectId: project.value.id,
          userId: item.userId,
          taskId: currentData.value.id
        }
      })

    if (isUpdate.value) {
      updateProjectTask(currentData.value).then((res) => {
        if (res) {
          message.success("更新成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    } else {
      addProjectTask(currentData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    }
  })
}

// 根据模块选择需求
const requirementList = ref<ProjectRequirement[]>([])
const loadRequirementList = (moduleId: string) => {
  if (moduleId) {
    currentData.value.moduleId = moduleId
    const module = modules.value.find((item) => item.id === moduleId)
    if (module) {
      // 从接口获取需求列表
      getProjectRequirementList({
        projectId: project.value.id,
        fullPath: module.fullPath,
        offset: -1,
        limit: -1
      }).then((res) => {
        if (res) {
          requirementList.value = res.items || []
        }
      })
    }
  } else {
    requirementList.value = []
    currentData.value.moduleId = ""
  }
}

onBeforeUpdate(() => {
  if (props.drawerActive) {
    loadRequirementList(currentData.value.moduleId || "")
    loadParentTask("", currentData.value.parentId || "")
    memberIdList.value = currentData.value.members?.map((item) => item.userId) || []
  }
})

// 选择需求后自动填充任务详情
const updateTaskDescByRequirement = (requirementId: string) => {
  currentData.value.requirementId = requirementId
  if (currentData.value.taskDesc) return
  const requirement = requirementList.value.find((item) => item.id === requirementId)
  if (requirement && requirement.detail) {
    currentData.value.taskDesc = requirement.detail
  }
}

onMounted(() => {
  resetCurrentData()
})

// 处理任务参与人
const memberIdList = ref<string[]>([])

const taskTime = computed<number[] | null>({
  get: () => {
    const startTime = currentData.value.startTime
    const endTime = currentData.value.endTime
    if (startTime && endTime) {
      return [startTime, endTime]
    }
    return null
    // return [currentData.value.startTime, currentData.value.endTime + 3600 * 1000 * 24]
  },
  set: (val: number[] | null) => {
    if (val) {
      if (val.length === 1) {
        currentData.value.startTime = val[0]
      } else if (val.length === 2) {
        currentData.value.startTime = val[0]
        currentData.value.endTime = val[1]
      }
    }
  }
})

// 规则定义
const rules = {
  name: [
    { required: true, message: "请输入任务名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  moduleId: [{ required: true, message: "请选择功能模块", trigger: "blur" }],
  requirementId: [{ required: true, message: "请选择所属需求", trigger: "blur" }],
  members: [
    {
      type: "array",
      required: true,
      validator: () => {
        return memberIdList.value.length > 0
      },
      message: "请选择任务参与人",
      trigger: ["blur", "change"]
    }
  ]
}

// 父级任务异步加载
const loadingParentTask = ref(false)
const parentTaskList = ref<ProjectTask[]>([])
const loadParentTask = (keyword = "", id = "") => {
  loadingParentTask.value = true
  const params = {
    id,
    projectId: project.value.id,
    offset: 0,
    limit: 10,
    name: keyword
  }
  getProjectTaskList(params)
    .then((res) => {
      if (res) {
        parentTaskList.value = res.items || []
      }
    })
    .finally(() => {
      loadingParentTask.value = false
    })
}
const parentTaskChanged = (parentId: string) => {
  currentData.value.parentId = parentId
  // 同时修改功能模块和需求为父级任务的
  const parentTask = parentTaskList.value.find((item) => item.id === parentId)
  if (parentTask) {
    currentData.value.moduleId = parentTask.moduleId
    currentData.value.requirementId = parentTask.requirementId
  }
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="460" :on-update:show="(show: boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="任务名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入任务名称" />
        </n-form-item>
        <n-form-item label="父级任务" path="parentId">
          <n-select
            v-model:value="currentData.parentId"
            placeholder="请搜索并选择父级任务"
            :loading="loadingParentTask"
            :options="parentTaskList"
            label-field="name"
            value-field="id"
            filterable
            clearable
            remote
            @search="loadParentTask"
            :on-update:value="parentTaskChanged"
          />
        </n-form-item>
        <n-form-item label="功能模块" path="moduleId">
          <n-cascader
            v-model:value="currentData.moduleId"
            :options="moduleTree"
            value-field="id"
            label-field="name"
            children-field="children"
            placeholder="请选择功能模块"
            clearable
            show-path
            :on-update:value="loadRequirementList"
            :disabled="!!currentData.parentId"
          />
        </n-form-item>
        <n-form-item label="所属需求" path="requirementId">
          <n-select
            v-model:value="currentData.requirementId"
            placeholder="请选择任务所属需求"
            :options="requirementList"
            label-field="name"
            value-field="id"
            filterable
            :on-update:value="updateTaskDescByRequirement"
            :disabled="!!currentData.parentId"
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
        <n-form-item label="预期任务时间" path="startTime">
          <n-date-picker v-model:value="taskTime" type="datetimerange" clearable placeholder="请选择预期任务时间" />
        </n-form-item>
        <n-form-item label="任务参与人" path="members">
          <n-select
            v-model:value="memberIdList"
            placeholder="请选择任务参与人"
            :options="memberList"
            label-field="realName"
            value-field="userId"
            filterable
            multiple
            clearable
          />
        </n-form-item>
        <n-form-item label="详情描述" path="taskDesc">
          <n-input v-model:value="currentData.taskDesc" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入描述" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData">重置</n-button>
        <n-button size="small" type="primary" @click="handleCommitData" v-resource-code="['project:detail:task:add', 'project:detail:task:update']"
          >提交</n-button
        >
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
