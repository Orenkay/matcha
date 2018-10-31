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
      :props="{placeholder: 'Your Email'}" 
      name="email" 
      label="Email" 
      type="input" />
    <m-form-field 
      :required="true" 
      :props="{type: 'password', placeholder: 'Your Password'}" 
      :validate="passwordMatch"
      name="pass" 
      label="Password" 
      type="input" />
    <m-form-field 
      :required="true" 
      :props="{type: 'password', placeholder: 'Re-type Your Password'}" 
      :validate="passwordMatch"
      name="pass2" 
      label="Re-Password" 
      type="input" />
  </m-form>
</template>

<script>
export default {
  methods: {
    passwordMatch(v, form) {
      if (form.pass !== form.pass2) {
        return ['password doesnt match']
      }
      return []
    },
    submit(cb) {
      this.$refs.form.submit(data => {
        this.$http
          .post('/auth/register', data)
          .then(res => {
            this.$toast.open('Account successfuly created')
            cb()
          })
          .catch(err => {
            if (err.response && err.response.status === 400) {
              const { data } = err.response.data
              if (data.validation !== undefined) {
                data.validation.keys.forEach((k, i) => {
                  this.$refs.form.fieldError(
                    k,
                    data.validation.details[i].message
                  )
                })
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
