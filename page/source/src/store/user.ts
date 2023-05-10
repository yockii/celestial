import { defineStore } from "pinia"

export interface User {
  username: string
}

export const useUserStore = defineStore("user", {
  state: (): {
    user: User,
    token: string,
  } => ({
    user: {
      username: ""
    },
    token: ''
  }),
  getters: {
    username: (state) => state.user.username,
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
    storage: sessionStorage
  }
})
