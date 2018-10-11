<template>
  <m-form ref="form">
    <m-form-field name="user" label="Username" type="input" :value="account.username" :required="true" :props="{placeholder: 'Your Username'}" />
    <m-form-field name="email" label="Email" type="input" :value="account.email" :required="true" :props="{placeholder: 'Your Email'}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(cb);
    },
    fieldError(k, v) {
      this.$refs.form.fieldError(k, v)
    }
  },
  computed: {
    account() {
      return this.$store.getters.userData.account;
    }
  }
};
</script>

<style>
</style>
