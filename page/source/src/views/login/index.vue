<template>
  <!-- <div class="h-screen w-screen relative flex justify-center items-center"> -->
  <n-layout>
    <n-layout-content class="h-screen w-screen relative">
      <Transition
        :duration="{ enter: 500, leave: 500 }"
        mode="out-in"
        enter-active-class="animate__animated animate__bounceInLeft"
        leave-active-class="animate__animated animate__bounceOutRight"
      >
        <n-card v-if="isLogin" hoverable title="登录" class="z-2 w-300px h-320px rounded-1 shadow-sm ma-auto top-0 left-0 right-0 bottom-0 absolute">
          <login-form @register="isLogin = false" />
        </n-card>
        <n-card v-else hoverable title="注册" class="z-2 w-300px h-320px rounded-1 shadow-sm ma-auto top-0 left-0 right-0 bottom-0 absolute">
          <register-form @back="isLogin = true" />
        </n-card>
      </Transition>
    </n-layout-content>
  </n-layout>
  <!-- </div> -->
</template>

<script lang="ts" setup>
import LoginForm from "./login/index.vue"
import RegisterForm from "./register/index.vue"
import { ref } from "vue"
const isLogin = ref(true)
</script>
