<template>
  <m-form ref="form">
    <m-form-field name="firstName" label="Firstname" :value="profile.firstName" type="input" :required="true" :props="{placeholder: 'Your Firstname'}" />
    <m-form-field name="lastName" label="Lastname" type="input" :value="profile.lastName" :required="true" :props="{placeholder: 'Your Lastname'}" />
    <m-form-field name="Gender" type="select" label="Gender" :value="profile.gender" :required="true" :props="{placeholder: 'Your Gender'}">
      <option value="male">Male</option>
      <option value="female">Female</option>
    </m-form-field>
    <m-form-field name="attraction" type="select" label="Sexual Orientation" :value="profile.attraction" :required="true" :props="{placeholder: 'Your Sexual Orientation'}">
      <option value="hetero">Heterosexual</option>
      <option value="bi">Bisexual</option>
      <option value="homo">Homosexual</option>
    </m-form-field>
    <m-form-field name="bio" label="Bio" type="input" :value="profile.bio" :required="true" :props="{type: 'textarea', placeholder: 'Describe you in few lines', maxlength: 200}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(cb);
    }
  },
  computed: {
    profile() {
      return this.$store.getters.userData.profile;
    }
  }
};
</script>

<style>
</style>
