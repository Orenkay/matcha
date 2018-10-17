<template>
  <m-form ref="form">
    <m-form-field name="pass" label="New Password" type="input" :required="true" :props="{type: 'password', placeholder: 'Your New Password'}" />
    <m-form-field name="pass2" label="New Password confirmation" type="input" :required="true" :props="{type: 'password', placeholder: 'Your New Password'}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.passConfirm(pass => {
          data.currPass = pass;
          this.$http
            .patch("/users/me/password", data)
            .then(res => {
              this.$toast.open("Password changed");
              this.$store.commit("logout");
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                this.$toast.error(err.response.data.error);
              }
            });
        });
      });
    },
    passConfirm(cb) {
      this.$dialog.prompt({
        message: `Type your current password`,
        inputAttrs: {
          type: "password"
        },
        onConfirm: cb
      });
    }
  }
};
</script>

<style>
</style>
