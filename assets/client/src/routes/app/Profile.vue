<template>
  <div v-if="loaded">
    <div v-if="notFound">
      <div class="middle-container">
        <h1 class="title">This user profile does not exists :(</h1>
      </div>
    </div>
    <div 
      v-else 
      class="container">
      <div class="middle-container">
        <user-pictures 
          :only-pp="true" 
          :pictures="userData.pictures" />

        <!-- Here we show the edit link if we own the profile -->
        <router-link 
          v-if="owner" 
          to="/app/profile/me/edit/">edit</router-link>

        <div v-if="!owner">
          <like-button :target="targetId" />
          <fake-button :target="targetId" />
          <block-button :target="targetId" />
        </div>

        <div class="info-name">
          <span class="is-inline">{{ profile.lastName }}</span>
          <span class="is-inline">{{ profile.firstName }}</span>
          <b-icon 
            :icon="profile.gender === 'male' ? 'mars' : 'venus'" 
            :class="'is-' + profile.gender" 
            pack="fas" 
            size="is-small" />
          <b-tooltip 
            label="popularity" 
            type="is-black">
            <span>
              <b-icon 
                icon="star" 
                size="is-small" />  
              <span>
                {{ Math.trunc(userData.popularity * 100) }}%
              </span>
            </span>
          </b-tooltip>
        </div>

        <user-presence 
          v-if="!owner" 
          :target="$route.params.id" />

        <user-interests 
          :tags="userData.interests" 
          type="is-centered" />

        <user-location :location="userData.loc" />

        Age: {{ age }}

        <div>
          <span class="is-inline">Attracted by:</span>
          <b-icon 
            v-for="gender in attractedBy" 
            :key="gender" 
            :icon="gender" 
            :class="'is-' + gender" 
            pack="fas" 
            size="is-small" />
        </div>

        <div class="info-bio">
          {{ profile.bio }}
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import moment from 'moment'

import UserLocation from '../../components/UserLocation'
import UserPictures from '../../components/UserPictures'
import UserInterests from '../../components/UserInterests'
import UserPresence from '../../components/UserPresence'

import LikeButton from '../../components/LikeButton'
import FakeButton from '../../components/FakeButton'
import BlockButton from '../../components/BlockButton'

export default {
  components: {
    UserLocation,
    UserPictures,
    UserInterests,
    UserPresence,

    LikeButton,
    FakeButton,
    BlockButton
  },
  data() {
    return {
      userData: undefined
    }
  },
  computed: {
    selfData() {
      return this.$store.getters.userData
    },
    targetId() {
      return parseInt(this.$route.params.id)
    },
    notFound() {
      if (this.userData === undefined) {
        return true
      }
    },
    owner() {
      const id = this.$route.params.id
      return id === 'me' || parseInt(id) === this.selfData.account.id
    },
    age() {
      return moment().diff(this.profile.birthdate * 1000, 'years', false)
    },
    loaded() {
      return this.$store.getters.loaded
    },
    profile() {
      return this.userData.profile
    },
    attractedBy() {
      const gender = this.userData.profile.gender
      switch (this.userData.profile.attraction) {
        case 'hetero':
          return gender === 'male' ? ['female'] : ['male']
        case 'homo':
          return gender === 'male' ? ['male'] : ['female']
        default:
          return ['male', 'female']
      }
    }
  },
  created() {
    if (!this.owner) {
      return this.fetchProfile(this.$route.params.id)
    }
    this.userData = this.selfData
  },
  methods: {
    fetchProfile(id, cb) {
      this.$store.commit('loading', true)
      this.$http.get(`/profiles/${id}/visit`)
      this.$store
        .dispatch('getProfile', id)
        .then(profile => {
          this.userData = profile
        })
        .finally(() => {
          if (typeof cb === 'function') {
            cb()
          }
          this.$store.commit('loading', false)
        })
    }
  },
  beforeRouteUpdate(to, from, next) {
    if (to.params.id === 'me') {
      this.userData = this.selfData
      return next()
    }
    this.fetchProfile(to.params.id, next)
  }
}
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
</style>
