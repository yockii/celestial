<script setup lang="ts">
import { loginByDingTalk } from "@/service/api/settings/thirdSource"
import { useUserStore } from "@/store/user"

const message = useMessage()
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
onMounted(() => {
  const authCode = route.query.authCode
  const sourceId = route.query.state
  if (authCode && sourceId) {
    loginByDingTalk(sourceId as string, authCode as string)
      .then((data) => {
        if (data.token === "") {
          message.error("登录失败, 未获取到token")
          return
        }
        userStore.setToken(data.token)
        userStore.setUserInfo(data.user)
        message.success("登录成功")
        router.push("/")
      })
      .catch((err) => {
        message.error(err)
      })
  } else {
    // 如果没有获取到authCode，但userStore中有用户信息，则直接跳转到首页
    if (userStore.token !== "" && userStore.user.id !== "") {
      router.push("/")
      return
    }
    message.error("登录失败, 未获取到authCode")
  }
})
</script>

<template>
  <div>正在验证用户中，请稍后....</div>
</template>
