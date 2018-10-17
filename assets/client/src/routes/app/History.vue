<template>
  <div>
    <article class="media" v-for="(h, index) in history" :key="index" @click="goto(h.fromId)">
      <figure class="media-left">
        <div class="image is-32x32" :style="`background-image: url('${pp(h)}')`"></div>
      </figure>
      <div class="media-content">
        <div class="content">
          <p>
            {{ h.profile.lastName }} {{ h.profile.firstName}} <small>visited your profile</small>
          </p>
        </div>
      </div>
    </article>
  </div>
</template>

<script>
export default {
  created() {
    this.fetchData();
  },
  data() {
    return {
      history: []
    };
  },
  methods: {
    fetchData() {
      this.$http.get("/users/me/history").then(res => {
        const history = res.data.data || [];
        history.forEach(h => {
          this.$store.dispatch("getProfile", h.fromId).then(data => {
            data.fromId = h.fromId;
            this.history.push(data);
          });
        });
      });
    },
    goto(to) {
      this.$router.push("/app/profile/" + to);
    },
    pp(h) {
      const pp = h.pictures.find(p => p.isPP);
      if (!pp) {
        return "";
      }
      return pp.path;
    }
  }
};
</script>


<style scoped lang="scss">
article.media {
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
</style>
