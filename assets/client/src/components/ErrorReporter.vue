<template>
  <div class="overlay">
    <section class="section">
      <h1 class="title">Oops something went wrong</h1>
      <span class="error-type">{{ errorType() }}</span>
      <h2 class="subtitle">
        {{ displayError() }}
      </h2>
      <button class="button" @click="retry">Click here to reload the page</button>
    </section>
  </div>
</template>

<script>
export default {
  props: ["error"],
  methods: {
    retry() {
      this.$router.go();
    },
    displayError() {
      if (RegExp("Network Error").test(this.error)) {
        return "Request response is empty, please check your connection internet";
      }
      if (this.error.response !== undefined) {
        if (this.error.response.status === 401) {
          return `You are not authorized to access this place`;
        }
        return `Unknwon error. Code: ${this.error.response.status}`;
      }
      return `Unknown error.`;
    },
    errorType() {
      if (RegExp("Network Error").test(this.error)) {
        return "Network Error";
      }
      if (this.error.response !== undefined) {
        return "Request Error";
      }
      return "Error";
    }
  }
};
</script>

<style scoped>
.overlay {
  z-index: 100;
  background: rgba(255, 255, 255, 0.98);
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  top: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.error-type {
  font-size: 12px;
}
</style>
