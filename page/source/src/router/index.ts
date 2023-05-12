import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router"
import { useUserStore } from "@/store/user"
import { createDiscreteApi } from "naive-ui"
import Primary from "@/layout/Primary.vue"
const { loadingBar } = createDiscreteApi(["loadingBar"])

const routes: Array<RouteRecordRaw> = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/login/index.vue")
  },
  {
    path: "/index",
    name: "Home",
    alias: "/",
    component: () => Primary,
    redirect: { name: "Dashboard" },
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/index/dashboard/index.vue")
      }
    ]
  },
  {
    path: "/project",
    name: "Project",
    component: () => Primary,
    children: [
      {
        path: "list",
        name: "ProjectList",
        component: () => import("@/views/project/list/index.vue")
      }
    ]
  },
  {
    path: "/task",
    name: "Task",
    component: () => Primary,
    children: [
      {
        path: "list",
        name: "TaskList",
        component: () => import("@/views/project/list/index.vue")
      }
    ]
  },
  {
    path: "/test",
    name: "Test",
    component: () => Primary,
    children: [
      {
        path: "list",
        name: "TestList",
        component: () => import("@/views/project/list/index.vue")
      }
    ]
  },
  {
    path: "/asset",
    name: "Asset",
    component: () => Primary,
    children: [
      {
        path: "list",
        name: "AssetList",
        component: () => import("@/views/project/list/index.vue")
      }
    ]
  },
  {
    path: "/system",
    name: "System",
    component: () => Primary,
    children: [
      {
        path: "user",
        name: "User",
        component: () => import("@/views/system/user/index.vue")
      },
      {
        path: "stage",
        name: "Stage",
        component: () => import("@/views/system/stage/index.vue")
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  loadingBar.start()
  const userStore = useUserStore()
  if (!userStore.token && to.name !== "Login") {
    next({ name: "Login" })
  } else {
    next()
  }
})
router.afterEach(() => {
  loadingBar.finish()
})

export default router
