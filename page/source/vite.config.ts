import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import * as path from "path"
import eslintPlugin from "vite-plugin-eslint"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { NaiveUiResolver } from "unplugin-vue-components/resolvers"
import Unocss from "unocss/vite"
import presetUno from "@unocss/preset-uno"
import { presetAttributify } from "unocss"
import { createSvgIconsPlugin } from "vite-plugin-svg-icons"

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: "../build",
    emptyOutDir: true
  },
  plugins: [
    vue(),
    eslintPlugin({
      exclude: ["/virtual:/**", "/node_modules/**"],
      include: ["src/**/*.{vue,js,ts,jsx,tsx}"],
    }),
    AutoImport({
      imports: [
        "vue",
        "vue-router",
        {
          "naive-ui": ["useDialog", "useMessage", "useNotification", "useLoadingBar"]
        }
      ],
      eslintrc: {
        enabled: false
      }
    }),
    Components({
      resolvers: [NaiveUiResolver()]
    }),
    Unocss({
      presets: [presetUno(), presetAttributify()]
    }),
    createSvgIconsPlugin({
      // Specify the icon folder to be cached
      iconDirs: [path.resolve(process.cwd(), "src/assets")],
      // Specify symbolId format
      symbolId: "icon-[dir]-[name]"
    }),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src")
    }
  },
  server: {
    port: 3000,
    proxy: {
      "/api/v1": {
        target: "http://127.0.0.1:8086",
        changeOrigin: true
        // rewrite: (path) => path.replace(/^\/api/, ''), // 一般我的接口用的是/api/v1/xxx，所以这里要不用去掉/api
      }
    }
  }
})
