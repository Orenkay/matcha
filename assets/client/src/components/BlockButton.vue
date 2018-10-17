<template>
  <b-tooltip :label="tooltip" type="is-black" position="is-bottom">
    <button class="button" :class="classes" @click="click" :disabled="loading">
      <b-icon icon="block-helper" size="is-small" />
    </button>
  </b-tooltip>
</template>

<script>
export default {
  props: ["target"],
  created() {
    this.fetchData();
  },
  data() {
    return {
      isBlocked: false,
      loading: false
    };
  },
  computed: {
    selfData() {
      return this.$store.getters.userData;
    },
    tooltip() {
      return this.isBlocked ? "Un-Block" : "Block";
    },
    classes() {
      return [this.isBlocked && "is-danger"];
    }
  },
  methods: {
    fetchData() {
      this.loading = true;
      this.$http
        .get(`/block/${this.target}/blocked`)
        .then(res => {
          this.isBlocked = res.data.data;
        })
        .finally(() => {
          this.loading = false;
        });
    },
    click() {
      this.loading = true;
      const req = !this.isBlocked
        ? this.$http.post(`/block/${this.target}`)
        : this.$http.delete(`/block/${this.target}`);
      req
        .then(() => {
          this.isBlocked = !this.isBlocked;
          this.$store.dispatch("removeMessageFrom", this.target);
        })
        .finally(() => {
          this.loading = false;
        });
    }
  }
};
</script>

<style scoped>
.middle-container {
  padding-top: 0;
}
</style>
