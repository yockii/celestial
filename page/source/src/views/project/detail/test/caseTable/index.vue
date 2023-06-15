<script setup lang="ts">
import { useProjectStore } from "@/store/project"
import { useUserStore } from "@/store/user"
import { ProjectTest, ProjectTestCase, ProjectTestCaseCondition, ProjectTestCaseItem } from "@/types/project"
import { NGrid, NButtonGroup, NPopconfirm, NButton, NIcon, NTooltip, NSpace, NText } from "naive-ui"
import { storeToRefs } from "pinia"
import { Edit, PlaylistAdd, Trash } from "@vicons/tabler"
import {
  deleteProjectTestCase,
  deleteProjectTestCaseItem,
  getProjectTestCaseListWithItems,
  updateProjectTestCaseItemStatus
} from "@/service/api/project/projectTest"
import CaseDrawer from "../caseDrawer/index.vue"
import ItemDrawer from "../itemDrawer/index.vue"
import { AiStatusComplete, AiStatusFailed } from "@vicons/carbon"

const message = useMessage()
const userStore = useUserStore()
const { project } = storeToRefs(useProjectStore())
type CombinedTestData = {
  testCase: ProjectTestCase
  item: ProjectTestCaseItem
}

const props = defineProps<{
  test: ProjectTest | null
  testCaseList?: ProjectTestCase[]
}>()

const dynamicTestCaseList = ref<ProjectTestCase[]>([])

const remote = computed(() => !props.testCaseList || props.testCaseList.length === 0)
const originList = computed(() => (remote.value ? dynamicTestCaseList.value : props.testCaseList))
const list = computed(() => {
  const result: CombinedTestData[] = []
  if (!originList.value) return result
  for (const tc of originList.value) {
    if (tc.items?.length) {
      for (const item of tc.items) {
        result.push({
          testCase: tc,
          item
        })
      }
    } else {
      result.push({
        testCase: tc,
        item: {}
      })
    }
  }
  return result
})
const loading = ref(false)

const handleDeleteCaseData = (id: string) => {
  // 删除用例
  deleteProjectTestCase(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}

const caseColumn = {
  title: "用例名称",
  key: "name",
  rowSpan: (row: CombinedTestData) => row.testCase.items?.length || 1,
  render(row: CombinedTestData) {
    const btnGroup: VNode[] = []
    if (remote.value) {
      btnGroup.push(
        h(
          NButtonGroup,
          {},
          {
            default: () => [
              h(
                NTooltip,
                {},
                {
                  default: () => "新增用例项",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        type: "primary",
                        size: "small",
                        disabled: !userStore.hasResourceCode("project:detail:test:testCase:item:add"),
                        onClick: () => {
                          handleNewItem(row.testCase.id)
                        }
                      },
                      {
                        icon: () =>
                          h(
                            NIcon,
                            {},
                            {
                              default: () => h(PlaylistAdd)
                            }
                          )
                      }
                    )
                }
              ),
              h(
                NTooltip,
                {},
                {
                  default: () => "编辑用例",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        type: "primary",
                        size: "small",
                        secondary: true,
                        disabled: !userStore.hasResourceCode("project:detail:test:testCase:update"),
                        onClick: () => {
                          currentTestCase.value = row.testCase
                          caseDrawerActive.value = true
                        }
                      },
                      {
                        icon: () => h(NIcon, {}, { default: () => h(Edit) })
                      }
                    )
                }
              ),
              h(
                NTooltip,
                {},
                {
                  default: () => "删除用例",
                  trigger: () =>
                    h(
                      NPopconfirm,
                      {
                        onPositiveClick: () => handleDeleteCaseData(row.testCase.id)
                      },
                      {
                        default: () => "确认删除该用例吗？",
                        trigger: () =>
                          h(
                            NButton,
                            {
                              type: "error",
                              disabled: !userStore.hasResourceCode("project:detail:test:testCase:delete"),
                              size: "small"
                            },
                            { icon: () => h(NIcon, {}, { default: () => h(Trash) }) }
                          )
                      }
                    )
                }
              )
            ]
          }
        )
      )
    }

    return h(
      NSpace,
      {
        justify: "space-between"
      },
      {
        default: () => [
          h(
            NTooltip,
            {},
            {
              default: () => row.testCase.remark || row.testCase.name,
              trigger: () =>
                h(
                  NText,
                  {
                    type: row.testCase.items
                      ? row.testCase.items?.filter((item) => item.status === -1).length
                        ? "error"
                        : row.testCase.items?.filter((item) => item.status === 1).length
                        ? "default"
                        : "success"
                      : "default"
                  },
                  { default: () => row.testCase.name }
                )
            }
          ),
          btnGroup
        ]
      }
    )
  }
}

