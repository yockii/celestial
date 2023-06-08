<script setup lang="ts">
import { ref, reactive, computed, onMounted, h } from "vue"
import { Search } from "@vicons/carbon"
import { OssConfig, OssConfigCondition } from "@/types/ossConfig"
import dayjs from "dayjs"
import { NButtonGroup, NButton, NPopconfirm, FormInst, useMessage, PaginationProps, NSwitch } from "naive-ui"
import { addOssConfig, deleteOssConfig, getOssConfigDetail, getOssConfigList, updateOssConfig, updateOssConfigStatus } from "@/service/api/settings/ossConfig"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
const condition = ref<OssConfigCondition>({
  offset: 0,
  limit: 10
})
const list = ref<OssConfig[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getOssConfigList(condition.value)
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
const statusColumn = {
  title: "状态",
  key: "status",
  render: (row: OssConfig) => {
    // return row.status === 1 ? "启用" : "禁用"
    return h(
      NSwitch,
      {
        value: row.status === 1,
        onUpdateValue: (value: boolean) => {
          updateOssConfigStatus(row.id, value ? 1 : -1).then((res) => {
            if (res) {
              message.success("修改成功")
              refresh()
            }
          })
        }
      },
      {
        checked: () => "启用",
        unchecked: () => "禁用"
      }
    )
  }
}
const columns = [
  {
    title: "名称",
    key: "name"
  },
  {
    title: "类型",
    key: "type"
  },
  {
    title: "桶",
    key: "bucket"
  },
  {
    title: "Endpoint",
    key: "endpoint"
  },
  statusColumn,
  {
    title: "创建时间",
    key: "createTime",
    // 时间戳转换为 yyyy-MM-dd HH:mm:ss的形式
    render: (row: OssConfig) => dayjs(row.createTime).fromNow()
  },
  {
    title: "操作",
    key: "operation",
    // 返回VNode, 用于渲染操作按钮
    render: (row: OssConfig) => {
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
const handleEditData = (row: OssConfig) => {
  getOssConfigDetail(row.id).then((res) => {
    if (res) {
      checkedData.value = res
      drawerActive.value = true
    }
  })
}
const handleDeleteData = (id: string) => {
  deleteOssConfig(id).then((res) => {
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
  checkedData.value = { id: "", name: "", endpoint: "", bucket: "", type: "" }
}
const handleAddOssConfig = () => {
  resetCheckedData()
  drawerActive.value = true
}
// 抽屉相关逻辑，用于新增/修改
const drawerActive = ref(false)
const isUpdate = computed(() => !!checkedData.value?.id)
const drawerTitle = computed(() => (isUpdate.value ? "修改oss配置" : "新增oss配置"))
const checkedData = ref<OssConfig>({ id: "", type: "", name: "" })
const formRef = ref<FormInst | null>(null)
const handleCommitData = () => {
  formRef.value?.validate((errors) => {
    if (errors) {
      return
    }
    if (isUpdate.value) {
      // 更新数据
      updateOssConfig(checkedData.value).then((res) => {
        if (res) {
          message.success("修改成功")
          refresh()
          drawerActive.value = false
        }
      })
    } else {
      // 新增数据
      addOssConfig(checkedData.value).then((res) => {
        if (res) {
          message.success("新增成功")
          refresh()
          drawerActive.value = false
        }
      })
    }
  })
}
const ossConfigTypeOptions = [
  {
    label: "华为云OBS",
    value: "obs"
  },
  {
    label: "阿里云OSS",
    value: "oss"
  },
  {
    label: "minio",
    value: "minio"
  }
]
// 规则定义
const rules = {
  name: {
    required: true,
    message: "请输入名称",
    min: 2,
    max: 20,
    trigger: "blur"
  },
  type: [
    {
      required: true,
      message: "请选择类型",
      trigger: "blur"
    }
  ],
  endpoint: [
    {
      required: true,
      message: "请输入Endpoint",
      trigger: "blur"
    }
  ],
  bucket: [
    {
      required: true,
      message: "请输入桶",
      trigger: "blur"
    }
  ],
  accessKeyId: [
    {
      required: true,
      message: "请输入accessKeyId",
      trigger: "blur"
    }
  ],
  secretAccessKey: [
    {
      required: true,
      message: "请输入secretAccessKey",
      trigger: "blur"
    }
  ]
}
</script>

<template>
  <n-grid :cols="1" y-gap="8">
    <n-gi>
      <n-grid :cols="2">
        <n-gi>
          <n-h3>oss配置管理</n-h3>
        </n-gi>
        <n-gi class="flex flex-justify-end">
          <n-button size="small" type="primary" @click="handleAddOssConfig" v-resource-code="'system:ossConfig:add'">新增</n-button>
        </n-gi>
      </n-grid>
    </n-gi>
    <n-gi>
      <n-grid :cols="8" x-gap="8">
        <n-gi :span="2">
          <n-input size="small" placeholder="名称" v-model:value="condition.name" @keydown.enter.prevent="refresh" />
        </n-gi>
        <n-gi>
          <n-input size="small" placeholder="类型" v-model:value="condition.type" @keydown.enter.prevent="refresh" />
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
        :row-key="(row: OssConfig) => row.id"
        :pagination="pagination"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :columns="columns"
      />
    </n-gi>
  </n-grid>

  <n-drawer v-model:show="drawerActive" :width="401">
    <n-drawer-content :title="drawerTitle" closable>
      <n-form size="small" ref="formRef" :model="checkedData" :rules="rules" label-width="90px">
        <n-form-item label="名称" path="name" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="类型" path="type" label-placement="left" label-align="left">
          <n-select v-model:value="checkedData.type" placeholder="请选择oss类型" :options="ossConfigTypeOptions" />
        </n-form-item>
        <n-form-item label="桶名称" path="bucket" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.bucket" placeholder="请输入桶名称" />
        </n-form-item>
        <n-form-item label="区域编码" path="region" label-placement="left" label-align="left">
          <n-input v-model:value="checkedData.region" placeholder="请输入区域编码" />
        </n-form-item>
        <n-form-item label="是否HTTPS" path="selfDomain" label-placement="left" label-align="left">
          <n-radio-group v-model:value="checkedData.selfDomain" name="selfDomain">
            <n-space>
              <n-radio :value="1">是</n-radio>
              <n-radio :value="2">否</n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="自定义域名" path="secure" label-placement="left" label-align="left">
          <n-radio-group v-model:value="checkedData.secure" name="secure">
            <n-space>
              <n-radio :value="1">是</n-radio>
              <n-radio :value="2">否</n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
        <n-form-item label="Endpoint" path="endpoint">
          <n-input v-model:value="checkedData.endpoint" placeholder="请输入Endpoint" />
        </n-form-item>
        <n-form-item label="accessKeyId" path="accessKeyId">
          <n-input v-model:value="checkedData.accessKeyId" placeholder="请输入accessKeyId" />
        </n-form-item>
        <n-form-item label="secretAccessKey" path="secretAccessKey">
          <n-input v-model:value="checkedData.secretAccessKey" placeholder="请输入secretAccessKey" />
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
