<template>
  <b-modal 
    :active="opened" 
    :on-cancel="close" 
    has-modal-card>
    <div 
      class="modal-card" 
      style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">Login</p>
      </header>
      <section class="modal-card-body">
        <login-form ref="form" />
        <br>
        <span>
          Forgotten password? <a @click="passwordReset">click here</a>
        </span>
      </section>
      <footer class="modal-card-foot">
        <button 
          class="button" 
          type="button" 
          @click="close">Close</button>
        <button 
          class="button is-primary" 
          @click="submit">Login</button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
import LoginForm from './forms/LoginForm'
export default {
  components: {
    LoginForm
  },
  props: {
    opened: {
      default: false,
      type: Boolean
    },
    close: {
      type: Function,
      required: true
    }
  },
  methods: {
    passwordReset() {
      this.$dialog.prompt({
        message: "What's your email address?",
        inputAttrs: {
          placeholder: 'Type your email address'
        },
        onConfirm: email => {
          this.$http
            .post('/users/pass/reset', { email })
            .then(res => {
              this.$toast.open(res.data.data)
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                this.$toast.error(err.response.data.error)
              }
            })
        }
      })
    },
    submit() {
      this.$refs.form.submit()
    }
  }
}
</script>

<style>
</style>
