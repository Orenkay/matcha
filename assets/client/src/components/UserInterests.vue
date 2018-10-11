<template>
  <div>
    <b-field grouped group-multiline>
      <div class="control" v-for="(tag, index) in tags" :key="index">
        <b-taglist attached>
          <b-tag type="is-dark">{{ tag.value }}</b-tag>
          <b-tag type="is-danger button" @click.native="remove(tag.value)">
            <b-icon icon="close" size="is-small" />
          </b-tag>
        </b-taglist>
      </div>
    </b-field>
    <b-field v-if="editable">
      <b-autocomplete v-model="name" :data="data" :loading="isFetching" @keyup.native="getAsyncData" expanded></b-autocomplete>
      <p class="control">
        <button class="button" @click="add">+</button>
      </p>
    </b-field>
  </div>
</template>

<script>
export default {
  props: ["editable", "tags"],
  data() {
    return {
      data: [],
      name: "",
      isFetching: false
    };
  },
  methods: {
    add() {
      this.$http
        .post("/interests/me", { value: this.name }, { errorHandle: false })
        .then(({ data }) => {
          this.name = "";
          console.log(data.data.value);
          this.$store.commit("addInterest", data.data);
        })
        .catch(err => {
          if (err.response && err.response.status === 400) {
            this.$toast.open(err.response.data.error);
          }
        });
    },
    remove(val) {
      this.$http.delete("/interests/me/" + val).then(res => {
        this.$store.commit("removeInterest", val);
      });
    },
    getAsyncData() {
      if (!this.name.length) {
        this.data = [];
        return;
      }
      this.isFetching = true;
      this.$http
        .get("/interests/" + this.name, { errorHandle: false })
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

<style>
</style>
