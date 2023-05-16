import { NIcon } from "naive-ui"
import { Component, h } from "vue"
import {RouterLink} from "vue-router";
import {useAppStore} from "@/store/app";

// 渲染icon的通用方法
export function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

export function renderLinkedLabel(menuName: string, routeName: string, mainMenu = false) {
  return () => {
      return h(
          RouterLink,
          {
              to: {
                  name: routeName
              },
              onClick: () => {
                  if (mainMenu) {
                      useAppStore().activeSubMenuKey = routeName
                  }
              }
          },
          {
              default: () => menuName
          }
      )
  }
}

// rgb颜色转16进制 "rgb(255,255,255)" => "#ffffff"
export function rgbStringToHex(rgb: string) {
    const reg = /^(rgb|RGB)/
    if (reg.test(rgb)) {
        let strHex = "#"
        const colorArr = rgb.replace(/(?:\(|\)|rgb|RGB)*/g, "").split(",")
        for (const color of colorArr) {
            const hex = Number(color).toString(16)
            if (hex.length === 1) {
                strHex += "0" + hex
            } else {
                strHex += hex
            }
        }
        if (strHex.length !== 7) {
            strHex = rgb
        }
        return strHex
    }
    return rgb
}

// rbg独立给出数值转16进制
export function rgbToHex(r: number, g: number, b: number) {
    let hex = ((r << 16) | (g << 8) | b).toString(16)
    if (hex.length < 6) {
        hex = "0" + hex
    }
    return "#" + hex
}