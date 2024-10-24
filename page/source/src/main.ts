import { createApp } from "vue"
import "./style.css"
import App from "./App.vue"
import router from "./router"
import { createPinia } from "pinia"
import piniaPluginPersistedstate from "pinia-plugin-persistedstate"

import "normalize.css/normalize.css" // 重置浏览器样式

// import "vfonts/Lato.css"// 通用字体
import "vfonts/FiraCode.css" // 等宽字体
import "uno.css"
import "animate.css"

import svgIcon from "./components/SvgIcon.vue"
import dayjs from "dayjs"
import "dayjs/locale/zh-cn"
import relativeTime from "dayjs/plugin/relativeTime"

import directives from "@/directives"

dayjs.extend(relativeTime)
dayjs.locale("zh-cn")

const app = createApp(App)
app.use(router)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
app.use(pinia)
app.component("svg-icon", svgIcon)
app.use(directives)
app.mount("#app")
