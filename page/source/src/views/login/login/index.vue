<script setup lang="ts">
import {ref} from "vue";
import {FormInst, FormItemInst, FormRules, useMessage} from "naive-ui";
import {UserPlus} from "@vicons/tabler";
import {login} from "@/service";
import {useUserStore} from "@/store/user";
import { useRouter} from "vue-router";
const formRef = ref<FormInst | null>(null)
const rPasswordFormItemRef = ref<FormItemInst | null>(null)
const message = useMessage()
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
            trigger: 'password-input'
        })
    }
}
const rules: FormRules = {
    username: [
        {
            required: true,
            min: 3,
            trigger: ['input', 'blur'],
            message: '请输入用户名'
        }
    ],
    password: [
        {
            required: true,
            min: 5,
            message: '请输入密码'
        }
    ]
}

const userStore = useUserStore()
const router = useRouter()
const handleLoginButtonClick = (): void => {
    formRef.value?.validate().then(() => {
        login(loginInfo.value.username, loginInfo.value.password).then((data) => {
            if (data.token === '') {
                message.error('登录失败, 未获取到token')
                return
            }
            userStore.setToken(data.token)
            userStore.setUserInfo(data.user)
            message.success('登录成功')
            router.push("/")
        }).catch((err) => {
            message.error(err)
        });
    })
}
const emits = defineEmits(['register'])
const handleRegisterClick = (): void => {
    emits('register')
}

</script>

<template>
    <n-form ref="formRef" :model="loginInfo" :rules="rules">
        <n-form-item path="username" label="用户名">
            <n-input v-model:value="loginInfo.username" placeholder="请输入用户名" @keydown.enter.prevent />
        </n-form-item>
        <n-form-item path="password" label="密码">
            <n-input
                    v-model:value="loginInfo.password"
                    type="password"
                    placeholder="请输入密码"
                    @input="handlePasswordInput"
                    @keydown.enter.prevent
            />
        </n-form-item>
        <n-grid :cols="2">
            <n-grid-item>
                <n-button text @click="handleRegisterClick">
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
                    <n-button
                            :disabled="loginInfo.username === '' || loginInfo.password === ''"
                            round
                            type="primary"
                            @click="handleLoginButtonClick"
                    >
                        登录
                    </n-button>
                </div>
            </n-grid-item>
        </n-grid>
    </n-form>
</template>

<style scoped>

</style>