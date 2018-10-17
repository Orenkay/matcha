<template>
  <div v-if="this.presence">
    {{ presenceText }}
  </div>
</template>

<script>
import moment from "moment";
export default {
  props: ["target"],
  created() {
    this.getPresence();
  },
  data() {
    return {
      presence: undefined
    };
  },
  computed: {
    isOnline() {
      return Date.now() - this.presence <= 1000 * 60 * 5;
    },
    presenceText() {
      if (!this.isOnline) {
        return (
          "offline since: " + moment(this.presence).format("DD/MM/YYYY - HH:mm")
        );
      }
      return "online";
    }
  },
  methods: {
    getPresence() {
      this.$http.get(`/profiles/${this.target}/presence`).then(res => {
        this.presence = new Date(res.data.data * 1000);
      });
    }
  }
};
</script>
