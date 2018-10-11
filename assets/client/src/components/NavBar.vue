<template>
  <div class="container">
    <nav class="navbar has-shadow is-spaced">
      <div class="navbar-brand">
        <router-link to="/" class="navbar-item" href="/app">
          <h1 class="brand">Matcha</h1>
        </router-link>
        <div class="navbar-burger" :class="menuOpened && 'is-active'" @click="toggleMenu">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>

      <div class="navbar-menu" :class="menuOpened && 'is-active'">
        <div class="navbar-start">
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="field is-grouped">
              <div class="buttons">
                <button class="button is-link is-outlined" @click="profile">
                  <span class="icon">
                    <b-icon icon="account" size="is-small" />
                  </span>
                  <span>Profile</span>
                </button>
                <button class="button is-black is-outlined">
                  <span class="icon">
                    <b-icon icon="bell" size="is-small" />
                  </span>
                </button>
                <b-dropdown position="is-bottom-left">
                  <button class="button is-black is-outlined" slot="trigger">
                    <span class="icon">
                      <b-icon icon="settings" size="is-small" />
                    </span>
                  </button>
                  <b-dropdown-item @click="accountEdit">
                    <b-icon icon="settings" size="is-small" />
                    <span>Account</span>
                  </b-dropdown-item>
                  <b-dropdown-item @click="toggleNotifications">
                    <div :class="notificationsEnabled ? 'has-text-link' : 'has-text-danger'">
                      <b-icon :icon="notificationsEnabled ? 'check' : 'close'" size="is-small" />
                      <span>Notifications</span>
                    </div>
                  </b-dropdown-item>
                  <b-dropdown-item @click="logout">
                    <div>
                      <b-icon icon="exit-to-app" size="is-small" />
                      <span>Logout</span>
                    </div>
                  </b-dropdown-item>
                </b-dropdown>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  </div>
</template>

<script>
export default {
  data() {
    return {
      menuOpened: false
    };
  },
  computed: {
    notificationsEnabled() {
      return this.$store.state.notificationsEnabled;
    }
  },
  methods: {
    toggleMenu() {
      this.menuOpened = !this.menuOpened;
    },
    profile() {
      this.$router.push("/app/profile");
    },
    accountEdit() {
      this.$router.push("/app/account/edit");
    },
    toggleNotifications() {
      this.$store.commit("toggle-notifications");
      this.$toast.open({
        queue: false,
        message: `Notifications ${
          this.notificationsEnabled ? "enabled" : "disabled"
        }`
      });
    },
    logout() {
      this.$store.dispatch("logout").catch(() => {
        this.$toast.open("Unable to logout, please retry");
      });
    }
  }
};
</script>

<style scoped>
.is-spaced {
  padding-left: 0;
  padding-right: 0;
}
</style>
