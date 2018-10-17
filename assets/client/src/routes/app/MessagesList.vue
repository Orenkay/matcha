<template>
  <div>
    <div v-for="m in matches" :key="m">
      <message-user :target="m" @click.native="show(m)" />
    </div>
  </div>
</template>

<script>
import MessageUser from "../../components/MessageUser";
export default {
  components: {
    MessageUser
  },
  computed: {
    matches() {
      return this.$store.state.userData.matches;
    }
  },
  methods: {
    show(id) {
      this.$router.push("/app/messages/" + id);
    }
  }
};
</script>
