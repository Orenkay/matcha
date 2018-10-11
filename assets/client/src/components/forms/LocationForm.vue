<template>
      <gmap-autocomplete class="input" @place_changed="setPlace"></gmap-autocomplete>
</template>

<script>
export default {
  data() {
    return {
      place: null
    }
  },
  methods: {
    setPlace(place) {
      this.place = place
    },
    submit(cb) {
      cb(this.place)
    },
  },
  computed: {
    loc() {
      return this.$store.getters.userData.loc;
    }
  }
};
</script>

<style>
</style>
