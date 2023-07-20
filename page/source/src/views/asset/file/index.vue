<script setup lang="ts">
import { AssetCategory, File, FileCondition } from "@/types/asset"
import { downloadAssetFile, getAssetCategoryList, getAssetFileList } from "@/service/api"
import dayjs from "dayjs"
import { NButton, NButtonGroup, NIcon, NTooltip, PaginationProps } from "naive-ui"
import NameAvatar from "@/components/NameAvatar.vue"
import { useUserStore } from "@/store/user"
import Drawer from "./drawer/index.vue"
import { Download, Edit, LicenseThirdParty, Version, WordCloud } from "@vicons/carbon"
import PermissionDrawer from "@/components/asset/permissionDrawer.vue"
import VersionDrawer from "./versionDrawer/index.vue"

const router = useRouter()
const userStore = useUserStore()
const assetCategoryList = ref<AssetCategory[]>([])
const treeSelected = (keys: string[]) => {
  if (keys.length > 0) {
    condition.value.categoryId = keys[0]
  } else {
    condition.value.categoryId = ""
  }
  refresh()
}
const handleLoad = (node: AssetCategory) => {
  return new Promise<void>((resolve) => {
    getAssetCategoryList({
      parentId: node.id,
      offset: -1,
      limit: -1
    }).then((res) => {
      node.children = res.items.map((item) => {
        item.isLeaf = !item.childrenCount || item.childrenCount === 0
        return item
      })
      resolve()
    })
  })
}

const condition = ref<FileCondition>({})
const list = ref<File[]>([])
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
const columns = [
  {
    title: "文件名",
    key: "name"
  },
  {
    title: "后缀",
    key: "suffix"
  },
  {
    title: "大小",
    key: "size",
    render: (row: File) => {
      const size = row.size || 0
      if (size < 1024) {
        return size + "B"
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + "KB"
      } else if (size < 1024 * 1024 * 1024) {
        return (size / (1024 * 1024)).toFixed(2) + "MB"
      } else {
        return (size / (1024 * 1024 * 1024)).toFixed(2) + "GB"
      }
    }
  },
  {
    title: "上传时间",
    key: "createTime",
    render: (row: File) => {
      return dayjs(row.createTime).format("YYYY-MM-DD HH:mm:ss")
    }
  },
  {
    title: "上传人",
    key: "creatorUsername",
    render: (row: File) => {
      return h(
        NTooltip,
        {},
        {
          trigger: () =>
            h(NameAvatar, {
              name: row.creator?.realName || "未知"
            }),
          default: () => row.creator?.realName || "未知"
        }
      )
    }
  },
  {
    title: "操作",
    key: "action",
    render: (row: File) => {
      const btnGroup: VNode[] = []
      if (userStore.hasResourceCode("asset:file:download") && row.permission && row.permission >= 3) {
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
                      downloadAssetFile(row.id).then((res: Blob) => {
                        const reader = new FileReader()
                        reader.readAsDataURL(res)
                        reader.onload = (e) => {
                          const a = document.createElement("a")
                          a.download = row.name + "." + row.suffix
                          a.href = e.target?.result as string
                          document.body.appendChild(a)
                          a.click()
                          document.body.removeChild(a)
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

      if (row.permission && row.permission >= 1) {
        // 可查看版本列表
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "版本列表",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    type: "primary",
                    onClick: () => {
                      handleVersionList(row)
                    }
                  },
                  {
                    default: () => h(NIcon, { component: Version })
                  }
                )
            }
          )
        )
      }

      if (userStore.hasResourceCode("asset:file:update") && row.permission && row.permission >= 2) {
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "修改信息",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    type: "info",
                    onClick: () => {
                      handleEditFile(row)
                    }
                  },
                  {
                    default: () => h(NIcon, { component: Edit })
                  }
                )
            }
          )
        )
      }

      // 分配权限
      if (row.permission && row.permission === 4) {
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "分配权限",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    type: "warning",
                    onClick: () => {
                      handleAssignPermission(row)
                    }
                  },
                  {
                    default: () => h(NIcon, { component: LicenseThirdParty })
                  }
                )
            }
          )
        )
      }

      if (row.permission && row.permission >= 3) {
        btnGroup.push(
          h(
            NTooltip,
            {},
            {
              default: () => "在线编辑",
              trigger: () =>
                h(
                  NButton,
                  {
                    size: "small",
                    type: "tertiary",
                    onClick: () => {
                      const { href } = router.resolve({
                        path: `/editor/${row.id}`
                      })
                      window.open(href, "_blank")
                    }
                  },
                  {
                    default: () => h(NIcon, { component: WordCloud })
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

const refresh = () => {
  loading.value = true
  getAssetFileList(condition.value)
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

// 版本抽屉
const versionDrawerActive = ref(false)

// 抽屉
const drawerActive = ref(false)
const currentData = ref<File>({
  id: "",
  name: "",
  categoryId: ""
})
const handleAddFile = () => {
  currentData.value = {
    id: "",
    name: "",
    categoryId: condition.value.categoryId || ""
  }
  drawerActive.value = true
}
const handleEditFile = (row: File) => {
  currentData.value = Object.assign({}, row)
  drawerActive.value = true
}
const handleVersionList = (row: File) => {
  currentData.value = Object.assign({}, row)
  versionDrawerActive.value = true
}
const handleAssignPermission = (row: File) => {
  currentData.value = Object.assign({}, row)
  permissionDrawerActive.value = true
}
const permissionDrawerActive = ref(false)

onMounted(() => {
  getAssetCategoryList({
    onlyParent: true,
    offset: -1,
    limit: -1
  }).then((res) => {
    assetCategoryList.value = res.items.map((item) => {
      item.isLeaf = !item.childrenCount || item.childrenCount === 0
      return item
    })
  })
  refresh()
})

const route = useRoute()
onBeforeUpdate(() => {
  if (route.query.id) {
    condition.value.id = route.query.id as string
    refresh()
  }
})
</script>

<template>
  <n-grid :cols="6" x-gap="16">
    <n-gi>
      <!--资产目录树-->
      <n-tree
        block-line
        block-node
        cancelable
        key-field="id"
        label-field="name"
        children-field="children"
        :on-update:selected-keys="treeSelected"
        :data="assetCategoryList"
        :on-load="handleLoad"
      />
    </n-gi>
    <n-gi :span="5">
      <n-grid :cols="1" y-gap="16">
        <n-gi>
          <n-space justify="space-between">
            <span></span>
            <n-button type="primary" @click="handleAddFile" v-resource-code="'asset:file:add'">新增文件</n-button>
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
            :row-key="(row: File) => row.id"
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </n-gi>
      </n-grid>
    </n-gi>
  </n-grid>

  <drawer v-if="drawerActive" v-model:drawer-active="drawerActive" v-model:data="currentData" @refresh="refresh" />
  <permission-drawer
    v-if="permissionDrawerActive"
    v-model:drawer-active="permissionDrawerActive"
    :fileId="currentData.id"
    :fileName="currentData.name"
    :creatorId="currentData.creatorId || ''"
  />
  <version-drawer v-if="versionDrawerActive" v-model:drawer-active="versionDrawerActive" :data="currentData" />
</template>
