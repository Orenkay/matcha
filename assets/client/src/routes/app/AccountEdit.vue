<template>
  <div>
    <b-tabs>
      <b-tab-item label="Basic">
        <account-form ref="form" />
        <br />
        <button class="button is-black" @click="submit">Edit</button>
      </b-tab-item>

      <b-tab-item label="Password">
        <pass-edit-form ref="formPass" />
        <br />
        <button class="button is-black" @click="submitPass">Edit</button>
      </b-tab-item>
    </b-tabs>
    <button class="button is-danger" @click="deleteAccount">Delete account</button>
  </div>
</template>

<script>
import AccountForm from "../../components/forms/AccountForm";
import PassEditForm from "../../components/forms/PassEditForm";
export default {
  components: {
    AccountForm,
    PassEditForm
  },
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.passConfirm(pass => {
          data.currPass = pass;
          this.$http
            .patch("/users/me/", data, { errorHandle: false })
            .then(res => {
              this.$toast.open("Account informations successfuly edited");
              this.$store.commit("logout");
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                const { data } = err.response;
                if (data.data && data.data.validation !== undefined) {
                  data.data.validation.keys.forEach((k, i) => {
                    this.$refs.form.fieldError(
                      k,
                      data.data.validation.details[i].message
                    );
                  });
                } else {
                  this.$toast.open({
                    message: data.error,
                    type: "is-danger"
                  });
                }
              }
            });
        });
      });
    },
    submitPass() {
      this.$refs.formPass.submit(data => {
        console.lo;
        this.passConfirm(pass => {
          data.currPass = pass;
          this.$http
            .patch("/users/me/password", data, { errorHandle: false })
            .then(res => {
              this.$toast.open("Password changed");
              this.$store.commit("logout");
            })
            .catch(err => {
              if (err.response && err.response.status === 400) {
                this.$toast.open({
                  message: err.response.data.error,
                  type: "is-danger"
                });
              }
            });
        });
      });
    },
    deleteAccount() {
      this.$dialog.prompt({
        message: `Type your current password`,
        inputAttrs: {
          type: "password"
        },
        onConfirm: pass => {
          this.$http.delete("/users/me/" + pass, { errorHandle: false}).then(res => {
            this.$toast.open("Account deleted");
            this.$store.commit("logout");
          }).catch(err => {
            this.$toast.open(err.response.data.error)
          });
        }
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
