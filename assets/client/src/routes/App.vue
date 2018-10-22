<template>
  <div class="app-container">
    <navbar/>
    <div class="app-body">
      <div class="container">
        <router-view/>
      </div>
    </div>
  </div>
</template>

<script>
import Navbar from '../components/NavBar'
export default {
  components: {
    Navbar
  },
  created() {
    this.notifInterval = setInterval(this.getNotifs, 1000)
    this.heartbeatInterval = setInterval(this.heartbeat, 1000 * 60)
    this.heartbeat()
  },
  beforeDestroy() {
    clearInterval(this.notifInterval)
    clearInterval(this.heartbeatInterval)
  },
  methods: {
    heartbeat() {
      this.$http
        .patch('/users/me/heartbeat')
        .then(() => {})
        .catch(() => {})
    },
    getNotifs() {
      this.$http.get('/notifs/me').then(res => {
        if (res.data === null) {
          return
        }
        const notifications = res.data.data
        if (notifications === null) {
          return
        }
        notifications.forEach(n => {
          switch (n.event) {
            case 'visite':
              this.$store
                .dispatch('getProfile', n.data.fromId)
                .then(userData => {
                  const { lastName, firstName } = userData.profile
                  this.$store.commit('addNotification', {
                    message: `${lastName} ${firstName} visited your profile`,
                    to: `/app/profile/${n.data.fromId}`
                  })
                })
              break
            case 'message':
              const convId =
                n.data.sender === this.$store.getters.userData.account.id
                  ? n.data.reciever
                  : n.data.sender
              this.$store.dispatch('addMessage', { message: n.data, convId })
              this.$store.dispatch('getProfile', convId).then(userData => {
                const { lastName, firstName } = userData.profile
                this.$store.commit('addNotification', {
                  message: `${lastName} ${firstName}: ${n.data.message}`,
                  to: `/app/messages/${convId}`
                })
              })
              break
            case 'like':
              this.$store.dispatch('getProfile', n.data).then(userData => {
                const { lastName, firstName } = userData.profile
                this.$store.commit('addNotification', {
                  message: `${lastName} ${firstName} liked you`,
                  to: `/app/profile/${n.data}`
                })
              })
              break
            case 'unlike':
              this.$store.dispatch('getProfile', n.data).then(userData => {
                const { lastName, firstName } = userData.profile
                this.$store.commit('addNotification', {
                  message: `${lastName} ${firstName} doesnt like you anymore`,
                  to: `/app/profile/${n.data}`
                })
              })
              break
            case 'match':
              this.$store.commit('addMatch', n.data)
              this.$store.dispatch('getProfile', n.data).then(userData => {
                const { lastName, firstName } = userData.profile
                this.$store.commit('addNotification', {
                  message: `You got a match with ${lastName} ${firstName}`,
                  to: `/app/profile/${n.data}`
                })
              })
              break
            case 'unmatch':
              this.$store.commit('removeMatch', n.data)
              this.$store.dispatch('getProfile', n.data).then(userData => {
                const { lastName, firstName } = userData.profile
                this.$store.commit('addNotification', {
                  message: `You dont match anymore with ${lastName} ${firstName}`,
                  to: `/app/profile/${n.data}`
                })
              })
              break
          }
        })
      })
    }
  }
}
</script>

<style scoped>
#error {
  display: flex;
  text-align: center;
  align-items: center;
  justify-content: center;
}

.app-container {
  width: 100%;
}
.app-body {
  margin-top: 2rem;
}
</style>