const handleDeleteItemData = (id: string) => {
  if (id === "") return
  // 删除用例项
  deleteProjectTestCaseItem(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}
const handleTestItemPass = (id: string | undefined, pass: boolean) => {
  if (!id) return
  // 删除用例项
  updateProjectTestCaseItemStatus(id, pass ? 2 : -1).then((res) => {
    if (res) {
      message.success("操作成功")
      refresh()
    }
  })
}
const caseItemColumn = {
  title: "用例项",
  key: "items",
  render(row: CombinedTestData) {
    if (!row.item.id) {
      return h(
        "div",
        {},
        {
          default: () => "暂无用例项"
        }
      )
    }

    const btnGroup: VNode[] = []
    if (remote.value) {
      if (!row.item.status || row.item.status === 1) {
        btnGroup.push(
          h(
            NButtonGroup,
            {},
            {
              default: () => [
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "测试通过",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "primary",
                          disabled: !userStore.hasResourceCode("project:detail:test:testCase:item:updateStatus"),
                          onClick: () => handleTestItemPass(row.item.id, true)
                        },
                        {
                          default: () => h(NIcon, { component: AiStatusComplete })
                        }
                      )
                  }
                ),
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "测试不通过",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          type: "error",
                          disabled: !userStore.hasResourceCode("project:detail:test:testCase:item:updateStatus"),
                          onClick: () => handleTestItemPass(row.item.id, false)
                        },
                        {
                          default: () => h(NIcon, { component: AiStatusFailed })
                        }
                      )
                  }
                )
              ]
            }
          )
        )
      }
      btnGroup.push(
        h(
          NButtonGroup,
          {},
          {
            default: () => [
              h(
                NTooltip,
                {},
                {
                  default: () => "编辑用例项",
                  trigger: () =>
                    h(
                      NButton,
                      {
                        size: "small",
                        type: "primary",
                        secondary: true,
                        disabled: !userStore.hasResourceCode("project:detail:test:testCase:item:update"),
                        onClick: () => {
                          currentItem.value = row.item
                          itemDrawerActive.value = true
                        }
                      },
                      {
                        icon: () => h(NIcon, {}, { default: () => h(Edit) })
                      }
                    )
                }
              ),
              h(
                NTooltip,
                {},
                {
                  default: () => "删除用例项",
                  trigger: () =>
                    h(
                      NPopconfirm,
                      {
                        onPositiveClick: () => handleDeleteItemData(row.item.id || "")
                      },
                      {
                        default: () => "确认删除该用例项吗？",
                        trigger: () =>
                          h(
                            NButton,
                            {
                              type: "error",
                              disabled: !userStore.hasResourceCode("project:detail:test:testCase:item:delete"),
                              size: "small"
                            },
                            { icon: () => h(NIcon, {}, { default: () => h(Trash) }) }
                          )
                      }
                    )
                }
              )
            ]
          }
        )
      )
    }

    return h(
      NSpace,
      {
        justify: "space-between"
      },
      {
        default: () => [
          h(
            NTooltip,
            {},
            {
              default: () => row.item.content || row.item.name || "",
              trigger: () =>
                h(NText, { type: row.item.status === -1 ? "error" : row.item.status === 2 ? "success" : "default" }, { default: () => row.item.name || "" })
            }
          ),
          h(
            NSpace,
            {},
            {
              default: () => btnGroup
            }
          )
        ]
      }
    )
  }
}
// 列定义
const columns = [caseColumn, caseItemColumn]
// 条件
const condition = ref<ProjectTestCaseCondition>({
  projectId: project.value.id
})
const refresh = () => {
  loading.value = true
  getProjectTestCaseListWithItems(condition.value)
    .then((res) => {
      if (res) {
        dynamicTestCaseList.value = res.items || []
        pagination.itemCount = res.total
        pagination.pageCount = Math.ceil(res.total / (condition.value.limit || 10))
        pagination.page = Math.ceil((condition.value.offset || 0) / (condition.value.limit || 10)) + 1
      }
    })
    .finally(() => {
      loading.value = false
    })
}
// 分页
const pagination = reactive({
  itemCount: 0,
  pageCount: 0,
  page: 1,
  pageSize: 10,
  onChange: (page: number) => {
    if (remote) {
      condition.value.offset = (page - 1) * (condition.value.limit || 10)
    } else {
      pagination.page = page
    }
  },
  onUpdatePageSize: (pageSize: number) => {
    if (remote) {
      condition.value.limit = pageSize
    } else {
      pagination.page = 1
      pagination.pageSize = pageSize
    }
  }
})

// 新增测试用例
const caseDrawerActive = ref(false)
const currentTestCase = ref<ProjectTestCase>({ id: "", projectId: project.value.id, name: "" })
const handleNewTestCase = () => {
  currentTestCase.value = { id: "", projectId: project.value.id, name: "" }
  caseDrawerActive.value = true
}

// 新增测试用例项
const itemDrawerActive = ref(false)
const currentItem = ref<ProjectTestCaseItem>({ id: "", testCaseId: "", name: "", content: "" })
const handleNewItem = (caseId: string) => {
  currentItem.value = { id: "", testCaseId: caseId, projectId: project.value.id, name: "", content: "" }
  itemDrawerActive.value = true
}

onMounted(() => {
  refresh()
})
defineExpose({
  refresh
})
</script>

<template>
  <n-empty v-if="test && test.endTime && test.endTime > 0 && (!testCaseList || testCaseList.length === 0)" description="暂无数据"> </n-empty>
  <n-grid v-else :cols="1" y-gap="16">
    <n-gi v-if="!testCaseList || testCaseList.length === 0">
      <n-space justify="space-between">
        <span></span>
        <n-button type="primary" @click="handleNewTestCase">新增用例</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        :remote="remote"
        :data="list"
        :loading="loading"
        :row-key="(row: CombinedTestData) => row.item.id"
        :pagination="pagination"
        :columns="columns"
        :single-line="false"
      />
    </n-gi>
  </n-grid>

  <case-drawer v-model:drawer-active="caseDrawerActive" :data="currentTestCase" @refresh="refresh" />
  <item-drawer v-model:drawer-active="itemDrawerActive" :data="currentItem" @refresh="refresh" />
</template>
