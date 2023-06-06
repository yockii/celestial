<script setup lang="ts">
import { deleteCommonTestCase, deleteCommonTestCaseItem, getCommonTestCaseList } from "@/service/api/asset/commonTestCase"
import { CommonTestCase, CommonTestCaseCondition, CommonTestCaseItem } from "@/types/asset"
import { PaginationProps, NGrid, NGridItem, NButton, NPopconfirm, useMessage, NButtonGroup, NIcon, NTooltip } from "naive-ui"
import CaseDrawer from "./caseDrawer/index.vue"
import ItemDrawer from "./itemDrawer/index.vue"
import { Edit, PlaylistAdd, Trash } from "@vicons/tabler"
import { useUserStore } from "@/store/user"
const message = useMessage()
const userStore = useUserStore()
type CommonTestCaseItemData = {
  testCase: CommonTestCase
  testCaseItem: CommonTestCaseItem
}

const condition = ref<CommonTestCaseCondition>({})
const loading = ref(false)
const list = ref<CommonTestCase[]>([])
const data = computed(() => {
  const result = []
  for (const item of list.value) {
    if (item.items?.length) {
      for (const subItem of item.items) {
        result.push({
          testCase: item,
          testCaseItem: subItem
        })
      }
    } else {
      result.push({
        testCase: item,
        testCaseItem: {}
      })
    }
  }
  return result
})
const caseNameColumn = {
  title: "用例名称",
  key: "name",
  rowSpan: (row: CommonTestCaseItemData) => row.testCase.items?.length || 1,
  render(row: CommonTestCaseItemData) {
    return h(
      NGrid,
      {
        cols: 2
      },
      [
        h(NGridItem, null, {
          default: () =>
            h(
              NTooltip,
              {},
              {
                default: () => row.testCase.remark,
                trigger: () => row.testCase.name
              }
            )
        }),
        h(
          NGridItem,
          {
            class: "flex flex-row justify-end"
          },
          {
            default: () =>
              h(NButtonGroup, {}, [
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "新增用例项",
                    trigger: h(
                      NButton,
                      {
                        type: "primary",
                        size: "small",
                        disabled: !userStore.hasResourceCode("asset:commonTestCase:addItem"),
                        onClick: () => {
                          handleNewItem(row.testCase.id)
                        }
                      },
                      {
                        icon: h(
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
                    trigger: h(
                      NButton,
                      {
                        type: "primary",
                        size: "small",
                        secondary: true,
                        disabled: !userStore.hasResourceCode("asset:commonTestCase:update"),
                        onClick: () => {
                          currentCase.value = row.testCase
                          caseDrawerActive.value = true
                        }
                      },
                      {
                        icon: h(NIcon, {}, { default: () => h(Edit) })
                      }
                    )
                  }
                ),
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "删除用例",
                    trigger: h(
                      NPopconfirm,
                      {
                        onPositiveClick: () => handleDeleteData(row.testCase.id)
                      },
                      {
                        default: () => "确认删除该用例吗？",
                        trigger: () =>
                          h(
                            NButton,
                            {
                              type: "error",
                              disabled: !userStore.hasResourceCode("asset:commonTestCase:delete"),
                              size: "small"
                            },
                            { icon: h(NIcon, {}, { default: () => h(Trash) }) }
                          )
                      }
                    )
                  }
                )
              ])
          }
        )
      ]
    )
  }
}
const caseItemColumn = {
  title: "用例项",
  key: "items",
  render(row: CommonTestCaseItemData) {
    if (!row.testCaseItem.id) {
      return h(
        "div",
        {},
        {
          default: () => "暂无用例项"
        }
      )
    }
    return h(
      NGrid,
      {
        cols: 2
      },
      [
        h(NGridItem, null, {
          default: () =>
            h(
              NTooltip,
              {},
              {
                default: () => row.testCaseItem.remark,
                trigger: () => row.testCaseItem.content
              }
            )
        }),
        h(
          NGridItem,
          {
            class: "flex flex-row justify-end"
          },
          {
            default: () =>
              h(NButtonGroup, {}, [
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "编辑用例项",
                    trigger: h(
                      NButton,
                      {
                        size: "small",
                        type: "primary",
                        secondary: true,
                        disabled: !userStore.hasResourceCode("asset:commonTestCase:updateItem"),
                        onClick: () => {
                          currentItem.value = row.testCaseItem
                          itemDrawerActive.value = true
                        }
                      },
                      {
                        icon: h(NIcon, {}, { default: () => h(Edit) })
                      }
                    )
                  }
                ),
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "删除用例项",
                    trigger: h(
                      NPopconfirm,
                      {
                        onPositiveClick: () => handleDeleteItem(row.testCaseItem.id)
                      },
                      {
                        default: () => "确认删除该用例项吗？",
                        trigger: () =>
                          h(
                            NButton,
                            {
                              type: "error",
                              disabled: !userStore.hasResourceCode("asset:commonTestCase:deleteItem"),
                              size: "small"
                            },
                            { icon: h(NIcon, {}, { default: () => h(Trash) }) }
                          )
                      }
                    )
                  }
                )
              ])
          }
        )
      ]
    )
  }
}
const columns = [caseNameColumn, caseItemColumn]

const refresh = () => {
  loading.value = true
  getCommonTestCaseList(condition.value)
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
  pageSize: 5,
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

// 删除用例
const handleDeleteData = (id: string) => {
  deleteCommonTestCase(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}

// 删除用例项
const handleDeleteItem = (id: string) => {
  deleteCommonTestCaseItem(id).then((res) => {
    if (res) {
      message.success("删除成功")
      refresh()
    }
  })
}

// 用例抽屉
const caseDrawerActive = ref(false)
const currentCase = ref<CommonTestCase>({
  id: "",
  name: ""
})
const handleNewCase = () => {
  currentCase.value = {
    id: "",
    name: ""
  }
  caseDrawerActive.value = true
}

// 用例项抽屉
const itemDrawerActive = ref(false)
const currentItem = ref<CommonTestCaseItem>({
  id: "",
  testCaseId: "",
  content: ""
})
const handleNewItem = (testCaseId: string) => {
  currentItem.value = {
    id: "",
    testCaseId: testCaseId,
    content: ""
  }
  itemDrawerActive.value = true
}

// 界面加载
onMounted(() => {
  refresh()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-button @click="handleNewCase" v-resource-code="'asset:commonTestCase:add'">新增用例</n-button>
    </n-gi>
    <n-gi v-resource-code="'asset:commonTestCase:list'">
      <n-data-table
        size="small"
        remote
        :data="data"
        :loading="loading"
        :row-key="(row: CommonTestCaseItemData) => row.testCaseItem.id"
        :pagination="pagination"
        :on-update:page="handlePageChange"
        :on-update:page-size="handlePageSizeChange"
        :columns="columns"
        :single-line="false"
      />
    </n-gi>
  </n-grid>

  <case-drawer v-model:drawer-active="caseDrawerActive" v-model:data="currentCase" @refresh="refresh" />
  <item-drawer v-model:drawer-active="itemDrawerActive" v-model:data="currentItem" @refresh="refresh" />
</template>
