import { NIcon } from "naive-ui"
import { Component, h } from "vue"
import {RouterLink} from "vue-router";
import {useAppStore} from "../store/app";

// 渲染icon的通用方法
export function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

export function renderLinkedLabel(menuName: string, routeName: string) {
  return () => {
      useAppStore().activeSubMenuKey = routeName
      return h(
          RouterLink,
          {
              to: {
                  name: routeName
              }
          },
          {
              default: () => menuName
          }
      )
  }
}