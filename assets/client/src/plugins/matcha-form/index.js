import Form from './components/Form'
import FormField from './components/FormField'

export default {
  install(Vue, options) {
    Vue.component('m-form', Form)
    Vue.component('m-form-field', FormField)
  }
}
