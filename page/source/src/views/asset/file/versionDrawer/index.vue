<script setup lang="ts">
import { File, FileVersion, FileVersionCondition } from "@/types/asset"
import { NButtonGroup, NButton, NIcon, NTooltip, PaginationProps } from "naive-ui"
import NameAvatar from "@/components/NameAvatar.vue"
import { WordCloud } from "@vicons/carbon"
import { getAssetFileVersionList } from "@/service/api"
import dayjs from "dayjs"

const router = useRouter()
const emit = defineEmits(["update:drawerActive"])
const props = defineProps<{
  drawerActive: boolean
  data: File
}>()

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
    title: "更新时间",
    key: "createTime",
    render: (row: FileVersion) => {
      return dayjs(row.createTime).format("YYYY-MM-DD HH:mm:ss")
    }
  },
  {
    title: "更新人",
    key: "creator",
    render: (row: FileVersion) => {
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
    render: (row: FileVersion) => {
      const btnGroup: VNode[] = []
      if (props.data.permission && props.data.permission >= 3) {
        btnGroup.push(
          h(
            NButton,
            {
              size: "small",
              type: "tertiary",
              onClick: () => {
                const { href } = router.resolve({
                  name: "Editor",
                  params: {
                    id: props.data.id,
                    versionId: row.id
                  }
                })
                window.open(href, "_blank")
              }
            },
            {
              default: () => h(NIcon, { component: WordCloud })
            }
          )
        )
      }
      return h(NButtonGroup, {}, () => btnGroup)
    }
  }
]

const condition = ref<FileVersionCondition>({
  fileId: props.data.id
})
const list = ref<FileVersion[]>([])
const loading = ref(false)
const refresh = () => {
  loading.value = true
  getAssetFileVersionList(condition.value)
    .then((res) => {
      list.value = res.items || []
      pagination.itemCount = res.total
      pagination.pageCount = Math.ceil(res.total / res.limit)
      pagination.page = Math.ceil(res.offset / res.limit) + 1
    })
    .finally(() => {
      loading.value = false
    })
}

const handlePageChange = (page: number) => {
  condition.value.offset = (page - 1) * (condition.value.limit || 10)
  refresh()
}
const handlePageSizeChange = (pageSize: number) => {
  condition.value.limit = pageSize
  refresh()
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" @update:show="(show:boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="`${data.name}.${data.suffix}`">
      <n-data-table
        size="small"
        remote
        :data="list"
        :loading="loading"
        :columns="columns"
        :pagination="pagination"
        :row-key="(row: FileVersion) => row.id"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </n-drawer-content>
  </n-drawer>
</template>
