// eslint-disable-next-line spaced-comment
/// <reference types="vite/client" />

import "vue-router"

declare module "*.vue" {
  import Vue from "vue"
  export default Vue
}

declare module "vue-router" {
  interface RouteMeta extends Record<string | number | symbol, unknown> {
    label?: () => string
    title?: string
    activeMenuKey?: string
    activeSubMenuKey?: string
  }
}
