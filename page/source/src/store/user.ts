import { defineStore } from "pinia"
import { User } from "@/types/user"
import { RouteHistory } from "@/types/app"
import { Project } from "@/types/project"

export const useUserStore = defineStore("user", {
  state: (): {
    user: User
    token: string
    isSuperAdmin: boolean
    resourceCodes: string[]
    dataPermission: number
    history: RouteHistory[]
    myProjectList: Project[]
  } => ({
    user: {
      id: "",
      username: "",
      status: 0
    },
    token: "",
    isSuperAdmin: false,
    resourceCodes: [],
    dataPermission: 0,
    history: [],
    myProjectList: []
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
    },
    addRoute(route: RouteHistory) {
      if (!route.url || route.url === "/" || route.url === "/login") return
      // 如果当前路由已经存在，删除之前的记录
      const index = this.history.findIndex((item) => item.url === route.url)
      if (index > -1) {
        this.history.splice(index, 1)
      }
      if (this.history.length >= 15) {
        // 超过15条记录，删除最早的一条
        this.history.shift()
      }
      // 添加新的记录
      this.history.push(route)
    }
  },
  persist: {
    key: "user",
    storage: sessionStorage
  }
})
