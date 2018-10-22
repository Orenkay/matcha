<template>
  <article 
    v-if="userData !== undefined" 
    class="media">
    <figure class="media-left">
      <div 
        :style="`background-image: url('${pp}')`" 
        class="image is-64x64"/>
    </figure>
    <div class="media-content">
      <div class="content">
        <p>
          <strong>{{ userData.profile.lastName }} {{ userData.profile.firstName }}</strong>
        </p>
      </div>
    </div>
  </article>
</template>

<script>
export default {
  props: ['target'],
  data() {
    return {
      userData: undefined
    }
  },
  computed: {
    pp() {
      return this.userData.pictures.find(p => p.isPP).path
    }
  },
  created() {
    this.$store.dispatch('getProfile', this.target).then(data => {
      this.userData = data
    })
  }
}
</script>

<style scoped>
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
  border-radius: 50px;
}
</style>
