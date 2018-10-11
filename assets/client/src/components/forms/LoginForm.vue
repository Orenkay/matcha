<template>
  <m-form ref="form">
    <m-form-field name="user" label="Username" type="input" :required="true" :props="{placeholder: 'Your Username'}" />
    <m-form-field name="pass" label="Password" type="input" :required="true" :props="{type: 'password', placeholder: 'Your Password'}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(cb);
    },
    fieldError(k, errors) {
      this.$refs.form.fieldError(k, errors);
    }
  }
};
</script>

<style>
</style>
