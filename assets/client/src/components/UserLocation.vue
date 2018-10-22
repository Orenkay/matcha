<template>
  <div>
    <b-icon 
      pack="fas" 
      icon="map-marker-alt" 
      size="is-small" />
    <span 
      :class="!location.address && 'nowhere'" 
      class="is-inline">{{ location.address || 'you have no location' }}</span>
    <div v-if="editable">
      <hr >
      <gmap-autocomplete 
        ref="ginput" 
        v-bind="gprops" 
        class="input" 
        @place_changed="setPlace"/>
    </div>
  </div>
</template>

<script>
export default {
  props: ['editable', 'location'],
  data() {
    return {
      gprops: {
        componentRestrictions: {
          country: 'FR'
        },
        types: ['address']
      }
    }
  },
  methods: {
    setPlace({ place_id }) {
      const req =
        this.$store.getters.loc.address === undefined
          ? this.$http.post('/loc/me', { placeId: place_id })
          : this.$http.patch('/loc/me/edit', { placeId: place_id })
      req.then(res => {
        this.$store.commit('setUserData', ['loc', res.data.data])
      })
    }
  }
}
</script>

<style scoped>
.nowhere {
  color: grey;
  font-style: italic;
}
</style>
