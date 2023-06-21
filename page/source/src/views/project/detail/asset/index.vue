<script setup lang="ts">
import { CheckmarkCircle } from "@vicons/ionicons5"
import { useUserStore } from "@/store/user"
import { ProjectAsset, ProjectAssetCondition } from "@/types/project"
import dayjs from "dayjs"
import { DataTableBaseColumn, DataTableFilterState, NButton, NButtonGroup, NIcon, NPopconfirm, NTooltip, PaginationProps } from "naive-ui"
import Drawer from "./drawer/index.vue"
import { useProjectStore } from "@/store/project"
import { storeToRefs } from "pinia"
import { downloadAssetFile, getAssetFile, getProjectAssetList } from "@/service/api"
import { Delete, Download, Edit } from "@vicons/carbon"
const userStore = useUserStore()
const projectStore = useProjectStore()
const { project } = storeToRefs(projectStore)

const types = ref<{ key: number; label: string }[]>([
  // 1-需求 2-设计 3-代码 4-测试 9-其他
  {
    key: 1,
    label: "需求"
  },
  {
    key: 2,
    label: "设计"
  },
  {
    key: 3,
    label: "代码"
  },
  {
    key: 4,
    label: "测试"
  },
  {
    key: 9,
    label: "其他"
  }
])
const tagSelected = (type: number) => {
  condition.value.type = condition.value.type === type ? 0 : type
  refresh()
}

const condition = ref<ProjectAssetCondition>({
  projectId: "",
  type: 0,
  status: 0
})
const list = ref<ProjectAsset[]>([])
const loading = ref(false)
const pagination = reactive({
  itemCount: 0,
  page: 1,
  pageCount: 1,
  pageSize: 10,
  prefix({ itemCount }: PaginationProps) {
    return `共${itemCount}条`
  }
})

const statusColumn = reactive({
  title: "状态",
  key: "status",
  render: (row: ProjectAsset) => {
    // 1-草稿 2-已审核 3-已发布 9-已归档
    switch (row.status) {
      case 1:
        return "草稿"
      case 2:
        return "已审核"
      case 3:
        return "已发布"
      case 9:
        return "已归档"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "草稿",
      value: 1
    },
    {
      label: "已审核",
      value: 2
    },
    {
      label: "已发布",
      value: 3
    },
    {
      label: "已归档",
      value: 9
    }
  ]
})
const typeColumn = reactive({
  title: "类型",
  key: "type",
  render: (row: ProjectAsset) => {
    // 1-需求 2-设计 3-代码 4-测试 9-其他
    switch (row.type) {
      case 1:
        return "需求"
      case 2:
        return "设计"
      case 3:
        return "代码"
      case 4:
        return "测试"
      case 9:
        return "其他"
      default:
        return "未知"
    }
  },
  filter: true,
  filterMultiple: false,
  filterOptionValue: 0,
  filterOptions: [
    {
      label: "需求",
      value: 1
    },
    {
      label: "设计",
      value: 2
    },
    {
      label: "代码",
      value: 3
    },
    {
      label: "测试",
      value: 4
    },
    {
      label: "其他",
      value: 9
    }
  ]
})
const columns = [
  {
    title: "资产",
    key: "name"
  },
  typeColumn,
  statusColumn,
  {
    title: "创建时间",
    key: "createTime",
    render: (row: ProjectAsset) => {
      return dayjs(row.createTime).format("YYYY-MM-DD HH:mm:ss")
    }
  },
  {
    title: "操作",
    key: "action",
    render: (row: ProjectAsset) => {
      const btnGroup: VNode[] = []
      if (userStore.hasResourceCode("asset:file:download")) {
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "下载",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    onClick: () => {
                      getAssetFile(row.fileId).then((file) => {
                        if (file) {
                          downloadAssetFile(file.id).then((res: Blob) => {
                            const reader = new FileReader()
                            reader.readAsDataURL(res)
                            reader.onload = (e) => {
                              const a = document.createElement("a")
                              a.download = row.name + "." + file.suffix
                              a.href = e.target?.result as string
                              document.body.appendChild(a)
                              a.click()
                              document.body.removeChild(a)
                            }
                          })
                        }
                      })
                    }
                  },
                  {
                    default: () => h(NIcon, { component: Download })
                  }
                )
            }
          )
        )
      }
      if (userStore.hasResourceCode("project:detail:asset:edit")) {
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "编辑",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    secondary: true,
                    type: "primary",
                    disabled: !userStore.hasResourceCode("project:detail:asset:edit"),
                    onClick: () => handleEditData(row)
                  },
                  {
                    default: () => h(NIcon, { component: Edit })
                  }
                )
            }
          )
        )
      }
      if (userStore.hasResourceCode("project:detail:asset:delete")) {
        btnGroup.push(
          h(
            NPopconfirm,
            {
              onPositiveClick: () => handleDeleteData(row.id)
            },
            {
              default: () => "确认删除",
              trigger: () =>
                h(
                  NTooltip,
                  {},
                  {
                    default: () => "删除",
                    trigger: () =>
                      h(
                        NButton,
                        {
                          size: "small",
                          disabled: !userStore.hasResourceCode("project:detail:asset:delete"),
                          type: "error"
                        },
                        {
                          default: () => h(NIcon, { component: Delete })
                        }
                      )
                  }
                )
            }
          )
        )
      }

      return h(NButtonGroup, {}, () => btnGroup)
    }
  }
]

