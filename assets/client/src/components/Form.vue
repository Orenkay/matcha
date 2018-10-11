<template>
  <form>
    <b-field v-for="(f, k) in fields" :key="k" :label="f.label" :type="f.errors.length > 0 ? 'is-danger' : ''" :message="f.errors">
      <component :is="getFieldComponent(f)" v-model="f.value" :placeholder="f.placeholder" @blur="validField(f, formData())" v-bind="f.props && f.props()">

        <template v-if="f.type === 'select'">
          <option v-for="(opts, index) in f.options" :key="index" :value="opts.value">{{ opts.label }}</option>
        </template>

      </component>
    </b-field>
  </form>
</template>

<script>
export default {
  props: ["fields"],
  beforeMount() {
    for (let k in this.fields) {
      if (this.fields[k].errors === undefined) this.fields[k].errors = [];
    }
  },
  methods: {
    getFieldComponent(f) {
      switch (f.type) {
        case "select":
          return "b-select";
        default:
          return "b-input";
      }
    },
    formData() {
      const data = {};
      for (let k in this.fields) {
        data[k] = this.fields[k].value;
      }
      return data;
    },
    validField(field, formData) {
      field.errors = [];
      if (typeof field.validation === "function") {
        field.errors = field.validation(field.value, formData);
      }
      if (
        field.required &&
        (field.value === undefined || field.value.trim() === "")
      ) {
        field.errors.push("field is required");
      }
      this.$forceUpdate();
      return field.errors.length === 0;
    },
    fieldError(field, error) {
      this.fields[field] && this.fields[field].errors.push(error);
      this.$forceUpdate();
    },
    submit(cb) {
      let submit = true;
      let formData = this.formData();
      for (let k in this.fields) {
        let valid = this.validField(this.fields[k], formData);
        submit && (submit = valid);
      }
      submit && cb(formData);
    }
  }
};
</script>
