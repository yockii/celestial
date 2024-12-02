<script setup lang="ts">
import {StatusChange} from "@vicons/carbon";
import {NButton} from "naive-ui";
import {Role, User} from "@/types/user"
import {Project, ProjectMember} from "@/types/project"
import {changeCharger} from "@/service";

const props = defineProps<{
  project: Project
  members: ProjectMember[]
}>()
const showDrawer = ref(false);
const charger = ref<User | null>(null);
const currentChargerId = ref<string[]>([]);
const chargerSelected = (user:User, checked:boolean) => {
  if (checked) {
    // 替换
    currentChargerId.value = [user.id]
    charger.value = user
  } else {
    currentChargerId.value = []
    charger.value = null
  }
}
const resetCharger = () => {
  currentChargerId.value = []
}
const handleCommitCharger = () => {
  // 提交
  if (charger.value && currentChargerId.value.length > 0) {
    changeCharger(props.project.id, currentChargerId.value[0]).then((res) => {
      if(res){
        showDrawer.value = false
        const memberCharger = props.members.find((member) => member.roleId === '0')
        if (memberCharger) {
          memberCharger.userId = charger.value.id
          memberCharger.username = charger.value.username
          memberCharger.realName = charger.value.realName
        }
      }
    })
  }
}
</script>

<template>
  <n-tooltip>
    <template #trigger>
      <n-icon class="cursor-pointer" @click="showDrawer = true">
        <StatusChange />
      </n-icon>
    </template>
    更换负责人
  </n-tooltip>
  <n-drawer v-model:show="showDrawer" width="300">
    <n-drawer-content title="选择新的负责人">
      <user-select-drawer-content :selected-user-ids="currentChargerId" @user-selected-changed="chargerSelected" />
      <template #footer>
        <n-button class="mr-a" @click="resetCharger">重置</n-button>
        <n-button type="primary" @click="handleCommitCharger">提交</n-button>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped>

</style>