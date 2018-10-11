import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'

import store from './store'

const instance = axios.create({
  baseURL: 'http://192.168.1.20:3000'
})

instance.interceptors.request.use(function (config) {
  config.headers['X-Auth-Token'] = store.state.authToken
  return config
})

instance.interceptors.response.use(null, function (err) {
  if (err.config.errorHandle === false) {
    return Promise.reject(err)
  } else {
    store.commit('setError', err)
  }
})

Vue.use(
  VueAxios,
  instance
)

export default instance
