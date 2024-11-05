<script lang="ts" setup>
import {useProjectStore} from "@/store/project";
import {ScreenSearchDesktopFilled} from '@vicons/material'
import {SendAlt, SendAltFilled} from "@vicons/carbon"
import {Robot} from "@vicons/tabler"
import {getAiData} from "@/service";
import { MdPreview } from 'md-editor-v3';

const projectStore = useProjectStore()
const showModal = ref(false)
const question = ref("")
const loading = ref(false)

type AIContent = {
  isMe: boolean
  content: string
}
const contentList = ref<AIContent[]>([
  {isMe: false, content: "你好，我是AI查询小助手，你可以向我提出查询需求，我尽力查找本系统的数据并反馈给你。"},
])
const send = async () => {
  if(question.value === '') {
    return
  }
  try {
    contentList.value.push({isMe: true, content: question.value})
    const robotContent = {isMe: false, content: "请稍等，我正在查询中..."}
    contentList.value.push(robotContent)
    loading.value = true
    const res = await getAiData(question.value)
    if (res) {
      robotContent.content = res
    } else {
      robotContent.content = "查询失败，请稍后再试"
    }
    question.value = ""
    loading.value = false
  } catch (e) {
    console.log(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <n-tooltip trigger="hover" placement="right" v-if="projectStore.hasResourceCodes(['ai'])">
    <template #trigger>
      <n-float-button class="z-999 right-60px bottom-100px" position="absolute" @click="showModal = true">
        <n-icon>
          <ScreenSearchDesktopFilled/>
        </n-icon>
      </n-float-button>
    </template>
    AI查询小助手
  </n-tooltip>

  <n-drawer v-model:show="showModal" :width="600" closable>
    <n-drawer-content>
      <div class="h-100% flex flex-col">
        <div class="flex-1">
          <div v-for="(c, i) in contentList" :key="i" class="flex items-start mb-16px" :class="c.isMe ? 'justify-end' : 'justify-start'">
            <div class="min-w-40px">
              <n-avatar v-if="!c.isMe">
                <n-icon>
                  <Robot />
                </n-icon>
              </n-avatar>
            </div>
            <div class="flex-1 flex" :class="c.isMe ? 'justify-end' : 'justify-start'">
              <!-- 模拟对话框 -->
              <div class="content" :class="c.isMe ? 'my-content' : 'robot-content'">
                <div class="ma-8px" :class="c.isMe ? 'mr-16px' : 'ml-24px'">
<!--                  {{ c.content }}-->
                  <md-preview :model-value="c.content" />
                </div>
              </div>
            </div>
            <div class="min-w-40px">
              <n-avatar v-if="c.isMe" round>我</n-avatar>
            </div>
          </div>
        </div>
        <div class="flex gap-16px items-center">
          <n-input v-model:value="question" autofocus @keydown.enter="send" :disabled="loading" :loading="loading"/>
          <n-button @click="send" type="info" :loading="loading" :disabled="loading || question === ''">
            <template #icon>
              <SendAltFilled v-if="question !== ''"/>
              <SendAlt v-else/>
            </template>
          </n-button>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped>
.content {
  width: 400px;
}
.my-content {
  margin-right: 4px;
  background-color: #e6f7ff;
  border: 1px solid #e6f7ff;
  clip-path: polygon(0 0, 384px 0, 384px 8px, 100% 16px, 384px 24px, 384px 100%, 0 100%);
}
.robot-content {
  margin-left: 4px;
  background-color: #C7EDCC;
  border: 1px solid #C7EDCC;
  clip-path: polygon(16px 0, 100% 0, 100% 100%, 16px 100%, 16px 24px, 0 16px, 16px 8px);
}
</style>