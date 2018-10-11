<template>
  <div class="container">
    <div class="middle-container">
      <figure class="image is-128x128">
        <!-- <figure class="image is-128x128 avatar-overlay" @click="avatarEdit">
          <b-icon icon="pencil" />
        </figure> -->
        <img class="is-rounded" src="https://bulma.io/images/placeholders/128x128.png" />
      </figure>
      <router-link to="/app/profile/edit">edit</router-link>
      <div class="info-name">
        <span class="is-inline">{{profile.lastName}}</span>
        <span class="is-inline">{{profile.firstName}}</span>
        <b-icon pack="fas" :icon="profile.gender === 'male' ? 'mars' : 'venus'" size="is-small" :class="'is-' + profile.gender" />
      </div>
      <div>
        <b-icon pack="fas" icon="map-marker-alt" size="is-small" />
        <span class="is-inline">{{loc.address}}</span>
      </div>
      <div>
        <span class="is-inline">Attracted by:</span>
        <b-icon v-for="gender in attractedBy" :key="gender" pack="fas" :icon="gender" size="is-small" :class="'is-' + gender" />
      </div>
      <div class="info-bio">
        {{profile.bio}}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  methods: {
    avatarEdit() {}
  },
  computed: {
    profile() {
      return this.$store.getters.profile;
    },
    loc() {
      return this.$store.getters.loc;
    },
    attractedBy() {
      const gender = this.$store.getters.profile.gender;
      switch (this.$store.getters.profile.attraction) {
        case "hetero":
          return gender === "male" ? ["female"] : ["male"];
        case "homo":
          return gender === "male" ? ["male"] : ["female"];
        default:
          return ["male", "female"];
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.tooltip {
  cursor: pointer;
}

.middle-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;

  & > * {
    margin-bottom: 20px;
  }
}
.is-male {
  color: #15aabf;
}
.is-female {
  color: #be4bdb;
}
.info-name {
  font-size: 20px;
}
.info-bio {
  font-size: 14px;
  text-align: center;
}
.avatar-overlay {
  position: absolute;

  display: flex;
  justify-content: center;
  align-items: center;

  background: rgba(0, 0, 0, 0.8);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  opacity: 0;

  &:hover {
    opacity: 1;
  }
}
</style>
