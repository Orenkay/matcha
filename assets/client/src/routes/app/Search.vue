<template>
  <div>
    <b-collapse class="filter-collapse" :open="true">
      <div class="notification">
        <div class="content">
          <filters-form @apply="filtersApply">
            <span slot="button-label">Search</span>
          </filters-form>
        </div>
      </div>
    </b-collapse>
    <br />
    <router-link v-for="(s, index) in searches" :key="index" :to="`/app/profile/${s.profile.userId}`">
      <article class="media">
        <figure class="media-left">
          <div class="image is-64x64" :style="`background-image: url('${s.pp}')`"></div>
        </figure>
        <div class="media-content">
          <div class="content">
            <p>
              <span>
                <span>{{s.profile.lastName }}</span>
                <span>{{ s.profile.firstName}}</span>
                <span class="meta">
                  <span>{{ s.profile.age }} yo</span>
                  <b-tooltip label="popularity" type="is-black">
                    <span class="popularity">
                      <b-icon icon="star" size="is-small" />
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
import moment from "moment";
import UserInterests from "../../components/UserInterests";
import UserLocation from "../../components/UserLocation";
import FiltersForm from "../../components/forms/FiltersForm";
export default {
  components: {
    UserInterests,
    UserLocation,
    FiltersForm
  },
  data() {
    return {
      searches: []
    };
  },
  methods: {
    filtersApply(params) {
      this.fetchData(params);
    },
    fetchData(params) {
      params = Object.assign(params || {}, { suggestion: false });
      this.$http.get("/matcher/1", { params }).then(res => {
        if (!res.data.data) {
          this.searches = [];
          return;
        }
        Promise.all(
          res.data.data.map(
            (targetId, i) =>
              new Promise((resolve, reject) => {
                this.$store.dispatch("getProfile", targetId).then(data => {
                  data.pp = data.pictures.find(p => p.isPP).path;
                  data.profile.age = moment().diff(
                    data.profile.birthdate * 1000,
                    "years",
                    false
                  );
                  resolve(data);
                });
              })
          )
        ).then(searches => {
          this.searches = searches;
        });
      });
    },
    goto(id) {
      this.$router.push("/app/profile/" + id);
    }
  }
};
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
