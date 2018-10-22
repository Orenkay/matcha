<template>
  <div>
    <b-tabs>
      <b-tab-item label="Basic">
        <account-form ref="form" />
        <br >
        <button 
          class="button is-black" 
          @click="submit">Edit</button>
      </b-tab-item>

      <b-tab-item label="Password">
        <pass-edit-form ref="formPass" />
        <br >
        <button 
          class="button is-black" 
          @click="submitPass">Edit</button>
      </b-tab-item>
    </b-tabs>
    <button 
      class="button is-danger" 
      @click="deleteAccount">Delete account</button>
  </div>
</template>

<script>
import AccountForm from '../../components/forms/AccountForm'
import PassEditForm from '../../components/forms/PassEditForm'
export default {
  components: {
    AccountForm,
    PassEditForm
  },
  methods: {
    submit() {
      this.$refs.form.submit()
    },
    submitPass() {
      this.$refs.formPass.submit()
    },
    deleteAccount() {
      this.$dialog.prompt({
        message: `Type your current password`,
        inputAttrs: {
          type: 'password'
        },
        onConfirm: pass => {
          this.$http
            .delete('/users/me/' + pass)
            .then(res => {
              this.$toast.open('Account deleted')
              this.$store.commit('logout')
            })
            .catch(err => {
              this.$toast.error(err.response.data.error)
            })
        }
      })
    }
  }
}
</script>
