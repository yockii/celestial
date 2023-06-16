<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { Search } from "@vicons/carbon"
import { ThirdSource, ThirdSourceCondition } from "@/types/thirdSource"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, FormItemRule } from "naive-ui"
import { addThirdSource, deleteThirdSource, getThirdSourceDetail, getThirdSourceList, updateThirdSource } from "@/service/api"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<ThirdSourceCondition>({
  offset: 0,
  limit: 10
})
const list = ref<ThirdSource[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getThirdSourceList(condition.value)
    .then((res) => {
      list.value = res.items || []
      pagination.itemCount = res.total
      pagination.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
      pagination.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
    })
    .finally(() => {
      loading.value = false
    })
}

const pagination = reactive({
  itemCount: 0,
  page: 1,
  pageCount: 1,
  pageSize: 10,
  prefix({ itemCount }: PaginationProps) {
    return `共${itemCount}条`
  }
})
const handlePageChange = (page: number) => {
  condition.value.offset = (page - 1) * (condition.value.limit || 10)
  refresh()
}
const handlePageSizeChange = (pageSize: number) => {
  condition.value.limit = pageSize
  refresh()
}
const columns = [
  {
    title: "ID",
    key: "id"
  },
  {
    title: "名称",
    key: "name"
  },
  {
    title: "代码",
    key: "code"
  },
  {
    title: "企业ID",
    key: "corpId"
  },
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: ThirdSource) => dayjs(row.createTime).fromNow()
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: ThirdSource) => {
      return h(NButtonGroup, {}, () => [
        h(
          NButton,
          {
            size: "small",
            secondary: true,
            disabled: !userStore.hasResourceCode("system:thirdSource:update"),
            type: "primary",
            onClick: () => handleEditData(row)
          },
          {
            default: () => "编辑"
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteData(row.id)
          },
          {
            default: () => "确认删除",
            trigger: () =>
              h(
                NButton,
                {
                  size: "small",
                  disabled: !userStore.hasResourceCode("system:thirdSource:delete"),
                  type: "error"
                },
                {
                  default: () => "删除"
                }
              )
          }
        )
      ])
    }
  }
]
const handleEditData = (row: ThirdSource) => {
  getThirdSourceDetail(row.id).then((res) => {
    if (res) {
      checkedData.value = res
      drawerActive.value = true
    }
  })
}
const handleDeleteData = (id: string) => {
  deleteThirdSource(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
onMounted(() => {
  refresh()
})
const resetCheckedData = () => {
  checkedData.value = { id: "", name: "", code: "" }
}
const handleAddThirdSource = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改三方源" : "新增三方源"))
const checkedData = ref<ThirdSource>({ id: "", code: "", name: "" })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateThirdSource(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addThirdSource(checkedData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          refresh()
          drawerActive.value = false
        }
      })
    }
  })
}
// 规则定义
const rules = {
  name: {
    required: true,
    message: "请输入名称",
    min: 2,
    max: 10,
    trigger: "blur"
  },
  code: [
    {
      required: true,
      message: "请输入代码",
      min: 3,
      max: 20,
      trigger: "blur"
    }
  ],
  configuration: [
    // 校验json是否合法
    {
      required: true,
      message: "请输入正确的配置",
      trigger: "blur",
      validator: (rule: FormItemRule, value: string) => {
        try {
          JSON.parse(value)
          return true
        } catch (error) {
          /* empty */
        }
        return false
      }
    }
  ],
  matchConfig: [
    // 校验json是否合法
    {
      required: true,
      message: "请输入正确的配置",
      trigger: "blur",
      validator: (rule: FormItemRule, value: string) => {
        try {
          const obj = JSON.parse(value)
          if (obj["username"] && obj["match"]) {
            return true
          }
        } catch (error) {
          /* empty */
        }
        return false
      }
    }
  ]
}
const thirdSourceOptions = [
  {
    label: "钉钉",
    value: "dingtalk"
  }
]
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>三方源管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddThirdSource" v-resource-code="'system:thirdSource:add'">新增</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi>
          <n-input size="small" placeholder="代码" v-model:value="condition.code" @keydown.enter.prevent="refresh" />
        </n-gi>

        <n-gi :offset="4" class="flex flex-justify-end">
          <n-button size="small" @click="refresh">
            <template #icon>
              <n-icon><search /></n-icon>
            </template>
          </n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :row-key="(row: ThirdSource) => row.id"
        :pagination="pagination"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form ref="formRef" :model="checkedData" :rules="rules" label-width="100px" label-placement="left">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="checkedData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="代码" path="code">
          <n-select v-model:value="checkedData.code" placeholder="请选择三方源" :options="thirdSourceOptions" />
        </n-form-item>
        <n-form-item label="企业ID" path="corpId">
          <n-input v-model:value="checkedData.corpId" placeholder="请输入企业ID" />
        </n-form-item>
        <n-form-item label="配置" path="configuration">
          <n-input type="textarea" v-model:value="checkedData.configuration" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入配置" />
        </n-form-item>
        <n-form-item label="匹配配置" path="matchConfig">
          <n-input type="textarea" v-model:value="checkedData.matchConfig" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="请输入匹配配置" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCheckedData">重置</n-button>
        <n-button type="primary" @click="handleCommitData" v-resource-code="['system:thirdSource:add', 'system:thirdSource:update']">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
