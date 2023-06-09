<script setup lang="ts">
import { getRoleList } from "@/service/api/settings/role"
import { assignRole } from "@/service/api/settings/user"
import { Role } from "@/types/user"
import { NTag } from "naive-ui"

const emit = defineEmits(["update:drawerActive"])
const message = useMessage()
const props = defineProps<{
  drawerActive: boolean
  userId: string
  userRoleIdList: string[]
}>()
watch(
  () => props.userRoleIdList,
  () => {
    resetCurrentData()
  }
)
const roles = ref<Role[]>([])
const roleIdList = ref<string[]>([])
const resetCurrentData = () => {
  roleIdList.value = props.userRoleIdList
}
const handleCommitData = () => {
  if (roleIdList.value.length === 0) {
    message.error("请选择角色")
    return
  }
  assignRole(props.userId, roleIdList.value).then((res) => {
    if (res) {
      message.success("更新成功")
      emit("update:drawerActive", false)
    }
  })
}
onMounted(() => {
  getRoleList({
    status: 1,
    offset: -1,
    limit: -1
  }).then((res) => {
    roles.value = res.items || []
  })
})

const renderLabel = (option: Role) => {
  const style: { textColor: string } = JSON.parse(option.style || "{}")
  return [
    h(
      "div",
      {
        style: {
          color: style.textColor
        }
      },
      option.name
    )
  ]
}
const renderTag = ({ option, handleClose }: { option: Role; handleClose: () => void }) => {
  return h(
    NTag,
    {
      closable: true,
      size: "small",
      color: JSON.parse(option.style || "{}"),
      onMousedown: (e: FocusEvent) => {
        e.preventDefault()
      },
      onClose: (e: MouseEvent) => {
        e.stopPropagation()
        handleClose()
      }
    },
    {
      default: () => option.name
    }
  )
}
</script>

<template>
  <n-drawer :show="drawerActive" :width="401" :on-update:show="(show: boolean) => $emit('update:drawerActive', show)">
    <n-drawer-content title="分配角色">
      <n-form label-width="100px">
        <n-form-item label="分配的角色" required>
          <n-select
            v-model:value="roleIdList"
            placeholder="请选择角色"
            :options="roles"
            :render-tag="renderTag"
            :render-label="renderLabel"
            label-field="name"
            value-field="id"
            multiple
            clearable
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-button class="mr-a" @click="resetCurrentData">重置</n-button>
        <n-button size="small" type="primary" @click="handleCommitData" v-resource-code="'system:user:dispatchRoles'">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>
