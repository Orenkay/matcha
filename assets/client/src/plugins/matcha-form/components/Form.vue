<template>
  <form onsubmit="return false">
    <slot />
  </form>
</template>

<script>
export default {
  data() {
    return {
      fields: {},
      _isMatchaForm: true
    }
  },
  methods: {
    getFormData() {
      const data = {}
      for (let k in this.fields) {
        data[k] = this.fields[k].v
      }
      return data
    },
    fieldError(k, errors) {
      this.fields[k].errors = errors
    },
    submit(cb) {
      const data = this.getFormData()
      let submit = true
      for (let k in this.fields) {
        const valid = this.fields[k]._validate(this.fields[k].v, data)
        submit && (submit = valid)
      }
      if (submit) {
        cb(data)
      }
    }
  }
}
</script>
