<script setup lang="ts">
import {ref} from "vue";
import {FormInst, FormItemInst, FormItemRule, FormRules, useMessage} from "naive-ui";
import {ArrowBack} from "@vicons/tabler";
const formRef = ref<FormInst | null>(null)
const rPasswordFormItemRef = ref<FormItemInst | null>(null)
const message = useMessage()

const registerInfo = ref<{
    username: string
    password: string
    reenteredPassword: string
}>({
    username: "",
    password: "",
    reenteredPassword: ""
})
const validatePasswordStartWith = (
    rule: FormItemRule,
    value: string
): boolean => {
    return (
        !!registerInfo.value.password &&
        registerInfo.value.password.startsWith(value) &&
        registerInfo.value.password.length >= value.length
    )
}
const validatePasswordSame = (rule: FormItemRule, value: string): boolean => {
    return value === registerInfo.value.password
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
    ],
    reenteredPassword: [
        {
            required: true,
            message: '请再次输入密码',
            trigger: ['input', 'blur']
        },
        {
            validator: validatePasswordStartWith,
            message: '两次密码输入不一致',
            trigger: 'input'
        },
        {
            validator: validatePasswordSame,
            message: '两次密码输入不一致',
            trigger: ['blur', 'password-input']
        }
    ]
}
const handlePasswordInput = (): void => {
    if (registerInfo.value.password) {
        rPasswordFormItemRef.value?.validate({
            trigger: 'password-input'
        })
    }
}
const handleValidateButtonClick = (): void => {
    formRef.value?.validate().then(() => {
        message.success(             '验证成功' )
    })
}
const emits = defineEmits(['back'])
const handleBackLoginClick = (): void => {
    emits('back')
}
</script>

<template>
    <n-form ref="formRef" :model="registerInfo" :rules="rules">
        <n-form-item path="username" label="用户名">
            <n-input v-model:value="registerInfo.username" placeholder="请输入用户名" @keydown.enter.prevent />
        </n-form-item>
        <n-form-item path="password" label="密码">
            <n-input
                    v-model:value="registerInfo.password"
                    type="password"
                    placeholder="请输入密码"
                    @input="handlePasswordInput"
                    @keydown.enter.prevent
            />
        </n-form-item>
        <n-form-item
                ref="rPasswordFormItemRef"
                first
                path="reenteredPassword"
                label="重复密码"
        >
            <n-input
                    v-model:value="registerInfo.reenteredPassword"
                    :disabled="!registerInfo.password"
                    type="password"
                    placeholder="请再次输入密码"
                    @keydown.enter.prevent
            />
        </n-form-item>
        <n-grid :cols="2">
            <n-grid-item>
                <n-button text @click="handleBackLoginClick">
                    <template #icon>
                        <n-icon>
                            <arrow-back />
                        </n-icon>
                    </template>
                    返回登录
                </n-button>
            </n-grid-item>
            <n-grid-item>
                <div style="display: flex; justify-content: flex-end">
                    <n-button
                            :disabled="registerInfo.reenteredPassword === ''"
                            round
                            type="primary"
                            @click="handleValidateButtonClick"
                    >
                        注册
                    </n-button>
                </div>
            </n-grid-item>
        </n-grid>
    </n-form>
</template>

<style scoped>

</style>