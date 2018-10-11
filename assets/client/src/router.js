import Vue from 'vue'
import VueRouter from 'vue-router'

import store from './store'
import Home from './routes/Home.vue'
import App from './routes/App.vue'
import AppNew from './routes/app/New.vue'
import AppProfile from './routes/app/Profile.vue'
import AppProfileEdit from './routes/app/ProfileEdit.vue'
import AppAccountEdit from './routes/app/AccountEdit.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: Home,
    beforeEnter (to, from, next) {
      store.state.authToken !== undefined ? next('/app') : next()
    }
  },
  {
    path: '/app',
    component: App,
    children: [
      {
        path: 'new',
        component: AppNew,
        beforeEnter (to, from, next) {
          store.getters.inited ? next('/app') : next()
        }
      },
      {
        path: 'profile',
        component: AppProfile
      },
      {
        path: 'profile/edit',
        component: AppProfileEdit
      },
      {
        path: 'account/edit',
        component: AppAccountEdit
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

router.beforeEach(async (to, from, next) => {
  if (RegExp('/app.*').test(to.path)) {
    if (store.state.authToken === undefined) {
      return next('/')
    }
    const userData = store.getters.userData
    if (userData.account.id === undefined) {
      await store.dispatch('restoreSession')
    }
    if (to.path !== '/app/new' && !store.getters.inited) {
      return next('/app/new')
    }
    next()
  } else {
    next()
  }
})

export default router
