import { useUserStore } from "@/store/user"
import { Directive } from "vue"

const resourceCode: Directive = {
  mounted(el, binding) {
    const { value } = binding
    // 判断value是否数组
    if (Array.isArray(value)) {
      if (!value.some((item) => useUserStore().hasResourceCode(item))) {
        el.parentNode?.removeChild(el)
      }
      return
    }
    if (!useUserStore().hasResourceCode(value)) {
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
