import { App, Directive } from "vue"
import resourceCode from "./modules/authorization"
import projectResourceCode from "./modules/projectAuthorization"

const directivesList: { [name: string]: Directive } = {
  resourceCode,
  projectResourceCode
}

const directives = {
  install: function (app: App<Element>) {
    Object.keys(directivesList).forEach((key) => {
      app.directive(key, directivesList[key])
    })
  }
}

export default directives
