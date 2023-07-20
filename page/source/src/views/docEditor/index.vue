<template>
  <div class="wscreen hscreen">
    <DocumentEditor v-if="configReady" id="docEditor" :documentServerUrl="documentServerUrl" :config="config" :events_onDocumentReady="onDocumentReady" />
  </div>
</template>

<script setup lang="ts">
import { getFileConfig } from "@/service"
import { DocumentEditor, IConfig } from "@onlyoffice/document-editor-vue"
const route = useRoute()

const documentServerUrl = import.meta.env.VITE_EDITOR_URI

const configReady = ref(false)

const config = ref<IConfig>({
  document: {
    fileType: "docx",
    key: "",
    title: "",
    url: ""
  }
})

const onDocumentReady = (e: any) => {
  console.log(e)
}

onMounted(() => {
  const fileId = route.params.id as string
  const versionId = route.params.versionId as string
  getFileConfig(fileId, versionId).then((res) => {
    const conf = res
    const currentUrl = `${window.location.protocol}//${window.location.host}`
    // conf.document.url = `${window.location.protocol}://${window.location.host}/api/v1/assetFile/download?id=${conf.document.key}`
    // conf.editorConfig.callbackUrl = `${window.location.protocol}://${window.location.host}/api/v1/office/callback`
    config.value = Object.assign(config.value, conf)
    if (config.value.editorConfig) {
      config.value.editorConfig.customization = {
        logo: {
          image: currentUrl + "/logo_172_40.png",
          imageDark: currentUrl + "/logo_172_40_dark.png",
          url: "https://oa.xhnic.com/"
        },
        customer: {
          address: "浙江省杭州市西湖区",
          info: "西湖新基建",
          logo: currentUrl + "/logo_432_70.png",
          logoDark: currentUrl + "/logo_432_70_dark.png",
          mail: "xuyq@xhnewi.com",
          name: "西湖新基建",
          phone: "0571-87609880",
          www: "www.xhnewi.com"
        }
      }
    }
    configReady.value = true
  })
})
</script>
