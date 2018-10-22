<template>
  <m-form ref="form">
    <m-form-field 
      :value="profile.firstName" 
      :required="true" 
      :props="{placeholder: 'Your Firstname'}" 
      name="firstName" 
      label="Firstname" 
      type="input" />
    <m-form-field 
      :value="profile.lastName" 
      :required="true" 
      :props="{placeholder: 'Your Lastname'}" 
      name="lastName" 
      label="Lastname" 
      type="input" />
    <m-form-field 
      :value="birthdate" 
      :required="true" 
      :props="{'min-date': minDate, 'max-date': maxDate, placeholder: 'Your birthdate', }" 
      name="birthdate" 
      label="Birthdate" 
      type="datepicker" />
    <m-form-field 
      :value="profile.gender" 
      :required="true" 
      :props="{placeholder: 'Your Gender'}" 
      name="gender" 
      type="select" 
      label="Gender">
      <option value="male">Male</option>
      <option value="female">Female</option>
    </m-form-field>
    <m-form-field 
      :value="profile.attraction" 
      :required="true" 
      :props="{placeholder: 'Your Sexual Orientation'}" 
      name="attraction" 
      type="select" 
      label="Sexual Orientation">
      <option value="hetero">Heterosexual</option>
      <option value="bi">Bisexual</option>
      <option value="homo">Homosexual</option>
    </m-form-field>
    <m-form-field 
      :value="profile.bio" 
      :required="true" 
      :props="{type: 'textarea', placeholder: 'Describe you in few lines', maxlength: 200}" 
      name="bio" 
      label="Bio" 
      type="input" />
  </m-form>
</template>

<script>
export default {
  data() {
    const maxDate = new Date()
    maxDate.setFullYear(maxDate.getFullYear() - 18)
    return {
      minDate: new Date((1 << 31) * 1000),
      maxDate: maxDate
    }
  },
  computed: {
    birthdate() {
      return this.profile.birthdate && new Date(this.profile.birthdate * 1000)
    },
    profile() {
      return this.$store.getters.userData.profile
    }
  },
  methods: {
    submit(cb) {
      this.$refs.form.submit(data => {
        data.birthdate = data.birthdate.getTime() / 1000
        const req =
          this.$store.getters.profile.lastName === undefined
            ? this.$http.post('/profiles/me', data)
            : this.$http.patch('/profiles/me/edit', data)
        req.then(res => {
          this.$store.commit('setUserData', ['profile', res.data.data])
          if (typeof cb === 'function') {
            cb()
          }
        })
        req.catch(err => {
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
