<template>
  <div class="app-container">
    <navbar></navbar>
    <div class="app-body">
      <div class="container">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script>
import Navbar from "../components/NavBar";
export default {
  components: {
    Navbar
  }
};
</script>

<style scoped>
#error {
  display: flex;
  text-align: center;
  align-items: center;
  justify-content: center;
}

.app-container {
  width: 100%;
}
.app-body {
  margin-top: 2rem;
}
</style>
