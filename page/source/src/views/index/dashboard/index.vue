<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-input
        clearable
        passively-activated
        :loading="searching"
        placeholder="请输入搜索内容"
        show-count
        v-model:value="keyword"
        @keydown.enter.prevent="search"
      >
        <template #suffix>
          <n-icon :component="Search" class="cursor-pointer" @click="search" />
        </template>
      </n-input>
    </n-gi>
    <n-gi v-if="luceneDocumentList.length > 0">
      <n-list clickable hoverable>
        <template v-if="total > luceneDocumentList.length" #footer>
          <n-space>
            <n-button :loading="searching" :disabled="searching" @click="loadData" type="primary" size="small"> 加载更多 </n-button>
          </n-space>
        </template>
        <n-list-item v-for="ld in luceneDocumentList" :key="ld.id" @click="router.push(ld.route)">
          <n-thing :title="ld.title">
            <p class="ws-pre-line">{{ ld.content }}</p>
            <template #footer>
              <n-space>
                <span>创建时间：{{ dayjs(ld.createTime).fromNow() }}</span>
                <span>更新时间：{{ dayjs(ld.updateTime).fromNow() }}</span>
              </n-space>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-gi>
  </n-grid>
</template>

<script setup lang="ts">
import { luceneSearch } from "@/service/api"
import { useAppStore } from "@/store/app"
import { LuceneDocument } from "@/types/lecene"
import { Search } from "@vicons/carbon"
import dayjs from "dayjs"
const router = useRouter()
const keyword = ref("")
const luceneDocumentList = ref<LuceneDocument[]>([])
const total = ref(0)
const searching = ref(false)
const search = () => {
  luceneDocumentList.value = []
  loadData()
}
const loadData = () => {
  searching.value = true
  luceneSearch(keyword.value, 10, luceneDocumentList.value.length)
    .then((res) => {
      total.value = res.total
      luceneDocumentList.value = luceneDocumentList.value.concat(res.items || [])
    })
    .finally(() => {
      searching.value = false
    })
}
</script>

<style scoped></style>
