<template>
  <m-form ref="form">
    <m-form-field name="pass" label="New Password" type="input" :required="true" :props="{type: 'password', placeholder: 'Your New Password'}" />
    <m-form-field name="pass2" label="New Password confirmation" type="input" :required="true" :props="{type: 'password', placeholder: 'Your New Password'}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(cb);
    }
  }
};
</script>

<style>
</style>
