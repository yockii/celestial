import axios, { AxiosInstance, AxiosRequestConfig } from "axios"
import { useUserStore } from "@/store/user"
import { createDiscreteApi, MessageReactive } from "naive-ui"
import { Result } from "@/types/common"
const { message } = createDiscreteApi(["message"])

let messageReactive: MessageReactive | null = null

const renderUnAuthMsg = (msg: string) => {
  messageReactive = message.error(
    () =>
      h("div", {}, [
        msg + ", 请点击",
        h(
          "a",
          {
            href: "/login",
            onClick: () => {
              if (messageReactive) {
                messageReactive.destroy()
                messageReactive = null
              }
            }
          },
          "重新登录"
        )
      ]),
    {
      duration: 0,
      closable: true
    }
  )
}

export class Request {
  instance: AxiosInstance
  baseConfig: AxiosRequestConfig = {
    baseURL: "/api/v1",
    timeout: 10000
  }

  stopResp = false

  stopRequest() {
    this.stopResp = true
    setTimeout(() => {
      this.stopResp = false
    }, 1000)
  }

  constructor(config: AxiosRequestConfig) {
    this.instance = axios.create(Object.assign(this.baseConfig, config))

    this.instance.interceptors.request.use(
      (config) => {
        const token = useUserStore().token
        if (token && token !== "") {
          config.headers.Authorization = "Bearer " + token
        }
        return config
      },
      (error: any) => {
        return Promise.reject(error)
      }
    )

    this.instance.interceptors.response.use(
      (response) => {
        const res = response.data as Result<any>
        if (res.code !== 0) {
          // 全局弹出框展示
          message.error(res.msg)
          return Promise.reject(res.msg)
        }
        return res.data
      },
      (error: any) => {
        if (this.stopResp) {
          Promise.reject("已停止请求")
          return
        }
        let msg: string
        let noPopup = false
        switch (error.response.status) {
          case 400:
            msg = "请求错误(400)"
            break
          case 401:
            msg = "未授权，请重新登录(401)"
            renderUnAuthMsg(msg)
            noPopup = true
            this.stopRequest()
            break
          case 403:
            msg = "拒绝访问(403)"
            renderUnAuthMsg(msg)
            noPopup = true
            this.stopRequest()
            break
          case 404:
            msg = "请求出错(404)"
            break
          case 408:
            msg = "请求超时(408)"
            break
          case 500:
            msg = "服务器错误(500)"
            break
          case 501:
            msg = "服务未实现(501)"
            break
          case 502:
            msg = "网络错误(502)"
            break
          case 503:
            msg = "服务不可用(503)"
            break
          case 504:
            msg = "网络超时(504)"
            break
          case 505:
            msg = "HTTP版本不受支持(505)"
            break
          default:
            msg = `连接出错(${error.response.status})!`
        }
        if (!noPopup) {
          // 全局弹出框展示
          message.error(msg)
        }
        return Promise.reject(msg)
      }
    )
  }

  public request<T>(config: AxiosRequestConfig): Promise<T> {
    return this.stopResp ? Promise.reject<T>("停止执行请求中") : this.instance.request(config)
  }

  public get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.stopResp ? Promise.reject<T>("停止执行请求中") : this.instance.get(url, config)
  }

  public post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return this.stopResp ? Promise.reject<T>("停止执行请求中") : this.instance.post(url, data, config)
  }

  public put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return this.stopResp ? Promise.reject<T>("停止执行请求中") : this.instance.put(url, data, config)
  }

  public delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return this.stopResp ? Promise.reject<T>("停止执行请求中") : this.instance.delete(url, config)
  }
}
