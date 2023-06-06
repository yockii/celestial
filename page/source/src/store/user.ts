import { defineStore } from "pinia"
import { User } from "@/types/user"

export const useUserStore = defineStore("user", {
  state: (): {
    user: User
    token: string
    isSuperAdmin: boolean
    resourceCodes: string[]
    dataPermission: number
  } => ({
    user: {
      id: "",
      username: "",
      status: 0
    },
    token: "",
    isSuperAdmin: false,
    resourceCodes: [],
    dataPermission: 0
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
    },
    logout() {
      this.user = {
        id: "",
        username: "",
        status: 0
      }
      this.token = ""
      this.isSuperAdmin = false
      this.resourceCodes = []
      this.dataPermission = 0
    },
    hasResourceCode(resourceCode: string) {
      return this.isSuperAdmin || this.resourceCodes.includes(resourceCode)
    }
  },
  persist: {
    key: "user",
    storage: localStorage // sessionStorage
  }
})
