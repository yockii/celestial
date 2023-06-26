import { useProjectStore } from "@/store/project"
import { Directive } from "vue"
import { useUserStore } from "@/store/user"

const projectResourceCode: Directive = {
  mounted(el, binding) {
    const { value } = binding
    const projectStore = useProjectStore()
    const userStore = useUserStore()
    // 判断value是否数组
    if (Array.isArray(value)) {
      if (!value.some((item) => projectStore.hasResourceCode(item) || userStore.hasResourceCode(item))) {
        el.parentNode?.removeChild(el)
      }
      return
    }
    if (!projectStore.hasResourceCode(value) && !userStore.hasResourceCode(value)) {
      el.parentNode?.removeChild(el)
    }
  },
  updated(el, binding) {
    const { value } = binding
    if (!useProjectStore().hasResourceCode(value) && !useUserStore().hasResourceCode(value)) {
      el.parentNode?.removeChild(el)
    }
  }
}

export default projectResourceCode
