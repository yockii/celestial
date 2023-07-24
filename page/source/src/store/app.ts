import { ThirdSourcePublic } from "@/types/thirdSource"
import { defineStore } from "pinia"

export const useAppStore = defineStore("app", {
  state: (): {
    theme: string
    activeMenuKey: string
    activeSubMenuKey: string
    collapsed: boolean
    thirdSourceList: ThirdSourcePublic[]
    editorUrl: string
    lastRequestTime: number
  } => ({
    theme: "",
    activeMenuKey: "Home",
    activeSubMenuKey: "Dashboard",
    collapsed: false,
    thirdSourceList: [],
    editorUrl: "",
    lastRequestTime: 0
  }),
  getters: {},
  actions: {
    setTheme(theme: string) {
      this.theme = theme
    }
  },
  persist: true
})
