<script setup lang="ts">
import { FileUser } from "@/types/asset"
import { getAssetFileUserList, updateAssetFileUser, deleteAssetFileUser } from "@/service/api"
import { SearchOutlined, PlusOutlined } from "@vicons/material"
import { ChevronDown } from "@vicons/carbon"
import { User } from "@/types/user"
const message = useMessage()

const emit = defineEmits(["update:drawerActive"])
const props = defineProps<{
  drawerActive: boolean
  fileId: string
  fileName: string
  creatorId: string
}>()

const fileUsers = ref<FileUser[]>([])
const creator = computed(() => {
  if (!fileUsers.value) return null
  return fileUsers.value.find((item) => item.userId === props.creatorId)
})
const userWithoutCreator = computed(() => {
  if (!fileUsers.value) return []
  return fileUsers.value.filter((item) => item.userId !== props.creatorId)
})
const condition = ref<FileUser>({
  fileId: props.fileId
})
const loading = ref(false)
const refresh = () => {
  if (loading.value) return
  if (props.fileId) {
    loading.value = true
    getAssetFileUserList(condition.value)
      .then((res) => {
        fileUsers.value = res
      })
      .finally(() => {
        loading.value = false
      })
  }
}
const options = ref([
  { label: "只读", key: 1 },
  { label: "可编辑", key: 2 },
  { label: "可下载", key: 3 },
  { label: "可管理", key: 4 },
  { type: "divider" },
  { label: "移除权限", key: 0, props: { style: { color: "red" } } }
])
const valueToLabel = (value: number) => {
  return options.value.find((item) => item.key === value)?.label || "无权限"
}
const handleSelect = (userId: string, value: number) => {
  if (userId === "") {
    return
  }
  const fileUser = fileUsers.value.find((item) => item.userId === userId)
  if (fileUser && fileUser.permission !== value) {
    if (value > 0) {
      updateAssetFileUser({ fileId: props.fileId, userId, permission: value }).then((res) => {
        if (res) {
          message.success("更新成功")
          fileUser.permission = value
        }
      })
    } else {
      if (fileUser.id) {
        deleteAssetFileUser(fileUser.id).then((res) => {
          if (res) {
            message.success("更新成功")
            fileUsers.value = fileUsers.value.filter((item) => item.id !== fileUser.id)
          }
        })
      }
    }
  }
}

watch(
  () => props.fileId,
  (v) => {
    if (v) {
      condition.value.fileId = v
      refresh()
    }
  }
)

// 增加用户抽屉
const addUserDrawerActive = ref(false)
const showAddUser = () => {
  addUserDrawerActive.value = true
}
const selectedUserList = ref<User[]>([])
const selectedUserIdList = computed(() => {
  return selectedUserList.value.map((item) => item.id)
})
const userSelect = (user: User, checked: boolean) => {
  if (checked) {
    selectedUserList.value.push(user)
  } else {
    selectedUserList.value = selectedUserList.value.filter((item) => item.id !== user.id)
  }
}
const confirmUsers = () => {
  // fileUsers中没有的，就加入fileUsers中，默认为0，无权限
  const newUsers = selectedUserList.value.filter((user) => !fileUsers.value.find((item) => item.userId === user.id))
  if (newUsers.length) {
    const newFileUsers = newUsers.map((user) => {
      return {
        fileId: props.fileId,
        userId: user.id,
        realName: user.realName,
        permission: 0
      }
    })
    fileUsers.value = [...fileUsers.value, ...newFileUsers]
  }
  addUserDrawerActive.value = false
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="502" @update:show="(show:boolean) => emit('update:drawerActive', show)">
    <n-drawer-content :title="'配置' + fileName + '权限'">
      <n-space vertical :size="32">
        <n-space justify="space-between">
          <n-input v-model:value="condition.realName" placeholder="请输入姓名搜索" @keydown.enter="refresh()">
            <template #suffix>
              <n-icon class="cursor-pointer" @click="refresh()">
                <SearchOutlined />
              </n-icon>
            </template>
          </n-input>
          <n-tooltip>
            <template #trigger>
              <n-button type="primary" @click="showAddUser()">
                <template #icon>
                  <n-icon>
                    <PlusOutlined />
                  </n-icon>
                </template>
              </n-button>
            </template>
            新增用户
          </n-tooltip>
        </n-space>
        <n-space vertical :size="16">
          <n-space justify="space-between" v-if="creator">
            <n-text depth="3">{{ creator.realName }} <n-tag size="small" type="info">所有者</n-tag></n-text>
            <n-text depth="3">可管理</n-text>
          </n-space>
          <n-space v-for="pu in userWithoutCreator" :key="pu.id" justify="space-between">
            <n-text>{{ pu.realName }}</n-text>
            <n-dropdown trigger="click" :options="options" @select="(v:number) => handleSelect(pu.userId || '', v)">
              <n-text class="cursor-pointer"
                >{{ valueToLabel(pu.permission || 0) }}<n-icon><ChevronDown /></n-icon
              ></n-text>
            </n-dropdown>
          </n-space>
        </n-space>
      </n-space>
    </n-drawer-content>
    <!-- 选择人的抽屉 -->
    <n-drawer v-model:show="addUserDrawerActive" :width="401">
      <n-drawer-content title="选择人员">
        <user-select-drawer-content :selected-user-ids="selectedUserIdList" @user-selected-changed="userSelect" />
        <template #footer>
          <n-button type="primary" @click="confirmUsers">确定</n-button>
        </template>
      </n-drawer-content>
    </n-drawer>
  </n-drawer>
</template>
