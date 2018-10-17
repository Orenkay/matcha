<template>
  <div :class="type">
    <b-field grouped group-multiline>
      <div class="control" v-for="(tag, index) in tags" :key="index">
        <b-taglist attached>
          <b-tag type="is-dark">{{ tag.value }}</b-tag>
          <b-tag v-if="editable" type="is-danger button" @click.native="remove(tag.value)">
            <b-icon icon="close" size="is-small" />
          </b-tag>
        </b-taglist>
      </div>
    </b-field>
    <b-field v-if="editable">
      <b-autocomplete v-model="value" :data="data" :loading="isFetching" @keyup.native="getAsyncData" expanded></b-autocomplete>
      <p class="control">
        <button class="button" @click="add">+</button>
      </p>
    </b-field>
  </div>
</template>

<script>
export default {
  props: ["editable", "tags", "type"],
  data() {
    return {
      data: [],
      value: "",
      isFetching: false
    };
  },
  computed: {
    getValue() {
      return this.value.toLowerCase();
    }
  },
  methods: {
    add() {
      this.$http
        .post("/interests/me", { value: this.getValue })
        .then(({ data }) => {
          this.value = "";
          this.$store.commit("addInterest", data.data);
        })
        .catch(err => {
          if (err.response && err.response.status === 400) {
            this.$toast.error(err.response.data.error);
          }
        });
    },
    remove(val) {
      this.$http.delete("/interests/me/" + val).then(res => {
        this.$store.commit("removeInterest", val);
      });
    },
    getAsyncData() {
      if (!this.value.length) {
        this.data = [];
        return;
      }
      this.isFetching = true;
      this.$http
        .get("/interests/" + this.getValue)
        .then(({ data }) => {
          this.data = [];
          if (data.data === null) {
            return;
          }
          data.data.forEach(t => {
            this.data.push(t.value);
          });
        })
        .catch(() => {
          this.data = [];
        })
        .finally(() => (this.isFetching = false));
    }
  }
};
</script>

<style scoped>
.is-centered > .field.is-grouped {
  justify-content: center;
}
</style>
