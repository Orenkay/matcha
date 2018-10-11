<template>
  <b-modal :active="opened" :onCancel="close" has-modal-card>
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">Login</p>
      </header>
      <section class="modal-card-body">
        <login-form ref="form" />
      </section>
      <footer class="modal-card-foot">
        <button class="button" type="button" @click="close">Close</button>
        <button class="button is-primary" @click="submit">Login</button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
import LoginForm from "./forms/LoginForm";
export default {
  components: {
    LoginForm
  },
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
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.$http
          .post("/auth/login", data, { errorHandle: false })
          .then(res => {
            this.$store.commit("login", res.data.token);
            this.$toast.open("Connected");
            this.close();
          })
          .catch(err => {
            if (err.response) {
              if (err.response.status === 400) {
                this.$toast.open({
                  message: err.response.data.error,
                  type: "is-danger"
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
