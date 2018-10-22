<template>
  <section class="section">
    <div class="brand">Matcha</div>
    <section class="section">
      <login-modal 
        :opened="loginOpened" 
        :close="() => loginOpened = false"/>
      <button 
        class="button" 
        @click="() => loginOpened = true">Login</button>

      <signup-modal 
        :opened="signupOpened" 
        :close="() => signupOpened = false"/>
      <button 
        class="button" 
        @click="() => signupOpened = true">Signup</button>
    </section>
  </section>
</template>

<script>
import LoginModal from '../components/LoginModal'
import SignupModal from '../components/SignupModal'
export default {
  components: {
    LoginModal,
    SignupModal
  },
  data() {
    return {
      loginOpened: false,
      signupOpened: false
    }
  }
}
</script>

<style scoped>
.brand {
  text-align: center;
}
</style>
