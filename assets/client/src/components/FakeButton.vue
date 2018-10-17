<template>
  <b-tooltip :label="tooltip" type="is-black" position="is-bottom">
    <button class="button" @click="click" :disabled="loading || isReported">
      <b-icon icon="alert" size="is-small" />
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
      isReported: false,
      loading: false
    };
  },
  computed: {
    selfData() {
      return this.$store.getters.userData;
    },
    tooltip() {
      return this.isReported ? "Already reported" : "Report";
    }
  },
  methods: {
    fetchData() {
      this.loading = true;
      this.$http
        .get(`/report/${this.target}/reported`)
        .then(res => {
          this.isReported = res.data.data;
        })
        .finally(() => {
          this.loading = false;
        });
    },
    click() {
      this.loading = true;
      this.$http
        .post(`/report/${this.target}`)
        .then(() => {
          this.isReported = true;
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
