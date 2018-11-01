<template>
  <m-form ref="form">
    <m-form-field 
      :required="true" 
      :props="{type: 'password', placeholder: 'Your New Password'}" 
      :validate="validate"
      name="pass" 
      label="New Password" 
      type="input" />
    <m-form-field 
      :required="true" 
      :props="{type: 'password', placeholder: 'Your New Password'}" 
      :validate="validate"
      name="pass2" 
      label="New Password confirmation" 
      type="input" />
  </m-form>
</template>

<script>
export default {
  methods: {
    validate(v, form) {
      if (form.pass !== form.pass2) {
        return ['password doesnt match']
      }
      return []
    },
    submit() {
      this.$refs.form.submit(data => {
        this.passConfirm(pass => {
          data.currPass = pass
          this.$http
            .patch('/users/me/password', data)
            .then(res => {
              this.$toast.open('Password changed')
              this.$store.commit('logout')
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                this.$toast.error(err.response.data.error)
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
