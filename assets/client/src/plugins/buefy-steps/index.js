import './bulma-steps.min.css'

import Steps from './components/Steps'
import StepItem from './components/StepItem'

export default {
  install(Vue, options) {
    Vue.component('b-steps', Steps)
    Vue.component('b-step-item', StepItem)
  }
}
