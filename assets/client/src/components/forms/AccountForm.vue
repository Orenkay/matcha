<template>
  <m-form ref="form">
    <m-form-field 
      :value="account.username" 
      :required="true" 
      :props="{placeholder: 'Your Username'}" 
      name="user" 
      label="Username" 
      type="input" />
    <m-form-field 
      :value="account.email" 
      :required="true" 
      :props="{placeholder: 'Your Email'}" 
      name="email" 
      label="Email" 
      type="input" />
  </m-form>
</template>

<script>
export default {
  computed: {
    account() {
      return this.$store.getters.userData.account
    }
  },
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.passConfirm(pass => {
          data.currPass = pass
          this.$http
            .patch('/users/me/', data)
            .then(res => {
              this.$toast.open('Account informations successfuly edited')
              this.$store.commit('logout')
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                const { data } = err.response
                if (data.data && data.data.validation !== undefined) {
                  data.data.validation.keys.forEach((k, i) => {
                    this.$refs.form.fieldError(
                      k,
                      data.data.validation.details[i].message
                    )
                  })
                } else {
                  this.$toast.error(err.response.data.error)
                }
              }
            })
        })
      })
    },
    passConfirm(cb) {
      this.$dialog.prompt({
        message: `Type your current password`,
        inputAttrs: {
          type: 'password'
        },
        onConfirm: cb
      })
    }
  }
}
</script>

<style>
</style>
