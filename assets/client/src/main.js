import 'babel-polyfill'
import 'buefy/dist/buefy.css'
import Vue from 'vue'
import Buefy from 'buefy'
import * as VueGoogleMaps from 'vue2-google-maps'

import BuefySteps from './plugins/buefy-steps'
import MatchaForm from './plugins/matcha-form'
import App from './App.vue'
import store from './store'
import router from './router'

import './axios'

Vue.use(Buefy)
Vue.use(BuefySteps)
Vue.use(MatchaForm)

Vue.use(VueGoogleMaps, {
  load: {
    key: 'AIzaSyBciG2mapnXT-z59x40gmE_cT_7W61Mb8M',
    libraries: 'places'
  }
})

// Just an alias ...
Vue.use({
  install (Vue) {
    Vue.prototype.$toast.error = function (message) {
      this.open({
        type: 'is-danger',
        queue: false,
        message
      })
    }
  }
})

// eslint-disable-next-line
new Vue({
  el: '#app',
  store,
  router,
  render: h => h(App)
})
