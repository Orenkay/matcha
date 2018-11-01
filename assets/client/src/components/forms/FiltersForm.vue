<template>
  <m-form ref="form">
    <b-field 
      horizontal 
      label="Distance">
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Min'}" 
        name="distanceMin" 
        type="input" />
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Max'}" 
        name="distanceMax" 
        type="input" />
    </b-field>
    <b-field 
      horizontal 
      label="Age">
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Min'}" 
        name="ageMin" 
        type="input" />
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Max'}" 
        name="ageMax" 
        type="input" />
    </b-field>
    <b-field 
      horizontal 
      label="Popularity">
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Min'}" 
        name="popularityMin" 
        type="input" />
      <m-form-field 
        :validate="validate" 
        :props="{placeholder: 'Max'}" 
        name="popularityMax" 
        type="input" />
    </b-field>
    <b-field 
      horizontal 
      label="Sort by">
      <b-field>
        <m-form-field 
          type="select" 
          name="sort" 
          value="Pertinence">
          <option>Pertinence</option>
          <option>Distance</option>
          <option>Age</option>
          <option>Popularity</option>
        </m-form-field>
        <b-radio-button 
          v-model="sortBy" 
          native-value="asc" 
          type="is-link">
          <b-icon icon="sort-ascending"/>
          <span>Asc</span>
        </b-radio-button>

        <b-radio-button 
          v-model="sortBy" 
          native-value="desc" 
          type="is-link">
          <b-icon icon="sort-descending"/>
          <span>Desc</span>
        </b-radio-button>
      </b-field>
    </b-field>
    <b-field horizontal>
      <p class="control">
        <button 
          class="button is-link is-outlined" 
          @click="submit">
          <slot name="button-label">
            Apply settings
          </slot>
        </button>
      </p>
    </b-field>
  </m-form>
</template>

<script>
export default {
  data() {
    return {
      sortBy: 'asc'
    }
  },
  methods: {
    validate(v) {
      if (v !== undefined && v !== '' && !/^[0-9]+$/.test(v)) {
        return ['must be only digit']
      }
      return []
    },
    submit() {
      this.$refs.form.submit(d => {
        this.$emit('apply', {
          filters: {
            distance: {
              min: parseInt(d.distanceMin) || 0,
              max: parseInt(d.distanceMax) || 0
            },
            age: {
              min: parseInt(d.ageMin) || 0,
              max: parseInt(d.ageMax) || 0
            },
            popularity: {
              min: (parseInt(d.popularityMin) || 0) / 100,
              max: (parseInt(d.popularityMax) || 0) / 100
            }
          },
          sort: {
            By: d.sort.toLowerCase(),
            Desc: this.sortBy === 'desc'
          }
        })
      })
    }
  }
}
</script>

<style>
</style>
