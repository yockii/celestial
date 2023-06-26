import { useUserStore } from "@/store/user"
import { Directive } from "vue"

const resourceCode: Directive = {
  mounted(el, binding) {
    const { value } = binding
    const userStore = useUserStore()
    // 判断value是否数组
    if (Array.isArray(value)) {
      if (!value.some((item) => userStore.hasResourceCode(item))) {
        el.parentNode?.removeChild(el)
      }
      return
    }
    if (!userStore.hasResourceCode(value)) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el, binding) {
    const { value } = binding
    if (!useUserStore().hasResourceCode(value)) {
      el.parentNode?.removeChild(el)
    }
  }
}

export default resourceCode
