<template>
  <b-field :label="label" :message="errors" :type="errors.length > 0 ? 'is-danger' : ''">
    <component :is="'b-'+type" v-bind="props" v-model="v">
      <slot />
    </component>
  </b-field>
</template>

<script>
export default {
  props: ["name", "label", "value", "type", "props", "required", "validate"],
  methods: {
    _validate(v, formData) {
      this.errors = [];
      if (this.required === true) {
        if (v === undefined) {
          this.errors = ["field is required"];
        }
      }
      if (typeof this.validate === "function") {
        this.errors.push(...this.validate(v, formData));
      }
      return this.errors.length === 0;
    },
    _getFormParent() {
      let c = this.$parent;
      while (c !== undefined) {
        if (c.$data._isMatchaForm) {
          return c;
        }
        c = c.$parent;
      }
    }
  },
  data() {
    return {
      v: this.value,
      errors: []
    };
  },
  created() {
    this.formParent = this._getFormParent();
    if (!this.formParent) {
      this.$destroy();
      throw new Error("You should wrap mFormField on a mForm");
    }
    this.formParent.fields[this.name] = this;
  },
  beforeDestroy() {
    delete this.formParent.fields[this.name];
  }
};
</script>

<style>
</style>
