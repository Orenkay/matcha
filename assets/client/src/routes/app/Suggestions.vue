<template>
  <div>
    <b-collapse 
      :open="false" 
      class="filter-collapse">
      <div 
        slot="trigger" 
        slot-scope="props" 
        class="card-header">
        <p class="card-header-title">
          Advanced settings
        </p>
        <a class="card-header-icon">
          <b-icon 
            :icon="props.open ? 'menu-down' : 'menu-up'" 
            type="is-black"/>
        </a>
      </div>
      <div class="notification">
        <div class="content">
          <filters-form @apply="filtersApply" />
        </div>
      </div>
    </b-collapse>
    <br >
    <router-link 
      v-for="(s, index) in suggestions" 
      :key="index" 
      :to="`/app/profile/${s.profile.userId}`">
      <article class="media">
        <figure class="media-left">
          <div 
            :style="`background-image: url('${s.pp}')`" 
            class="image is-64x64"/>
        </figure>
        <div class="media-content">
          <div class="content">
            <p>
              <span>
                <span>{{ s.profile.lastName }}</span>
                <span>{{ s.profile.firstName }}</span>
                <span class="meta">
                  <span>{{ s.profile.age }} yo</span>
                  <b-tooltip 
                    label="popularity" 
                    type="is-black">
                    <span class="popularity">
                      <b-icon 
                        icon="star" 
                        size="is-small" />
                      <span>{{ Math.trunc(100 * s.popularity) }}%</span>
                    </span>
                  </b-tooltip>
                </span>
              </span>
              <user-location :location="s.loc" />
              <user-interests :tags="s.interests" />
            </p>
          </div>
        </div>
      </article>
    </router-link>
  </div>
</template>

<script>
import moment from 'moment'
import UserInterests from '../../components/UserInterests'
import UserLocation from '../../components/UserLocation'
import FiltersForm from '../../components/forms/FiltersForm'
export default {
  components: {
    UserInterests,
    UserLocation,
    FiltersForm
  },
  data() {
    return {
      suggestions: []
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    filtersApply(params) {
      this.fetchData(params)
    },
    fetchData(params) {
      params = Object.assign(params || {}, { suggestion: true })
      this.$http.get('/matcher/0', { params }).then(res => {
        if (!res.data.data) {
          this.suggestions = []
          return
        }
        Promise.all(
          res.data.data.map(
            (targetId, i) =>
              new Promise((resolve, reject) => {
                this.$store.dispatch('getProfile', targetId).then(data => {
                  data.pp = data.pictures.find(p => p.isPP).path
                  data.profile.age = moment().diff(
                    data.profile.birthdate * 1000,
                    'years',
                    false
                  )
                  resolve(data)
                })
              })
          )
        ).then(suggestions => {
          this.suggestions = suggestions
        })
      })
    },
    goto(id) {
      this.$router.push('/app/profile/' + id)
    }
  }
}
</script>

<style scoped lang="scss">
.filter-collapse {
  user-select: none;
}
article.media {
  color: initial;
  margin-top: 0;
  padding: 10px 0;
  align-items: center;
}
article.media:hover {
  background: rgba(0, 0, 0, 0.03);
  cursor: pointer;
}
.image {
  background-position: 50% 50%;
  background-size: 100%;
  border-radius: 32px;
}
.meta {
  font-size: 12px;
  margin-left: 5px;
  color: grey;
  > * {
    margin: 0 2px;
  }
}
.popularity,
.is-small {
  font-size: 11px !important;
  &:hover {
    color: black;
  }
}

.more-params {
  width: 100%;
}
</style>