const handlePageChange = (page: number) => {
  condition.value.offset = (page - 1) * (condition.value.limit || 10)
  refresh()
}
const handlePageSizeChange = (pageSize: number) => {
  condition.value.limit = pageSize
  refresh()
}
const handleFiltersChange = (filters: DataTableFilterState, sourceColumn: DataTableBaseColumn) => {
  if (!loading.value) {
    switch (sourceColumn.key) {
      case "status":
        if (filters["status"] instanceof Array) {
          statusColumn.filterOptionValue = (filters["status"][0] as number) || 0
        } else {
          statusColumn.filterOptionValue = filters["status"] as number
        }
        condition.value.status = statusColumn.filterOptionValue
        break
      case "type":
        if (filters["type"] instanceof Array) {
          typeColumn.filterOptionValue = (filters["type"][0] as number) || 0
        } else {
          typeColumn.filterOptionValue = filters["type"] as number
        }
        condition.value.type = typeColumn.filterOptionValue
        break
    }

    refresh()
  }
}

const refresh = () => {
  loading.value = true
  getProjectAssetList(condition.value)
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

// 增改资产
const drawerActive = ref(false)
const currentData = ref<ProjectAsset>({
  id: "",
  name: "",
  fileId: "",
  type: 1,
  status: 1,
  projectId: project.value.id
})
const handleAddProjectAsset = () => {
  currentData.value = {
    id: "",
    name: "",
    fileId: "",
    type: 1,
    status: 1,
    projectId: project.value.id
  }
  drawerActive.value = true
}
const handleEditData = (row: ProjectAsset) => {
  currentData.value = Object.assign({}, row)
  drawerActive.value = true
}

// 删除资产
const handleDeleteData = (id: string) => {
  console.log("handleDeleteData", id)
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-grid :cols="1" y-gap="16">
    <n-gi>
      <n-space justify="space-between">
        <n-space>
          <n-tag v-for="t in types" :key="t.key" @click="tagSelected(t.key)" class="cursor-pointer" :type="t.key === condition.type ? 'primary' : ''">
            <template v-if="t.key === condition.type" #icon>
              <n-icon :component="CheckmarkCircle" />
            </template>
            {{ t.label }}
          </n-tag>
        </n-space>
        <n-button type="primary" size="small" @click="handleAddProjectAsset">新增资产</n-button>
      </n-space>
    </n-gi>
    <n-gi>
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :columns="columns"
        :pagination="pagination"
        :row-key="(row: ProjectAsset) => row.id"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        @update:filters="handleFiltersChange"
      />
    </n-gi>
  </n-grid>

  <drawer v-model:drawer-active="drawerActive" v-model:data="currentData" @refresh="refresh" />
</template>
