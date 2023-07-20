import { createRouter, createWebHistory, isNavigationFailure, RouteRecordRaw } from "vue-router"
import { useUserStore } from "@/store/user"
import { createDiscreteApi } from "naive-ui"
import { useProjectStore } from "@/store/project"
import { useAppStore } from "@/store/app"
const { loadingBar } = createDiscreteApi(["loadingBar"])

const routes: Array<RouteRecordRaw> = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/index.vue")
  },
  {
    path: "/auth",
    name: "Auth",
    component: () => import("@/views/login/auth/index.vue")
  },
  {
    path: "/index",
    name: "Home",
    alias: "/",
    component: () => import("@/layout/Primary.vue"),
    redirect: { name: "Dashboard" },
    meta: {
      label: () => "首页",
      activeMenuKey: "Home",
      activeSubMenuKey: "Dashboard"
    },
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/index/dashboard/index.vue"),
        meta: {
          label: () => "首页",
          activeMenuKey: "Home",
          activeSubMenuKey: "Dashboard"
        }
      }
    ]
  },
  {
    path: "/project",
    name: "Project",
    component: () => import("@/layout/Primary.vue"),
    meta: {
      label: () => "项目列表",
      activeMenuKey: "ProjectList",
      activeSubMenuKey: "ProjectList"
    },
    children: [
      {
        path: "list",
        name: "ProjectList",
        component: () => import("@/views/project/list/index.vue"),
        meta: {
          label: () => "项目列表",
          activeMenuKey: "ProjectList",
          activeSubMenuKey: "ProjectList"
        }
      },
      {
        path: "detail/:id",
        name: "ProjectDetail",
        redirect: {
          name: "ProjectDashboard"
        },
        component: () => import("@/views/project/detail/index.vue"),
        meta: {
          label: () => `【${useProjectStore().project.name}】总览`,
          activeMenuKey: "ProjectList",
          activeSubMenuKey: "ProjectList"
        },
        children: [
          {
            path: "dashboard",
            name: "ProjectDashboard",
            meta: {
              title: "项目总览",
              label: () => `【${useProjectStore().project.name}】总览`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/dashboard/index.vue")
          },
          {
            path: "plan",
            name: "ProjectPlan",
            meta: {
              title: "项目计划",
              label: () => `【${useProjectStore().project.name}】计划`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/plan/index.vue")
          },
          {
            path: "module",
            name: "ProjectModule",
            meta: {
              title: "功能模块",
              label: () => `【${useProjectStore().project.name}】功能模块`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/module/index.vue")
          },
          {
            path: "requirement",
            name: "ProjectRequirement",
            meta: {
              title: "项目需求",
              label: () => `【${useProjectStore().project.name}】需求`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/requirement/index.vue")
          },
          {
            path: "task",
            name: "ProjectTask",
            meta: {
              title: "工作任务",
              label: () => `【${useProjectStore().project.name}】工作任务`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/task/index.vue")
          },
          {
            path: "test",
            name: "ProjectTest",
            meta: {
              title: "项目测试",
              label: () => `【${useProjectStore().project.name}】项目测试`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/test/index.vue")
          },
          {
            path: "issue",
            name: "ProjectIssue",
            meta: {
              title: "项目缺陷",
              label: () => `【${useProjectStore().project.name}】缺陷`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/issue/index.vue")
          },
          {
            path: "risk",
            name: "ProjectRisk",
            meta: {
              title: "项目风险",
              label: () => `【${useProjectStore().project.name}】风险`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/risk/index.vue")
          },
          {
            path: "change",
            name: "ProjectChange",
            meta: {
              title: "项目变更",
              label: () => `【${useProjectStore().project.name}】变更`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/change/index.vue")
          },
          {
            path: "asset",
            name: "ProjectAsset",
            meta: {
              title: "项目资产",
              label: () => `【${useProjectStore().project.name}】资产`,
              activeMenuKey: "ProjectList",
              activeSubMenuKey: "ProjectList"
            },
            component: () => import("@/views/project/detail/asset/index.vue")
          }
        ]
      }
    ]
  },
  {
    path: "/task",
    name: "Task",
    redirect: {
      name: "TaskIndex"
    },
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "index",
        name: "TaskIndex",
        component: () => import("@/views/task/index.vue"),
        meta: {
          label: () => "我的任务",
          activeMenuKey: "Task",
          activeSubMenuKey: "Task"
        }
      }
    ]
  },
  {
    path: "/issue",
    name: "Issue",
    redirect: {
      name: "IssueIndex"
    },
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "index",
        name: "IssueIndex",
        component: () => import("@/views/issue/index.vue"),
        meta: {
          label: () => "我的缺陷",
          activeMenuKey: "Issue",
          activeSubMenuKey: "Issue"
        }
      }
    ]
  },
  {
    path: "/asset",
    name: "Asset",
    component: () => import("@/layout/Primary.vue"),
    meta: {
      activeMenuKey: "Asset",
      activeSubMenuKey: "File"
    },
    children: [
      {
        path: "file",
        name: "File",
        component: () => import("@/views/asset/file/index.vue"),
        meta: {
          label: () => "资产文件管理",
          activeMenuKey: "Asset",
          activeSubMenuKey: "File"
        }
      },
      {
        path: "testcaselib",
        name: "TestCaseLib",
        component: () => import("@/views/asset/testCaseLib/index.vue"),
        meta: {
          label: () => "测试用例库",
          activeMenuKey: "Asset",
          activeSubMenuKey: "TestCaseLib"
        }
      }
    ]
  },
  {
    path: "/system",
    name: "System",
    component: () => import("@/layout/Primary.vue"),
    meta: {
      activeMenuKey: "System",
      activeSubMenuKey: "User"
    },
    children: [
      {
        path: "user",
        name: "User",
        component: () => import("@/views/system/user/index.vue"),
        meta: {
          label: () => "用户管理",
          activeMenuKey: "System",
          activeSubMenuKey: "User"
        }
      },
      {
        path: "department",
        name: "Department",
        component: () => import("@/views/system/department/index.vue"),
        meta: {
          label: () => "部门管理",
          activeMenuKey: "System",
          activeSubMenuKey: "Department"
        }
      },
      {
        path: "role",
        name: "Role",
        component: () => import("@/views/system/role/index.vue"),
        meta: {
          label: () => "角色管理",
          activeMenuKey: "System",
          activeSubMenuKey: "Role"
        }
      },
      {
        path: "stage",
        name: "Stage",
        component: () => import("@/views/system/stage/index.vue"),
        meta: {
          label: () => "阶段管理",
          activeMenuKey: "System",
          activeSubMenuKey: "Stage"
        }
      },
      {
        path: "assetCategory",
        name: "AssetCategory",
        component: () => import("@/views/system/assetCategory/index.vue"),
        meta: {
          label: () => "资产分类管理",
          activeMenuKey: "System",
          activeSubMenuKey: "AssetCategory"
        }
      },
      {
        path: "thirdSource",
        name: "ThirdSource",
        component: () => import("@/views/system/thirdSource/index.vue"),
        meta: {
          label: () => "第三方登录源管理",
          activeMenuKey: "System",
          activeSubMenuKey: "ThirdSource"
        }
      },
      {
        path: "ossConfig",
        name: "OssConfig",
        component: () => import("@/views/system/ossConfig/index.vue"),
        meta: {
          label: () => "OSS配置管理",
          activeMenuKey: "System",
          activeSubMenuKey: "OssConfig"
        }
      }
    ]
  },
  {
    path: "/editor/:id/:versionId?",
    name: "Editor",
    component: () => import("@/views/docEditor/index.vue")
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (!userStore.token && to.name !== "Login" && to.name !== "Auth") {
    localStorage.setItem("redirect", to.fullPath)
  }
  userStore.addRoute({
    url: from.fullPath as string,
    label: (from.meta?.label ? from.meta?.label() : from.name || from.fullPath) as string,
    time: new Date().getTime()
  })
  loadingBar.start()
  if (!userStore.token && to.name !== "Login" && to.name !== "Auth") {
    next({ name: "Login" })
  } else {
    next()
  }
})
router.afterEach((to, from, failure) => {
  loadingBar.finish()
  if (isNavigationFailure(failure)) {
    console.log("failed navigation", failure)
    return
  }
  const appStore = useAppStore()
  appStore.activeMenuKey = to.meta.activeMenuKey as string
  appStore.activeSubMenuKey = to.meta.activeSubMenuKey as string
})

export default router
