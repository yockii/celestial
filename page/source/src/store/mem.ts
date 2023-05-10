import { defineStore } from "pinia"
import {renderIcon} from "../utils/Render";
import {Folders, Home, ToolKit, SettingsServices, TaskView, Dashboard, UserMultiple} from "@vicons/carbon";
import {ProjectOutlined} from "@vicons/antd";
import {MenuOption} from "naive-ui";
import {renderLinkedLabel} from '../utils/Render'
import {useAppStore} from "./app";

export const useMemStore = defineStore("mem", {
    state: () => ({
        menus: [
            {
            label: renderLinkedLabel("首页", "Dashboard"),
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
            label: renderLinkedLabel("项目", "ProjectList"),
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
            label: "资产",
            key: "Asset",
            icon: renderIcon(Folders)
        },
            {
                label: renderLinkedLabel("系统", "User"),
                key: "System",
                icon: renderIcon(SettingsServices),
                children: [
                    {
                        label: "用户管理",
                        key: "User",
                        icon: renderIcon(UserMultiple)
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
                    icon: menu.icon,
                })
            }
            return menuOptions
        },
        sideMenus: (state) => {
            const menuOptions: MenuOption[] = []
            const menu = state.menus.find(menu => menu.key === useAppStore().activeMenuKey)
            if (menu) {
                for (const child of menu.children || []) {
                    menuOptions.push({
                        label: child.label,
                        key: child.key,
                        icon: child.icon,
                    })
                }
            }
            return menuOptions
        }
    },
    actions: {
    }
})