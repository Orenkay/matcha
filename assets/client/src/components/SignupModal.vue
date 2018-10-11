<template>
  <b-modal :active="opened" :onCancel="close" has-modal-card>
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">Register</p>
      </header>
      <section class="modal-card-body">
        <signup-form ref="form" />
      </section>
      <footer class="modal-card-foot">
        <button class="button" type="button" @click="close">Close</button>
        <button class="button is-primary" @click="submit">Register</button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
import SignupForm from "./forms/SignupForm";
export default {
  props: {
    opened: {
      default: false,
      type: Boolean
    },
    close: {
      type: Function,
      required: true
    }
  },
  components: {
    SignupForm
  },
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.$http
          .post("/auth/register", data, { errorHandle: false })
          .then(res => {
            this.$toast.open("Account successfuly created");
            this.close();
          })
          .catch(err => {
            if (err.response.status === 400) {
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
