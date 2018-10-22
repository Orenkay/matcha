<template>
  <div id="app">
    <b-loading 
      :is-full-page="true" 
      :active="!isLoaded" />
    <error-reporter 
      v-if="globalError" 
      :error="globalError" />
    <router-view />
  </div>
</template>

<script>
import ErrorReporter from './components/ErrorReporter'

export default {
  name: 'App',
  components: {
    ErrorReporter
  },
  computed: {
    isLoaded() {
      return this.$store.state.loaded
    },
    globalError() {
      return this.$store.state.error
    }
  }
}
</script>

<style lang="scss">
html,
body,
#app {
  height: 100%;
  width: 100%;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  display: flex;
  justify-content: center;
  color: #2c3e50;
}

.brand {
  font-size: 32px;
  font-weight: bold;
}

h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
