import 'babel-polyfill'
import 'buefy/dist/buefy.css'
import Vue from 'vue'
import Buefy from 'buefy'

import BuefySteps from './plugins/buefy-steps'
import MatchaForm from './plugins/matcha-form'
import App from './App.vue'
import store from './store'
import router from './router'

import './axios'

Vue.use(Buefy)
Vue.use(BuefySteps)
Vue.use(MatchaForm)

// eslint-disable-next-line
new Vue({
  el: '#app',
  store,
  router,
  render: h => h(App)
})
