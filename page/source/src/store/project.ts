import { Project, ProjectMember, ProjectModule } from "@/types/project"
import { defineStore } from "pinia"

export const useProjectStore = defineStore("project", {
  state: (): {
    project: Project
    modules: ProjectModule[]
  } => ({
    project: {
      id: "",
      name: "",
      code: "",
      description: "",
      stageId: ""
    },
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
    },
    memberList(state) {
      const userIdSet = new Set<string>()
      const set = new Set<ProjectMember>()
      // 将members中的成员根据userId进行去重并添加到set中
      state.project.members?.forEach((member) => {
        if (!userIdSet.has(member.userId)) {
          userIdSet.add(member.userId)
          set.add(member)
        }
      })
      return Array.from(set)
    }
  },
  actions: {},
  persist: {
    key: "project",
    storage: sessionStorage
  }
})
