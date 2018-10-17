<template>
  <m-form ref="form">
    <m-form-field name="user" label="Username" type="input" :required="true" :props="{placeholder: 'Your Username'}" />
    <m-form-field name="email" label="Email" type="input" :required="true" :props="{placeholder: 'Your Email'}" />
    <m-form-field name="pass" label="Password" type="input" :required="true" :props="{type: 'password', placeholder: 'Your Password'}" />
    <m-form-field name="pass2" label="Re-Password" type="input" :required="true" :props="{type: 'password', placeholder: 'Re-type Your Password'}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(data => {
        this.$http
          .post("/auth/register", data)
          .then(res => {
            this.$toast.open("Account successfuly created");
            cb();
          })
          .catch(err => {
            if (err.response && err.response.status === 400) {
              const { data } = err.response.data;
              if (data.validation !== undefined) {
                data.validation.keys.forEach((k, i) => {
                  this.$refs.form.fieldError(
                    k,
                    data.validation.details[i].message
                  );
                });
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
