<template>
  <m-form ref="form">
    <b-field horizontal label="Distance">
      <m-form-field name="distanceMin" type="input" :props="{placeholder: 'Min'}" />
      <m-form-field name="distanceMax" type="input" :props="{placeholder: 'Max'}" />
    </b-field>
    <b-field horizontal label="Age">
      <m-form-field name="ageMin" type="input" :props="{placeholder: 'Min'}" />
      <m-form-field name="ageMax" type="input" :props="{placeholder: 'Max'}" />
    </b-field>
    <b-field horizontal label="Popularity">
      <m-form-field name="popularityMin" type="input" :props="{placeholder: 'Min'}" />
      <m-form-field name="popularityMax" type="input" :props="{placeholder: 'Max'}" />
    </b-field>
    <b-field horizontal label="Sort by">
      <m-form-field type="select" name="sort" value="0">
        <option value="0">Default</option>
        <option value="1">Distance</option>
        <option value="2">Age</option>
        <option value="3">Popularity</option>
      </m-form-field>
    </b-field>
    <b-field horizontal>
      <p class="control">
        <button class="button is-link is-outlined" @click="submit">
          Apply settings
        </button>
      </p>
    </b-field>
  </m-form>
</template>

<script>
export default {
  methods: {
    submit() {
      this.$refs.form.submit(d => {
        this.$emit("apply", {
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
          sort: parseInt(d.sort)
        });
      });
    }
  }
};
</script>

<style>
</style>
