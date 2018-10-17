<template>
  <b-dropdown position="is-bottom-left">
    <div slot="trigger">
      <span class="bulle" v-if="displayBulle">{{count}}</span>
      <button class="button is-black is-outlined">
        <span class="icon">
          <b-icon icon="bell" size="is-small" />
        </span>
      </button>
    </div>
    <b-dropdown-item @click="clear">
      <b-icon icon="delete" size="is-small" />
      <span>Clear all notifications</span>
    </b-dropdown-item>
    <b-dropdown-item v-for="(notif, index) in notifications" :key="index" @click="trigger(index, notif)">
      <span>{{ notif.message }}</span>
    </b-dropdown-item>
  </b-dropdown>
</template>

<script>
export default {
  computed: {
    notifications() {
      return this.$store.getters.notifications;
    },
    displayBulle() {
      return this.notifications.length > 0;
    },
    count() {
      const count = this.notifications.length;
      return count < 10 ? count : "9+";
    }
  },
  methods: {
    trigger(id, n) {
      this.$store.commit("removeNotification", id);
      if (n.to) {
        this.$router.push(n.to);
      }
    },
    clear() {
      this.$store.commit("clearNotifications");
    }
  }
};
</script>

<style scoped>
.bulle {
  position: absolute;
  background: red;
  width: 20px;
  height: 20px;
  border-radius: 30px;
  z-index: 10;
  transform: translate(20px, 20px);
  color: white;
  display: flex;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  justify-content: center;
  align-items: center;
}
</style>
