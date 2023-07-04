<script setup lang="ts">
import { getAssetFile, updateAssetFile, getAssetCategoryList } from "@/service/api"
import { File } from "@/types/asset"
import { FormInst, FormItemRule, UploadFileInfo, UploadInst } from "naive-ui"
import { ArchiveOutline } from "@vicons/ionicons5"
import { useUserStore } from "@/store/user"
import { AssetCategory } from "@/types/asset"

const message = useMessage()
const userStore = useUserStore()

const emit = defineEmits(["update:drawerActive", "update:data", "refresh"])
const props = defineProps<{
  drawerActive: boolean
  data: File
}>()

const currentData = computed({
  get: () => {
    return props.data
  },
  set: (val: File) => {
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
const upload = ref<UploadInst | undefined>()
const handleCommitData = () => {
  formRef.value?.validate().then(() => {
    if (isUpdate.value) {
      updateAssetFile(currentData.value).then((res) => {
        if (res) {
          message.success("更新成功")
          emit("refresh")
          emit("update:drawerActive", false)
        }
      })
    } else {
      upload.value?.submit()
    }
  })
}

const categoryList = ref<AssetCategory[]>([])

// 上传文件
const fileList = ref<UploadFileInfo[]>([])
watch(
  () => currentData.value.id,
  (val) => {
    if (val) {
      getAssetFile(val).then((res) => {
        fileList.value = [
          {
            id: val,
            name: res.name,
            status: "finished"
          }
        ]
      })
    } else {
      fileList.value = []
    }
  }
)
const handleBeforeUpload = () => {
  if (!currentData.value.name) {
    message.error("请先输入资产名称")
    return false
  }
  if (!currentData.value.categoryId) {
    message.error("请选择资产目录")
    return false
  }

  return true
}
const handleUploadError = (options: { file: UploadFileInfo; event?: ProgressEvent }) => {
  const req = options.event?.target as XMLHttpRequest
  let msg = "上传文件失败"
  if (req.status === 200) {
    msg += "上传文件失败, " + req.response?.msg
  }
  message.error(msg)
  options.file.status = "error"
}
const handleUploadFinished = (options: { file: UploadFileInfo; event?: ProgressEvent }) => {
  const req = options.event?.target as XMLHttpRequest
  currentData.value.id = req.response.data.id
  options.file.name = req.response.data.name

  message.success("上传成功")
  emit("refresh")
  emit("update:drawerActive", false)
}

// 规则定义
const rules = {
  name: [
    { required: true, message: "请输入资产名称", trigger: "blur" },
    { min: 2, max: 20, message: "长度在 2 到 20 个字符", trigger: "blur" }
  ],
  type: {
    required: true,
    validator(rule: FormItemRule, value: number) {
      return value > 0 && (value <= 4 || value === 9)
    },
    message: "请选择资产类型",
    trigger: "blur"
  }
}

// 加载资产目录
const loadAssetCategory = () => {
  getAssetCategoryList({
    onlyParent: true,
    offset: -1,
    limit: -1,
    type: 2
  }).then((res) => {
    categoryList.value = res.items.map((item) => {
      item.isLeaf = !item.childrenCount || item.childrenCount === 0
      return item
    })
  })
}
const loadCategory = (option: AssetCategory) => {
  return new Promise<void>((resolve) => {
    getAssetCategoryList({
      parentId: option.id,
      offset: -1,
      limit: -1,
      type: 2
    }).then((res) => {
      option.children = res.items || []
      resolve()
    })
  })
}

onMounted(() => {
  loadAssetCategory()
})
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" @update:show="(show:boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="drawerTitle">
      <n-form ref="formRef" :model="currentData" label-width="100px" :rules="rules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="currentData.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="资产目录" required>
          <n-cascader
            v-model:value="currentData.categoryId"
            placeholder="请选择资产目录"
            :options="categoryList"
            label-field="name"
            value-field="id"
            clearable
            remote
            :on-load="loadCategory"
          />
        </n-form-item>
        <n-form-item label="附件" path="attachment" required v-if="!isUpdate">
          <n-upload
            ref="upload"
            action="/api/v1/assetFile/add"
            :default-upload="false"
            multiple
            directory-dnd
            accept="application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/pdf,image/*,.csv,text/plain,audio/*,application/zip,application/x-*-compressed"
            :max="1"
            :headers="{
              authorization: 'Bearer ' + userStore.token
            }"
            :data="{
              categoryId: currentData.categoryId,
              name: currentData.name
            }"
            response-type="json"
            :is-error-state="(xhr: XMLHttpRequest) => xhr.status !== 200 || xhr.response.code !== 0"
            v-model:file-list="fileList"
            :on-before-upload="handleBeforeUpload"
            :on-finish="handleUploadFinished"
            :on-error="handleUploadError"
          >
            <n-upload-dragger v-if="fileList.length === 0">
              <div class="mb-12px">
                <n-icon size="48" :depth="1">
                  <ArchiveOutline />
                </n-icon>
              </div>
              <n-text class="">点击或拖拽文件到此处上传</n-text>
              <n-p depth="3" class="mt-8px"> </n-p>
            </n-upload-dragger>
          </n-upload>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" v-if="!isUpdate" @click="resetCurrentData">重置</n-button>
        <n-button
          size="small"
          type="primary"
          @click="handleCommitData"
          :disabled="!isUpdate && !fileList.length"
          v-resource-code="['asset:file:add', 'asset:file:update']"
          >提交</n-button
        >
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
