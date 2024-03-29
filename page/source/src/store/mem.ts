import { defineStore } from "pinia"
import { renderIcon } from "@/utils/Render"
import {
  IdManagement,
  SkillLevelIntermediate,
  Folders,
  Home,
  Debug,
  SettingsServices,
  TaskView,
  Dashboard,
  UserMultiple,
  CategoryNewEach,
  LicenseGlobal,
  ObjectStorage,
  TimePlot
} from "@vicons/carbon"
import { MeetingRoomOutlined } from "@vicons/material"
import { Organization28Regular } from "@vicons/fluent"
import { ProjectOutlined } from "@vicons/antd"
import { Files, TestPipe } from "@vicons/tabler"
import { MenuOption } from "naive-ui"
import { renderLinkedLabel } from "@/utils/Render"
import { useAppStore } from "./app"
import { useUserStore } from "./user"

// const TICKER_INTERVAL = 30 * 1000;
// let tickerPointer:number;

export const useMemStore = defineStore("mem", {
  state: () => ({
    // now: new Date,
    menus: [
      {
        label: renderLinkedLabel("首页", "Dashboard", true),
        key: "Home",
        code: "home",
        icon: renderIcon(Home),
        children: [
          {
            label: renderLinkedLabel("仪表盘", "Dashboard"),
            key: "Dashboard",
            code: "home:dashboard",
            icon: renderIcon(Dashboard)
          }
        ]
      },
      {
        label: renderLinkedLabel("项目", "ProjectList", true),
        key: "ProjectList",
        code: "project",
        icon: renderIcon(ProjectOutlined)
      },
      {
        label: renderLinkedLabel("任务", "Task", true),
        key: "Task",
        code: "task",
        icon: renderIcon(TaskView)
      },
      {
        label: renderLinkedLabel("缺陷", "Issue", true),
        key: "Issue",
        code: "issue",
        icon: renderIcon(Debug)
      },
      {
        label: renderLinkedLabel("资产", "File", true),
        key: "Asset",
        code: "asset",
        icon: renderIcon(Folders),
        children: [
          {
            label: renderLinkedLabel("文件管理", "File"),
            key: "File",
            code: "asset:file",
            icon: renderIcon(Files)
          },
          {
            label: renderLinkedLabel("测试用例库", "TestCaseLib"),
            key: "TestCaseLib",
            code: "asset:commonTestCase",
            icon: renderIcon(TestPipe)
          }
        ]
      },
      {
        label: renderLinkedLabel("工时", "WorkTime", true),
        key: "WorkTime",
        code: "workTime",
        icon: renderIcon(TimePlot)
      },
      {
        label: renderLinkedLabel("会议室", "MeetingRoom", true),
        key: "MeetingRoom",
        code: "meetingRoom",
        icon: renderIcon(MeetingRoomOutlined)
      },
      {
        label: renderLinkedLabel("系统", "User", true),
        key: "System",
        code: "system",
        icon: renderIcon(SettingsServices),
        children: [
          {
            label: renderLinkedLabel("用户管理", "User"),
            key: "User",
            code: "system:user",
            icon: renderIcon(UserMultiple)
          },
          {
            label: renderLinkedLabel("部门管理", "Department"),
            key: "Department",
            code: "system:department",
            icon: renderIcon(Organization28Regular)
          },
          {
            label: renderLinkedLabel("角色管理", "Role"),
            key: "Role",
            code: "system:role",
            icon: renderIcon(IdManagement)
          },
          {
            label: renderLinkedLabel("阶段管理", "Stage"),
            key: "Stage",
            code: "system:stage",
            icon: renderIcon(SkillLevelIntermediate)
          },
          {
            label: renderLinkedLabel("资产目录", "AssetCategory"),
            key: "AssetCategory",
            code: "system:assetCategory",
            icon: renderIcon(CategoryNewEach)
          },
          {
            label: renderLinkedLabel("登录源管理", "ThirdSource"),
            key: "ThirdSource",
            code: "system:thirdSource",
            icon: renderIcon(LicenseGlobal)
          },
          {
            label: renderLinkedLabel("Oss配置管理", "OssConfig"),
            key: "OssConfig",
            code: "system:ossConfig",
            icon: renderIcon(ObjectStorage)
          }
        ]
      }
    ]
  }),
  getters: {
    mainMenus: (state) => {
      const menuOptions: MenuOption[] = []
      const userStore = useUserStore()
      for (const menu of state.menus) {
        if (userStore.hasResourceCode(menu.code || "") === false) continue
        menuOptions.push({
          label: menu.label,
          key: menu.key,
          icon: menu.icon
        })
      }
      return menuOptions
    },
    sideMenus: (state) => {
      const menuOptions: MenuOption[] = []
      const userStore = useUserStore()
      const menu = state.menus.find((menu) => menu.key === useAppStore().activeMenuKey)
      if (menu) {
        for (const child of menu.children || []) {
          if (userStore.hasResourceCode(menu.code || "") === false) continue
          menuOptions.push({
            label: child.label,
            key: child.key,
            icon: child.icon
          })
        }
      }
      return menuOptions
    }
  },
  actions: {
    // startTicker() {
    //     if (!tickerPointer) {
    //         tickerPointer = window.setInterval(() => {
    //             this.now = new Date
    //         }, TICKER_INTERVAL)
    //     }
    // },
    // haltTicker() {
    //     window.clearInterval(tickerPointer)
    //     tickerPointer = 0
    // }
  }
})
