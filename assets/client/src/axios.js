import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

import store from './store'

const instance = axios.create({
  baseURL: 'http://0.0.0.0:3000'
})

instance.interceptors.request.use(function(config) {
  config.headers['X-Auth-Token'] = store.state.authToken
  return config
})

instance.interceptors.response.use(null, function(err) {
  if (!err.response || err.response.status === 500) {
    store.commit('setError', err)
  }
  return Promise.reject(err)
})

Vue.use(VueAxios, instance)

export default instance
