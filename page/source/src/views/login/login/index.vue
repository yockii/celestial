<script setup lang="ts">
import { ref } from "vue"
import { FormInst, FormItemInst, FormRules, useMessage } from "naive-ui"
import { UserPlus } from "@vicons/tabler"
import { login, loginInDingTalk } from "@/service"
import { useUserStore } from "@/store/user"
import { useRouter } from "vue-router"
import { storeToRefs } from "pinia"
import { useAppStore } from "@/store/app"
import { getThirdSourcePublic } from "@/service/api/settings/thirdSource"
import dd from "dingtalk-jsapi"

const formRef = ref<FormInst | null>(null)
const rPasswordFormItemRef = ref<FormItemInst | null>(null)
const message = useMessage()
const appStore = useAppStore()
const loginInfo = ref<{
  username: string
  password: string
}>({
  username: "",
  password: ""
})
const handlePasswordInput = (): void => {
  if (loginInfo.value.password) {
    rPasswordFormItemRef.value?.validate({
      trigger: "password-input"
    })
  }
}
const rules: FormRules = {
  username: [
    {
      required: true,
      min: 3,
      trigger: ["input", "blur"],
      message: "请输入用户名"
    }
  ],
  password: [
    {
      required: true,
      min: 5,
      message: "请输入密码"
    }
  ]
}

const userStore = useUserStore()
const router = useRouter()
const handleLoginButtonClick = (): void => {
  formRef.value?.validate().then(() => {
    login(loginInfo.value.username, loginInfo.value.password)
      .then((data) => {
        if (data.token === "") {
          message.error("登录失败, 未获取到token")
          return
        }
        userStore.setToken(data.token)
        userStore.setUserInfo(data.user)
        message.success("登录成功")
        const to = {
          name: "Home"
        }
        if (appStore.activeSubMenuKey !== "") {
          to.name = appStore.activeSubMenuKey
        } else if (appStore.activeMenuKey !== "") {
          to.name = appStore.activeMenuKey
        }
        router.push(to)
      })
      .catch((err) => {
        message.error(err)
      })
  })
}
const emits = defineEmits(["register"])
const handleRegisterClick = (): void => {
  emits("register")
}

// 三方登录
const { thirdSourceList } = storeToRefs(useAppStore())
const handleDingtalkLogin = (id: string) => {
  const thirdSource = thirdSourceList.value.find((item) => item.id === id)
  if (thirdSource) {
    const query = encodeURI(
      `redirect_uri=${import.meta.env.VITE_REDIRECT_URI}/auth&response_type=code&scope=openid corpid&corpId=${thirdSource.corpId}&prompt=consent&state=${
        thirdSource.id
      }&client_id=${thirdSource.appKey}`
    )

    window.location.href = `https://login.dingtalk.com/oauth2/auth?${query}`
  }
}

// 加载初始化
onMounted(() => {
  getThirdSourcePublic().then((res) => {
    if (res) {
      thirdSourceList.value = res
      loginInDingtalk()
    }
  })
})

const dingtalkThirdSource = computed(() => {
  return thirdSourceList.value.find((item) => item.code === "dingtalk")
})
const loginInDingtalk = () => {
  if (dd.env.platform === "notInDingTalk") {
    console.warn("当前非钉钉环境，无法使用钉钉自动免登，将呈现登录页面")
  } else {
    if (dingtalkThirdSource.value) {
      // 尝试钉钉自动免登
      dd.ready(() => {
        dd.runtime.permission
          .requestAuthCode({
            corpId: import.meta.env.VITE_DD_CORPID
          })
          .then((info: { code: string }) => {
            dingtalkThirdSource.value &&
              loginInDingTalk(dingtalkThirdSource.value.id, info.code).then((data) => {
                if (data.token === "") {
                  message.error("登录失败, 未获取到token")
                  return
                }
                userStore.setToken(data.token)
                userStore.setUserInfo(data.user)
                message.success("登录成功")
                const to = {
                  name: "Home"
                }
                if (appStore.activeSubMenuKey !== "") {
                  to.name = appStore.activeSubMenuKey
                } else if (appStore.activeMenuKey !== "") {
                  to.name = appStore.activeMenuKey
                }
                router.push(to)
              })
          })
          .catch((e) => {
            console.warn(e)
          })
      })
    }
  }
}
</script>

<template>
  <n-form ref="formRef" :model="loginInfo" :rules="rules">
    <n-form-item path="username" label="用户名">
      <n-input v-model:value="loginInfo.username" placeholder="请输入用户名" @keydown.enter.prevent />
    </n-form-item>
    <n-form-item path="password" label="密码">
      <n-input v-model:value="loginInfo.password" type="password" placeholder="请输入密码" @input="handlePasswordInput" @keydown.enter.prevent />
    </n-form-item>
    <n-grid :cols="2">
      <n-grid-item>
        <n-button text @click="handleRegisterClick" v-if="false">
          <template #icon>
            <n-icon>
              <user-plus />
            </n-icon>
          </template>
          注册新用户
        </n-button>
      </n-grid-item>
      <n-grid-item>
        <div style="display: flex; justify-content: flex-end">
          <n-button :disabled="loginInfo.username === '' || loginInfo.password === ''" round type="primary" @click="handleLoginButtonClick"> 登录 </n-button>
        </div>
      </n-grid-item>
    </n-grid>
    <n-grid :cols="3">
      <n-gi v-for="item in thirdSourceList" :key="item.id">
        <n-button type="primary" @click="() => handleDingtalkLogin(item.id)">
          <n-icon>
            <user-plus />
          </n-icon>
          {{ item.name }}
        </n-button>
      </n-gi>
    </n-grid>
  </n-form>
</template>

<style scoped></style>
