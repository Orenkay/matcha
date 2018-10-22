<template>
  <div v-if="isLoaded">
    <router-link to="/app/messages">back</router-link>
    <h1 class="title">
      {{ targetData.profile.lastName }} {{ targetData.profile.firstName }}
    </h1>
    <hr >
    <div v-if="blocked">
      You have blocked this user
    </div>
    <div v-else>
      <article 
        v-for="(m, index) in messages" 
        v-if="messages.length > 0" 
        :key="index" 
        class="media">
        <figure class="media-left">
          <div 
            v-if="!isMe(m.sender)" 
            :style="`background-image: url('${pp(m.sender)}')`" 
            class="image is-32x32"/>
        </figure>
        <div class="media-content">
          <div 
            :class="isMe(m.sender) && 'is-me'" 
            class="content">
            <b-tooltip 
              :label="messageDate(m.date)" 
              type="is-white" 
              position="is-bottom">
              <p>
                {{ m.message }}
              </p>
            </b-tooltip>
          </div>
        </div>
        <figure class="media-right">
          <div 
            v-if="isMe(m.sender)" 
            :style="`background-image: url('${pp(m.sender)}')`" 
            class="image is-32x32"/>
        </figure>
      </article>
      <article class="media">
        <div class="media-content">
          <b-field>
            <b-input 
              v-model="message" 
              maxlength="200" 
              placeholder="Type your message here" 
              type="textarea"/>
          </b-field>
          <b-field>
            <button 
              class="button" 
              @click="send">Send message</button>
          </b-field>
        </div>
      </article>
    </div>
  </div>
</template>

<script>
import moment from 'moment'

export default {
  data() {
    return {
      blocked: false,
      targetData: undefined,
      message: '',
      messages: []
    }
  },
  beforeRouteUpdate(to, from, next) {
    const id = parseInt(to.params.id)
    if (this.userData.matches.findIndex(m => m === id) < 0) {
      return next('/app/messages')
    }
    return next()
  },
  computed: {
    isLoaded() {
      return this.targetData !== undefined
    },
    userData() {
      return this.$store.getters.userData
    },
    targetId() {
      return this.$route.params.id
    }
  },
  created() {
    this.$store.dispatch('getProfile', this.targetId).then(data => {
      if (!data) {
        return this.$router.push('/app/messages')
      }
      this.targetData = data
      this.$http.get(`/block/${this.targetId}/blocked`).then(res => {
        this.blocked = res.data.data
        if (!this.blocked) {
          this.$store.dispatch('getMessagesFrom', this.targetId).then(data => {
            this.messages = data
          })
        }
      })
    })
  },
  updated() {
    window.scrollTo(0, document.body.scrollHeight)
  },
  methods: {
    pp(id) {
      const data = this.isMe(id) ? this.userData : this.targetData
      return data.pictures.find(p => p.isPP).path
    },
    messageDate(time) {
      return moment(time * 1000).format('DD/MM/YYYY HH:mm:ss')
    },
    isMe(id) {
      return this.userData.account.id === id
    },
    send() {
      this.$http
        .post('/messages/me/' + this.targetId, { message: this.message })
        .then(res => {
          this.$store.dispatch('addMessage', {
            message: res.data.data,
            convId: parseInt(this.targetId)
          })
          this.message = ''
        })
        .catch(err => {
          if (err.response && err.response.status === 400) {
            this.$toast.error(err.response.data.error)
          }
        })
    }
  }
}
</script>

<style scoped lang="scss">
.image {
  background-position: 50% 50%;
  background-size: 100%;
  border-radius: 32px;
}
.media {
  align-items: center;
}
.content {
  display: flex;
  p {
    background: rgba(0, 0, 0, 0.05);
    padding: 0 7px;
    border-radius: 20px;
  }
  &.is-me {
    justify-content: flex-end;
    p {
      background: #0084ff;
      color: white;
    }
  }
}
</style>
