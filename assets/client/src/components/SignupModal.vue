<template>
  <b-modal 
    :active="opened" 
    :on-cancel="close" 
    has-modal-card>
    <div 
      class="modal-card" 
      style="width: auto">
      <header class="modal-card-head">
        <p class="modal-card-title">Register</p>
      </header>
      <section class="modal-card-body">
        <signup-form ref="form" />
      </section>
      <footer class="modal-card-foot">
        <button 
          class="button" 
          type="button" 
          @click="close">Close</button>
        <button 
          class="button is-primary" 
          @click="submit">Register</button>
      </footer>
    </div>
  </b-modal>
</template>

<script>
import SignupForm from './forms/SignupForm'
export default {
  components: {
    SignupForm
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
      this.$refs.form.submit(() => this.close())
    }
  }
}
</script>

<style>
</style>
