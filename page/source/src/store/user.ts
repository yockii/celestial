import { defineStore } from "pinia"
import { User } from "@/types/user"
import { RouteHistory } from "@/types/app"
import { Project } from "@/types/project"

// 后端 userTokenExpire 为 86400 秒，cookie 有效期与之保持一致
const TOKEN_COOKIE_NAME = "token"
const TOKEN_COOKIE_MAX_AGE = 86400

// 登录态同步写入 cookie，供 <img>/<a> 等无法携带 Authorization 头的请求使用
// （如富文本中的文件下载链接），后端 TokenLookup 同时支持 header 与 cookie
export function setTokenCookie(token: string) {
  document.cookie = `${TOKEN_COOKIE_NAME}=${encodeURIComponent(token)}; path=/; max-age=${TOKEN_COOKIE_MAX_AGE}; SameSite=Lax`
}

export function clearTokenCookie() {
  document.cookie = `${TOKEN_COOKIE_NAME}=; path=/; max-age=0; SameSite=Lax`
}

export function hasTokenCookie() {
  return document.cookie.split("; ").some((c) => c.startsWith(TOKEN_COOKIE_NAME + "="))
}

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
      setTokenCookie(token)
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
      clearTokenCookie()
    },
    hasResourceCode(resourceCode: string) {
      return this.isSuperAdmin || this.resourceCodes.includes(resourceCode)
    },
    addRoute(route: RouteHistory) {
      if (!route.url || route.url === "/" || route.url === "/login" || route.url.startsWith("/auth")) return
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
