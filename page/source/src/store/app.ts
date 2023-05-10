import { defineStore } from "pinia"

export const useAppStore = defineStore("app", {
  state: () => ({
    theme: "",
    activeMenuKey: "Home",
    activeSubMenuKey: "Dashboard",
    collapsed: false,
  }),
  getters: {
  },
  actions: {
    setTheme(theme: string) {
      this.theme = theme
    }
  },
  persist: true
})
