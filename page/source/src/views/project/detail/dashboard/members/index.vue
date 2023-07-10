<script setup lang="ts">
import { addProjectMembers } from "@/service/api"
import { Role, User } from "@/types/user"
import { Project, ProjectMember } from "@/types/project"
import { computed } from "vue"
import { UserFollow, IdManagement } from "@vicons/carbon"
import { rgbToHex } from "@/utils/Render"
import UserSelectDrawerContent from "@/components/user/UserSelectDrawerContent.vue"
import { NButton } from "naive-ui"

const emits = defineEmits(["projectMemberChanged"])
const props = defineProps<{
  project: Project
  members: ProjectMember[]
  roles: Role[]
}>()

const charger = computed(() => {
  // 负责人，找出角色ID为0的
  return props.members.find((m) => m.roleId === "0")
})
const noChargerMembers = computed(() => {
  // 非负责人
  return props.members.filter((m) => m.roleId !== "0")
})
const roleColor = computed(() => {
  return (roleId: string) => {
    const role = props.roles.find((r) => r.id === roleId)
    return role?.style ? JSON.parse(role.style) : undefined
  }
})
const roleName = computed(() => {
  return (roleId: string) => {
    const role = props.roles.find((r) => r.id === roleId)
    return role?.name
  }
})

// 抽屉
const manageActive = ref<boolean>(false)
const roleTextColor = (bgColor: string) => {
  // 计算灰度
  const gray = parseInt(bgColor.slice(1, 3), 16) * 0.299 + parseInt(bgColor.slice(3, 5), 16) * 0.587 + parseInt(bgColor.slice(5, 7), 16) * 0.114
  // 根据灰度，给出前景色
  return gray > 192 ? "#000000" : "#ffffff"
}
const roleTagStyle = computed(() => {
  return (role: Role) => {
    let result = {
      backgroundColor: "#fff",
      color: "#000",
      border: ""
    }
    if (role?.style) {
      const color = JSON.parse(role.style).color
      if (color) {
        const rgbaAttr = color.match(/[\d.]+/g)
        if (rgbaAttr.length >= 3) {
          const hexedColor = rgbToHex(rgbaAttr[0], rgbaAttr[1], rgbaAttr[2])
          result.backgroundColor = hexedColor
          result.color = roleTextColor(hexedColor)
        }
      }
    }
    // 判断是否当前选中的角色
    if (role.id === selectedRoleId.value) {
      result["border"] = "3px dashed " + result.color
    }
    return result
  }
})
const selectedRoleId = ref<string>("")
const memberIdsInRole = ref<string[]>([])
const roleSelected = (roleId: string) => {
  selectedRoleId.value = roleId
  memberIdsInRole.value = props.members.filter((m) => m.roleId === roleId).map((m) => m.userId)
}
const roleMembers = computed(() => {
  return props.members.filter((m) => m.roleId === selectedRoleId.value)
})
const selectedRoleStyle = computed(() => {
  const currentRole = props.roles.find((r) => r.id === selectedRoleId.value)
  if (currentRole && currentRole.style) {
    return JSON.parse(currentRole.style)
  }
  return undefined
})

// 人员选择
const userActive = ref<boolean>(false)
const roleUserSelect = (user: User, checked: boolean) => {
  if (checked) {
    // 添加
    memberIdsInRole.value.push(user.id)
  } else {
    // 删除
    const index = memberIdsInRole.value.findIndex((m) => m === user.id)
    if (index !== -1) {
      memberIdsInRole.value.splice(index, 1)
    }
  }
}
const resetRoleUsers = () => {
  memberIdsInRole.value = props.members.filter((m) => m.roleId === selectedRoleId.value).map((m) => m.userId)
}
const handleCommitRoleUsers = () => {
  // 提交服务器
  addProjectMembers(props.project.id, [selectedRoleId.value], memberIdsInRole.value).then((res) => {
    if (res) {
      // 成功后，关闭抽屉
      userActive.value = false
      emits("projectMemberChanged")
    }
  })
}

// 加载处理
onMounted(() => {
  selectedRoleId.value = props.roles[0]?.id || ""
})
</script>

<template>
  <n-grid :cols="2" class="mb-16px">
    <n-gi class="text-1.2em op-90 font-500">项目组</n-gi>
    <n-gi class="flex flex-justify-end">
      <n-tooltip v-project-resource-code="['project:detail:member:add', 'project:detail:member:delete']">
        <template #trigger>
          <n-button size="small" @click="manageActive = true">
            <template #icon>
              <n-icon>
                <IdManagement />
              </n-icon>
            </template>
          </n-button>
        </template>
        管理项目组成员
      </n-tooltip>
    </n-gi>
  </n-grid>
  <n-grid :cols="5">
    <n-gi>
      <n-text tag="div" class="mb-4px">项目负责人： </n-text>
      <n-text>{{ charger?.realName }}</n-text>
    </n-gi>
    <n-gi :span="4">
      <n-text tag="div" class="mb-4px">项目组成员： </n-text>
      <n-text v-if="noChargerMembers.length === 0">无</n-text>
      <n-space v-else>
        <n-tooltip trigger="hover" v-for="member in noChargerMembers" :key="member.userId">
          <template #trigger>
            <n-tag :color="roleColor(member.roleId)">
              {{ member.realName }}
            </n-tag>
          </template>
          {{ roleName(member.roleId) }}
        </n-tooltip>
      </n-space>
    </n-gi>
  </n-grid>

  <!-- 项目成员管理抽屉 -->
  <n-drawer v-model:show="manageActive" :width="401">
    <n-drawer-content>
      <!-- 角色标签探出 -->
      <div class="absolute w-80px left--80px top-60px">
        <div v-for="role in roles" :key="role.id" class="my-8px pa-8px text-right cursor-pointer" :style="roleTagStyle(role)" @click="roleSelected(role.id)">
          {{ role.name }}
        </div>
      </div>
      <div>
        <n-button v-show="selectedRoleId !== ''" @click="userActive = true">
          <n-icon>
            <UserFollow />
          </n-icon>
        </n-button>
      </div>
      <n-space class="mt-8px">
        <n-tag v-for="member in roleMembers" :key="member.userId" :color="selectedRoleStyle">
          {{ member.realName }}
        </n-tag>
      </n-space>
    </n-drawer-content>
    <!-- 选择人的抽屉 -->
    <n-drawer v-model:show="userActive" :width="300">
      <n-drawer-content title="选择人员">
        <user-select-drawer-content :selected-user-ids="memberIdsInRole" @user-selected-changed="roleUserSelect" />
        <template #footer>
          <n-button class="mr-a" @click="resetRoleUsers">重置</n-button>
          <n-button type="primary" @click="handleCommitRoleUsers">提交</n-button>
        </template>
      </n-drawer-content>
    </n-drawer>
  </n-drawer>
</template>

<style scoped></style>
