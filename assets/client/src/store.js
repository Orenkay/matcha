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
  interests: [],
  matches: []
})

const store = new Vuex.Store({
  state: {
    loaded: true,
    error: undefined,
    notifications: JSON.parse(cookies.get('notifications') || '[]'),
    profilesCache: {},
    messages: {},
    authToken: cookies.get('auth-token'),
    userData: createUserData()
  },
  actions: {
    getProfile: (store, id) => {
      return new Promise((resolve, reject) => {
        if (store.state.profilesCache[id] !== undefined) {
          return resolve(store.state.profilesCache[id])
        }
        axios
          .get(`/profiles/${id}`)
          .then(res => {
            if (!res.data.data) return resolve(undefined)
            const data = res.data.data

            if (!data.pictures || data.pictures.findIndex(p => p.isPP) < 0) {
              return resolve(undefined)
            }

            if (!data.interests || data.interests.length <= 0) {
              return resolve(undefined)
            }

            data.birthdate = data.birthdate * 1000
            store.state.profilesCache[id] = data
            resolve(store.state.profilesCache[id])
          })
          .catch(err => {
            if (err.response && err.response.status === 400) {
              return resolve(undefined)
            }
            reject(err)
          })
      })
    },
    addMessage: (store, d) => {
      store.dispatch('getMessagesFrom', d.convId).then(messages => {
        if (messages.findIndex(m => m.id === d.message.id) < 0) {
          messages.push(d.message)
          messages = messages.sort((a, b) => a.date > b.date)
        }
      })
    },
    removeMessageFrom: (store, id) => {
      store.state.messages[id] = undefined
    },
    getMessagesFrom: (store, id) => {
      return new Promise((resolve, reject) => {
        if (store.state.messages[id]) {
          return resolve(store.state.messages[id])
        }
        axios.get(`/messages/me/${id}`).then(res => {
          if (res.data.data === null) {
            res.data.data = []
          }
          store.state.messages[id] = res.data.data.sort((a, b) => a.date > b.date)
          return resolve(store.state.messages[id])
        })
      })
    },
    restoreSession: store => {
      return new Promise((resolve, reject) => {
        if (store.getters.isRestored) {
          return resolve()
        }
        store.state.loaded = false
        axios.get('/users/me', { timeout: 5000 })
          .then(res => {
            const data = res.data.data
            store.state.loaded = true

            store.commit('setUserData', ['account', data.account])
            store.commit('setUserData', ['profile', data.meta.profile])
            store.commit('setUserData', ['loc', data.meta.loc])
            store.commit('setUserData', ['interests', data.meta.interests])
            store.commit('setUserData', ['pictures', data.meta.pictures])
            store.commit('setUserData', ['matches', data.meta.matches])

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
        state.userData.pictures.findIndex(p => p.isPP) >= 0 &&
        state.userData.interests.length
    },
    notifications: state => {
      return state.notifications
    },
    isRestored: state => {
      return state.userData.account.id !== undefined
    },
    loaded: state => {
      return state.loaded
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
    loading: (state, v) => {
      state.loaded = !v
    },
    removeMatch: (state, v) => {
      const index = state.userData.matches.findIndex(m => m === v)
      if (index >= 0) {
        state.userData.matches.splice(index, 1)
      }
    },
    addMatch: (state, v) => {
      state.userData.matches.push(v)
    },
    removeNotification: (state, index) => {
      state.notifications.splice(index, 1)
      cookies.set('notifications', state.notifications)
    },
    addNotification: (state, data) => {
      state.notifications.push(data)
      cookies.set('notifications', state.notifications)
    },
    clearNotifications: state => {
      state.notifications = []
      cookies.remove('notifications')
    },
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
      state.notifications = []
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
    setUserData: (state, args) => {
      const [k, v] = args
      if (v !== null) {
        state.userData[k] = v
      }
    },
    addPicture: (state, pic) => {
      state.userData.pictures.push(pic)
    },
    removePicture: (state, id) => {
      const index = state.userData.pictures.findIndex(p => p.id === id)
      if (index >= 0) {
        state.userData.pictures.splice(index, 1)
      }
    },
    setPP: (state, id) => {
      const pp = state.userData.pictures.find(p => p.isPP)
      const newPP = state.userData.pictures.find(p => p.id === id)

      pp && (pp.isPP = false)
      newPP && (newPP.isPP = true)
    },
    setError: (state, err) => {
      state.error = err
    }
  }
})

export default store
