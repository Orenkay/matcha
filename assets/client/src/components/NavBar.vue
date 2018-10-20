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
          <router-link class="navbar-item" to="/app/search">
            <span class="icon">
              <b-icon icon="magnify" size="is-small" />
            </span>
            <span>Search</span>
          </router-link>
          <router-link class="navbar-item" to="/app/history">
            <span class="icon">
              <b-icon icon="history" size="is-small" />
            </span>
            <span>History</span>
          </router-link>
          <router-link class="navbar-item" to="/app/messages">
            <span class="icon">
              <b-icon icon="message" size="is-small" />
            </span>
            <span>Messages</span>
          </router-link>
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="field is-grouped">
              <div class="buttons">
                <router-link class="navbar-item" to="/app/profile/me">
                  <button class="button is-link is-outlined">
                    <span class="icon">
                      <b-icon icon="account" size="is-small" />
                    </span>
                    <span>Profile</span>
                  </button>
                </router-link>
                <notifications-dropdown />
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
import NotificationsDropdown from "./NotificationsDropdown";
export default {
  components: {
    NotificationsDropdown
  },
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
    accountEdit() {
      this.$router.push("/app/account/edit");
    },
    logout() {
      this.$store.dispatch("logout").catch(() => {
        this.$router.go();
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
