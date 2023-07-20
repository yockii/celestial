<script setup lang="ts">
import { useAppStore } from "@/store/app"
import { useUserStore } from "@/store/user"
import { Edit } from "@vicons/carbon"
import { storeToRefs } from "pinia"
import Vditor from "vditor"
import "vditor/dist/index.css"
const appStore = useAppStore()
const userStore = useUserStore()
const { theme } = storeToRefs(appStore)
const props = defineProps<{
  value: string
  label?: string
  placeholder?: string
}>()
const emit = defineEmits(["update:value"])

const showModal = ref(false)
const switchModal = (show: boolean, cancel = false) => {
  showModal.value = show
  if (showModal.value) {
    nextTick(() => {
      initVditor()
    })
  } else {
    if (!cancel) {
      const v = vditor.value?.getValue()
      if (v) {
        emit("update:value", v)
      }
    }
    vditor.value?.destroy()
  }
}
// vditor
const vditor = ref<Vditor | null>(null)
const initVditor = () => {
  vditor.value = new Vditor("vditor", {
    theme: theme.value === "dark" ? "dark" : "classic",
    height: "100vh",
    width: "100vw",
    placeholder: props.placeholder || "请输入内容",
    after: () => {
      vditor.value?.setValue(props.value)
    },
    upload: {
      url: "/api/v1/file",
      headers: {
        Authorization: `Bearer ${userStore.token}`
      },
      format: (files: File[], responseText: string): string => {
        const response = JSON.parse(responseText)
        const resultJson: { [key: string]: any } = {}
        if (response.code === 0) {
          const succMap: { [key: string]: string } = {}
          for (const key in response.data.success) {
            succMap[key] = window.location.origin + "/api/v1/file?objName=" + response.data.success[key]
          }

          resultJson["msg"] = ""
          resultJson["code"] = 0
          resultJson["data"] = {
            errFiles: response.data.failed,
            succMap
          }
        } else {
          const errFiles = []
          for (const f of files) {
            errFiles.push(f.name)
          }
          resultJson["msg"] = response.msg
          resultJson["code"] = response.code
          resultJson["data"] = {
            errFiles
          }
        }

        return JSON.stringify(resultJson)
      }
    }
  })
}
// TODO 更新值
// TODO 主题同步切换？
</script>

<template>
  <div class="relative w-full">
    <n-input type="textarea" :value="props.value" :placeholder="props.placeholder" @update:value="emit('update:value', $event)" />
    <n-icon class="absolute top-5px right-5px cursor-pointer" @click="switchModal(true)">
      <Edit />
    </n-icon>
  </div>
  <n-modal :show="showModal">
    <div class="relative">
      <n-button type="default" class="absolute bottom-10px right-80px z-999" @click="switchModal(false, true)"> 取消 </n-button>
      <n-button type="primary" class="absolute bottom-10px right-10px z-999" @click="switchModal(false)"> 确定 </n-button>
      <div id="vditor"></div>
    </div>
  </n-modal>
</template>
