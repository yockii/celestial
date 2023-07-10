<template>
  <DocumentEditor id="docEditor" documentServerUrl="http://192.168.1.117:8700/" :config="config" :events_onDocumentReady="onDocumentReady" />
</template>

<script setup lang="ts">
import { DocumentEditor } from "@onlyoffice/document-editor-vue"
const config = ref({
  document: {
    fileType: "docx",
    key: "Khirz6zTPdfd7",
    title: "Example Document Title.docx",
    url: "https://example.com/url-to-example-document.docx"
  },
  documentType: "word", // cell, slide
  // token: jwtToken,
  // type: "embedded",
  editorConfig: {
    callbackUrl: "https://example.com/url-to-callback.ashx",
    lang: "zh-CN"
    // user: {} // 用户信息
  }
})
const onDocumentReady = () => {
  console.log("Document is ready")
}
</script>
