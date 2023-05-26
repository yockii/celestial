import { Project, ProjectModule } from "@/types/project"
import { defineStore } from "pinia"

export const useProjectStore = defineStore("project", {
  state: (): {
    project: Project | null
    tab: string
    modules: ProjectModule[]
  } => ({
    project: null,
    tab: "项目总览",
    modules: []
  }),
  getters: {
    moduleTree(state) {
      const tree = state.modules.filter((module) => !module.parentId)
      const findChildren = (parent: ProjectModule) => {
        const children = state.modules.filter((module) => module.parentId === parent.id)
        if (children.length) {
          parent.children = children
          children.forEach((child) => findChildren(child))
        }
      }
      tree.forEach((module) => findChildren(module))
      return tree
    }
  },
  actions: {},
  persist: {
    key: "project",
    storage: sessionStorage
  }
})
