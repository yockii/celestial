<script setup lang="ts">
import { getUserList } from "@/service/api/settings/user"
import { User, UserCondition } from "@/types/user"
import { Search } from "@vicons/carbon"
const emits = defineEmits(["userSelectedChanged"])
const props = defineProps<{
  selectedUserIds: string[]
}>()
const condition = ref<UserCondition>({})
const userList = ref<User[]>([])
const listUser = () => {
  condition.value.offset = 0
  condition.value.limit = 100
  getUserList(condition.value).then((res) => {
    userList.value = res.items || []
  })
}
// const appendUserList = () => {
//     condition.value.offset = userList.value.length
//     getUserList(condition.value).then(res => {
//         userList.value = userList.value.concat(res.items)
//     })
// }
const userChecked = computed(() => {
  return (userId: string) => props.selectedUserIds.includes(userId)
})
const checkUser = (userId: string, checked: boolean) => {
  const user = userList.value.find((item) => item.id === userId)
  emits("userSelectedChanged", user, checked)
}
</script>

<template>
  <div class="mb-16px">
    <n-input v-model:value="condition.realName" round placeholder="搜索姓名" @keydown.enter.prevent="listUser">
      <template #suffix>
        <n-icon @click="listUser" class="cursor-pointer">
          <Search />
        </n-icon>
      </template>
    </n-input>
  </div>
  <n-grid :cols="2" x-gap="8" y-gap="8">
    <n-gi v-for="user in userList" :key="user.id">
      <n-tag :checked="userChecked(user.id)" checkable :on-update:checked="(val:boolean) => checkUser(user.id, val)">{{ user.realName }}</n-tag>
    </n-gi>
  </n-grid>
</template>

<style scoped></style>
