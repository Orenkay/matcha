import Vue from 'vue'
import Vuex from 'vuex'
import cookies from 'js-cookie'

import axios from './axios'
import router from './router'

Vue.use(Vuex)

const createUserData = () => ({
  account: {
    id: undefined,
    user: undefined,
    email: undefined
  },
  profile: {
    firstName: undefined,
    lastName: undefined,
    gender: undefined,
    attraction: undefined,
    bio: undefined
  },
  loc: {
    address: undefined
  },
  pictures: [],
  interests: []
})

const store = new Vuex.Store({
  state: {
    loaded: true,
    error: undefined,
    authToken: cookies.get('auth-token'),
    notificationsEnabled: false,
    userData: createUserData()
  },
  actions: {
    restoreSession: store => {
      return new Promise((resolve, reject) => {
        if (store.getters.isRestored) {
          return resolve()
        }
        store.state.loaded = false
        axios.get('/users/me', { timeout: 5000, errorHandle: false })
          .then(res => {
            const data = res.data.data
            store.state.loaded = true

            console.log(data)
            store.commit('setUserData', ['account', data.account])
            store.commit('setUserData', ['profile', data.profile])
            store.commit('setUserData', ['loc', data.loc])
            store.commit('setUserData', ['interests', data.interests])

            resolve()
          })
          .catch((err) => {
            if (err.response && err.response.status === 401) {
              store.state.loaded = true
              store.commit('logout')
            }
          })
      })
    },
    logout: store => {
      return new Promise((resolve, reject) => {
        axios
          .get('/auth/logout')
          .then(res => {
            store.commit('logout')
            resolve()
          })
          .catch((err) => {
            reject(err)
          })
      })
    }
  },
  getters: {
    inited: state => {
      return state.userData.profile.firstName !== undefined &&
        state.userData.loc.address !== undefined &&
        state.userData.pictures.length &&
        state.userData.interests.length
    },
    isRestored: state => {
      return state.userData.account.id !== undefined
    },
    userData: state => {
      return state.userData
    },
    account: state => {
      return state.userData.account
    },
    profile: state => {
      return state.userData.profile
    },
    loc: state => {
      return state.userData.loc
    },
    pictures: state => {
      return state.userData.pictures
    },
    interests: state => {
      return state.userData.interests
    }
  },
  mutations: {
    login: (state, token) => {
      state.authToken = token
      cookies.set('auth-token', token)
      router.replace('/app')
    },
    logout: (state) => {
      state.authToken = undefined
      cookies.remove('auth-token')
      router.replace('/')
      state.userData = createUserData()
    },
    addInterest: (state, v) => {
      state.userData.interests.push(v)
    },
    removeInterest: (state, v) => {
      const index = state.userData.interests.findIndex(t => t.value === v)
      if (index >= 0) {
        state.userData.interests.splice(index, 1)
      }
    },
    'toggle-notifications': (state) => {
      state.notificationsEnabled = !state.notificationsEnabled
    },
    setUserData: (state, args) => {
      const [k, v] = args
      if (v !== null) {
        state.userData[k] = v
      }
    },
    setError: (state, err) => {
      state.error = err
    }
  }
})

export default store
