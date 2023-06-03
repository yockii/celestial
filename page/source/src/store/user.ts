import { defineStore } from "pinia"
import { User } from "@/types/user"

export const useUserStore = defineStore("user", {
  state: (): {
    user: User
    token: string
  } => ({
    user: {
      id: "",
      username: "",
      status: 0
    },
    token: ""
  }),
  getters: {
    username: (state) => state.user.username,
    realName: (state) => state.user.realName
  },
  actions: {
    setUserInfo(user: User) {
      this.user = user
    },
    setToken(token: string) {
      this.token = token
    }
  },
  persist: {
    key: "user",
    storage: localStorage // sessionStorage
  }
})
