<template>
  <m-form ref="form">
    <m-form-field 
      :required="true" 
      :props="{placeholder: 'Your Username'}" 
      name="user" 
      label="Username" 
      type="input" />
    <m-form-field 
      :required="true" 
      :props="{type: 'password', placeholder: 'Your Password'}" 
      name="pass" 
      label="Password" 
      type="input" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(data => {
        this.$http
          .post('/auth/login', data)
          .then(res => {
            this.$store.commit('login', res.data.token)
            this.$toast.open('Connected')
            this.close()
          })
          .catch(err => {
            if (err.response) {
              if (err.response.status === 400) {
                this.$toast.error(err.response.data.error)
              }
            }
          })
      })
    }
  }
}
</script>

<style>
</style>
