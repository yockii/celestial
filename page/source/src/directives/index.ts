import { App, Directive } from "vue"
import resourceCode from "./modules/authorization"

const directivesList: { [name: string]: Directive } = {
  resourceCode
}

const directives = {
  install: function (app: App<Element>) {
    Object.keys(directivesList).forEach((key) => {
      app.directive(key, directivesList[key])
    })
  }
}

export default directives
