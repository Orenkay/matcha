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
      this.$refs.form.submit(data => {
        this.$http
          .post("/auth/login", data)
          .then(res => {
            this.$store.commit("login", res.data.token);
            this.$toast.open("Connected");
            this.close();
          })
          .catch(err => {
            if (err.response) {
              if (err.response.status === 400) {
                this.$toast.error(err.response.data.error);
              }
            }
          });
      });
    }
  }
};
</script>

<style>
</style>
