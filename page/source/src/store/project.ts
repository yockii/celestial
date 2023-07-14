import { Project, ProjectMember, ProjectModule } from "@/types/project"
import { defineStore } from "pinia"
import { useUserStore } from "./user"
import { useRoute } from "vue-router"

type ProjectModules = {
  [key: string]: ProjectModule[]
}
type projectResourceCodes = {
  [key: string]: string[]
}

export const useProjectStore = defineStore("project", {
  state: (): {
    projects: Project[]
    projectModules: ProjectModules
    projectResourceCodes: projectResourceCodes
  } => ({
    projects: [],
    projectModules: {},
    projectResourceCodes: {}
  }),
  getters: {
    project: (state): Project => {
      const route = useRoute()
      const projectId = route.params.id as string
      return (
        state.projects.find((project) => project.id === projectId) || {
          id: "",
          name: "",
          code: "",
          description: "",
          stageId: ""
        }
      )
    },
    modules: (state): ProjectModule[] => {
      const route = useRoute()
      const projectId = route.params.id as string
      if (state.projectModules[projectId]) {
        return state.projectModules[projectId]
      }
      return []
    },
    resourceCodes: (state): string[] => {
      const route = useRoute()
      const projectId = route.params.id as string
      if (state.projectResourceCodes[projectId]) {
        return state.projectResourceCodes[projectId]
      }
      return []
    },
    moduleTree(): ProjectModule[] {
      const tree = this.modules.filter((module) => !module.parentId && module.status !== 1 && module.status !== -1)
      const findChildren = (parent: ProjectModule) => {
        const children = this.modules.filter((module) => module.parentId === parent.id && module.status !== 1 && module.status !== -1)
        if (children.length) {
          parent.children = children
          children.forEach((child) => findChildren(child))
        }
      }
      tree.forEach((module) => findChildren(module))
      return tree
    },
    memberList(state): ProjectMember[] {
      const userIdSet = new Set<string>()
      const set = new Set<ProjectMember>()
      // 将members中的成员根据userId进行去重并添加到set中
      this.project.members?.forEach((member) => {
        if (!userIdSet.has(member.userId)) {
          userIdSet.add(member.userId)
          set.add(member)
        }
      })
      return Array.from(set)
    },
    isOwner(state) {
      return this.project.ownerId === useUserStore().user.id
    }
  },
  actions: {
    addProject(project: Project) {
      // 如果id相同则更新，否则新增
      const index = this.projects.findIndex((item) => item.id === project.id)
      if (index !== -1) {
        this.projects[index] = project
      } else {
        this.projects.push(project)
      }
    },
    setProjectModules(projectId: string, modules: ProjectModule[]) {
      this.projectModules[projectId] = modules
    },
    setProjectResourceCodes(projectId: string, resourceCodes: string[]) {
      this.projectResourceCodes[projectId] = resourceCodes
    },
    hasResourceCode(resourceCode: string) {
      const userStore = useUserStore()
      return (
        userStore.user.id === this.project.ownerId ||
        userStore.hasResourceCode("allProjectDetail") ||
        (this.resourceCodes && this.resourceCodes.includes(resourceCode))
      )
    },
    hasResourceCodes(resourceCodes: string[]) {
      const userStore = useUserStore()
      return (
        userStore.user.id === this.project.ownerId ||
        userStore.hasResourceCode("allProjectDetail") ||
        resourceCodes.some((resourceCode) => this.resourceCodes.includes(resourceCode))
      )
    }
  },
  persist: {
    key: "project",
    storage: sessionStorage
  }
})
