import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router"
import { useUserStore } from "@/store/user"
import { createDiscreteApi } from "naive-ui"
import { useProjectStore } from "@/store/project"
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
      label: () => "首页"
    },
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/index/dashboard/index.vue"),
        meta: {
          label: () => "首页"
        }
      }
    ]
  },
  {
    path: "/project",
    name: "Project",
    component: () => import("@/layout/Primary.vue"),
    meta: {
      label: () => "项目列表"
    },
    children: [
      {
        path: "list",
        name: "ProjectList",
        component: () => import("@/views/project/list/index.vue"),
        meta: {
          label: () => "项目列表"
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
          label: () => `【${useProjectStore().project.name}】总览`
        },
        children: [
          {
            path: "dashboard",
            name: "ProjectDashboard",
            meta: {
              title: "项目总览",
              label: () => `【${useProjectStore().project.name}】总览`
            },
            component: () => import("@/views/project/detail/dashboard/index.vue")
          },
          {
            path: "plan",
            name: "ProjectPlan",
            meta: {
              title: "项目计划",
              label: () => `【${useProjectStore().project.name}】计划`
            },
            component: () => import("@/views/project/detail/plan/index.vue")
          },
          {
            path: "module",
            name: "ProjectModule",
            meta: {
              title: "功能模块",
              label: () => `【${useProjectStore().project.name}】功能模块`
            },
            component: () => import("@/views/project/detail/module/index.vue")
          },
          {
            path: "requirement",
            name: "ProjectRequirement",
            meta: {
              title: "项目需求",
              label: () => `【${useProjectStore().project.name}】需求`
            },
            component: () => import("@/views/project/detail/requirement/index.vue")
          },
          {
            path: "task",
            name: "ProjectTask",
            meta: {
              title: "工作任务",
              label: () => `【${useProjectStore().project.name}】工作任务`
            },
            component: () => import("@/views/project/detail/task/index.vue")
          },
          {
            path: "test",
            name: "ProjectTest",
            meta: {
              title: "项目测试",
              label: () => `【${useProjectStore().project.name}】项目测试`
            },
            component: () => import("@/views/project/detail/test/index.vue")
          },
          {
            path: "issue",
            name: "ProjectIssue",
            meta: {
              title: "项目缺陷",
              label: () => `【${useProjectStore().project.name}】缺陷`
            },
            component: () => import("@/views/project/detail/issue/index.vue")
          },
          {
            path: "risk",
            name: "ProjectRisk",
            meta: {
              title: "项目风险",
              label: () => `【${useProjectStore().project.name}】风险`
            },
            component: () => import("@/views/project/detail/risk/index.vue")
          },
          {
            path: "change",
            name: "ProjectChange",
            meta: {
              title: "项目变更",
              label: () => `【${useProjectStore().project.name}】变更`
            },
            component: () => import("@/views/project/detail/change/index.vue")
          },
          {
            path: "asset",
            name: "ProjectAsset",
            meta: {
              title: "项目资产",
              label: () => `【${useProjectStore().project.name}】资产`
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
      name: "TaskList"
    },
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "list",
        name: "TaskList",
        component: () => import("@/views/task/list/index.vue"),
        meta: {
          label: () => "我的任务"
        }
      }
    ]
  },
  {
    path: "/test",
    name: "Test",
    redirect: {
      name: "TestList"
    },
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "list",
        name: "TestList",
        component: () => import("@/views/test/list/index.vue"),
        meta: {
          label: () => "我的测试"
        }
      }
    ]
  },
  {
    path: "/asset",
    name: "Asset",
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "file",
        name: "File",
        component: () => import("@/views/asset/file/index.vue"),
        meta: {
          label: () => "资产文件管理"
        }
      },
      {
        path: "testcaselib",
        name: "TestCaseLib",
        component: () => import("@/views/asset/testCaseLib/index.vue"),
        meta: {
          label: () => "测试用例库"
        }
      }
    ]
  },
  {
    path: "/system",
    name: "System",
    component: () => import("@/layout/Primary.vue"),
    children: [
      {
        path: "user",
        name: "User",
        component: () => import("@/views/system/user/index.vue"),
        meta: {
          label: () => "用户管理"
        }
      },
      {
        path: "role",
        name: "Role",
        component: () => import("@/views/system/role/index.vue"),
        meta: {
          label: () => "角色管理"
        }
      },
      {
        path: "stage",
        name: "Stage",
        component: () => import("@/views/system/stage/index.vue"),
        meta: {
          label: () => "阶段管理"
        }
      },
      {
        path: "assetCategory",
        name: "AssetCategory",
        component: () => import("@/views/system/assetCategory/index.vue"),
        meta: {
          label: () => "资产分类管理"
        }
      },
      {
        path: "thirdSource",
        name: "ThirdSource",
        component: () => import("@/views/system/thirdSource/index.vue"),
        meta: {
          label: () => "第三方登录源管理"
        }
      },
      {
        path: "ossConfig",
        name: "OssConfig",
        component: () => import("@/views/system/ossConfig/index.vue"),
        meta: {
          label: () => "OSS配置管理"
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
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
router.afterEach(() => {
  loadingBar.finish()
})

export default router
