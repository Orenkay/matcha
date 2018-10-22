<template>
  <b-tooltip 
    :label="tooltip" 
    type="is-black" 
    position="is-bottom">
    <button 
      :class="buttonClasses" 
      :disabled="loading" 
      class="button" 
      @click="click">
      <b-icon 
        :icon="buttonIcon" 
        size="is-small" />
    </button>
  </b-tooltip>
</template>

<script>
export default {
  props: ['target'],
  data() {
    return {
      isLiked: false,
      likeMe: false,
      loading: false
    }
  },
  computed: {
    selfData() {
      return this.$store.getters.userData
    },
    tooltip() {
      if (this.isLiked) {
        return 'Unlike'
      }
      if (this.likeMe) {
        return 'Like-back'
      }
      return 'Like'
    },
    buttonClasses() {
      return [this.likeMe && 'is-danger', this.loading && 'is-loading']
    },
    buttonIcon() {
      return this.isLiked ? 'heart-broken' : 'heart'
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      this.$http
        .get(`/likes/me/${this.target}`)
        .then(res => {
          const { data } = res.data
          this.isLiked = data.liked
          this.likeMe = data.likeMe
        })
        .finally(() => {
          this.loading = false
        })
    },
    click() {
      this.loading = true
      const req = !this.isLiked
        ? this.$http.post(`/likes/me/${this.target}`)
        : this.$http.delete(`/likes/me/${this.target}`)
      req
        .then(() => {
          this.isLiked = !this.isLiked
        })
        .finally(() => {
          this.loading = false
        })
    }
  }
}
</script>

<style scoped>
.middle-container {
  padding-top: 0;
}
</style>
