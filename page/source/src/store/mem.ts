import { defineStore } from "pinia"
import { renderIcon } from "@/utils/Render"
import {
  IdManagement,
  SkillLevelIntermediate,
  Folders,
  Home,
  ToolKit,
  SettingsServices,
  TaskView,
  Dashboard,
  UserMultiple,
  CategoryNewEach,
  LicenseGlobal
} from "@vicons/carbon"
import { ProjectOutlined } from "@vicons/antd"
import { Files, TestPipe } from "@vicons/tabler"
import { MenuOption } from "naive-ui"
import { renderLinkedLabel } from "@/utils/Render"
import { useAppStore } from "./app"

// const TICKER_INTERVAL = 30 * 1000;
// let tickerPointer:number;

export const useMemStore = defineStore("mem", {
  state: () => ({
    // now: new Date,
    menus: [
      {
        label: renderLinkedLabel("首页", "Dashboard", true),
        key: "Home",
        icon: renderIcon(Home),
        children: [
          {
            label: renderLinkedLabel("仪表盘", "Dashboard"),
            key: "Dashboard",
            icon: renderIcon(Dashboard)
          }
        ]
      },
      {
        label: renderLinkedLabel("项目", "ProjectList", true),
        key: "ProjectList",
        icon: renderIcon(ProjectOutlined)
      },
      {
        label: "任务",
        key: "Task",
        icon: renderIcon(TaskView)
      },
      {
        label: "测试",
        key: "Test",
        icon: renderIcon(ToolKit)
      },
      {
        label: renderLinkedLabel("资产", "File", true),
        key: "Asset",
        icon: renderIcon(Folders),
        children: [
          {
            label: renderLinkedLabel("文件管理", "File"),
            key: "File",
            icon: renderIcon(Files)
          },
          {
            label: renderLinkedLabel("测试用例库", "TestCaseLib"),
            key: "TestCaseLib",
            icon: renderIcon(TestPipe)
          }
        ]
      },
      {
        label: renderLinkedLabel("系统", "User", true),
        key: "System",
        icon: renderIcon(SettingsServices),
        children: [
          {
            label: renderLinkedLabel("用户管理", "User"),
            key: "User",
            icon: renderIcon(UserMultiple)
          },
          {
            label: renderLinkedLabel("角色管理", "Role"),
            key: "Role",
            icon: renderIcon(IdManagement)
          },
          {
            label: renderLinkedLabel("阶段管理", "Stage"),
            key: "Stage",
            icon: renderIcon(SkillLevelIntermediate)
          },
          {
            label: renderLinkedLabel("资产目录", "AssetCategory"),
            key: "AssetCategory",
            icon: renderIcon(CategoryNewEach)
          },
          {
            label: renderLinkedLabel("登录源管理", "ThirdSource"),
            key: "ThirdSource",
            icon: renderIcon(LicenseGlobal)
          }
        ]
      }
    ]
  }),
  getters: {
    mainMenus: (state) => {
      const menuOptions: MenuOption[] = []
      for (const menu of state.menus) {
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
      const menu = state.menus.find((menu) => menu.key === useAppStore().activeMenuKey)
      if (menu) {
        for (const child of menu.children || []) {
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
