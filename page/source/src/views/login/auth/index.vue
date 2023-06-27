<script setup lang="ts">
import { loginByDingTalk } from "@/service/api"
import { useUserStore } from "@/store/user"
import { useAppStore } from "@/store/app"

const message = useMessage()
const appStore = useAppStore()
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
        if (localStorage.getItem("redirect")) {
          router.push(localStorage.getItem("redirect") as string)
          localStorage.removeItem("redirect")
          return
        }
        const to = {
          name: "Home"
        }
        if (appStore.activeSubMenuKey && appStore.activeSubMenuKey !== "") {
          to.name = appStore.activeSubMenuKey
        } else if (appStore.activeMenuKey && appStore.activeMenuKey !== "") {
          to.name = appStore.activeMenuKey
        }
        router.push(to)
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
