<script setup lang="ts">
import { Add } from "@vicons/carbon"
import { ProjectTest, ProjectTestCase } from "@/types/project"
import { addProjectTest, closeProjectTest, getProjectTest, getProjectTestList } from "@/service/api"
import { useProjectStore } from "@/store/project"
import { storeToRefs } from "pinia"
import dayjs from "dayjs"
import CaseTable from "./caseTable/index.vue"
import { CaretRight } from "@vicons/carbon"

const message = useMessage()
const projectStore = useProjectStore()
const { project } = storeToRefs(projectStore)
const testList = ref<ProjectTest[]>([])
const caseTableRef = ref<typeof CaseTable>()
const needCloseTest = computed(() => {
  // 判断是否需要关闭本轮测试
  if (testList.value.length === 0) {
    return false
  }
  return testList.value.find((item) => !item.endTime || item.endTime === 0)
})

const refreshTest = () => {
  getProjectTestList(project.value.id).then((res) => {
    if (res) {
      testList.value = [...res].reverse()
    }
  })
}

// 当前选中的测试
const currentTest = ref<ProjectTest | null>(null)
const currentTestCaseList = ref<ProjectTestCase[]>([])
const handleSelectTest = (id: string | undefined) => {
  if (!id) {
    return
  }
  const test = testList.value.find((item) => item.id === id)
  if (test) {
    currentTest.value = test
    if (test.endTime && test.endTime > 0) {
      // 该测试已封版，获取完整的测试记录
      getProjectTest(id).then((res) => {
        if (res && res.testRecord) {
          currentTestCaseList.value = JSON.parse(res.testRecord)
        }
      })
    } else {
      // 该测试未封版，让测试用例列表为空，由组件自行获取动态数据
      currentTestCaseList.value = []
    }
  }
}
const openNewRound = () => {
  // 开启新一轮测试
  message.info("正在开启新一轮测试")
  addProjectTest({
    projectId: project.value.id
  }).then((res) => {
    if (res) {
      refreshTest()
      caseTableRef.value?.refresh()
    }
  })
}
const handleNextRound = () => {
  if (needCloseTest.value && needCloseTest.value.id) {
    // 关闭本轮测试
    message.info(`正在关闭第${needCloseTest.value.round}轮测试`)
    closeProjectTest(needCloseTest.value.id).then((res) => {
      if (res) {
        // 开启新一轮测试
        openNewRound()
      }
    })
  } else {
    // 开启新一轮测试
    openNewRound()
  }
}

// 加载时
onMounted(() => {
  refreshTest()
})
</script>

<template>
  <n-grid :cols="6" x-gap="16">
    <n-gi>
      <n-h6 prefix="bar">
        <n-space justify="space-between">
          测试轮次
          <n-popconfirm @positive-click="handleNextRound" v-if="projectStore.hasResourceCode('project:detail:test:add')">
            <template #trigger>
              <n-tooltip>
                <template #trigger>
                  <n-button text>
                    <template #icon>
                      <n-icon>
                        <Add />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                开启下一轮
              </n-tooltip>
            </template>
            确定要{{ needCloseTest ? `关闭第${needCloseTest.round}轮并` : "" }}开启新一轮测试吗？
          </n-popconfirm>
        </n-space>
      </n-h6>
      <n-list hoverable clickable show-divider v-if="projectStore.hasResourceCode('project:detail:test:list')">
        <n-tooltip v-for="test in testList" :key="test.id">
          <template #trigger>
            <n-list-item @click="handleSelectTest(test.id)">
              <n-space justify="space-between">
                <n-text :type="!test.endTime || test.endTime === 0 ? 'info' : ''">第{{ test.round }}轮</n-text>
                <n-icon v-if="currentTest ? test.id === currentTest.id : !test.endTime || test.endTime === 0">
                  <CaretRight />
                </n-icon>
              </n-space>
            </n-list-item>
          </template>
          {{ dayjs(test.startTime).format("YYYY-MM-DD HH:mm:ss") }} ~ {{ test.endTime ? dayjs(test.endTime).format("YYYY-MM-DD HH:mm:ss") : "进行中" }}
        </n-tooltip>
      </n-list>
    </n-gi>
    <n-gi :span="5">
      <case-table ref="caseTableRef" :test="currentTest" :testCaseList="currentTestCaseList" />
    </n-gi>
  </n-grid>
</template>
