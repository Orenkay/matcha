<template>
  <div>
    <div class="has-text-centered">
      <figure v-for="(pic, index) in pictures" :key="index" class="image is-128x128">
        <div class="edit-overlay image is-128x128">
          <div class="edit-overlay-content">
            <div class="button close-button" @click="remove(pic.id)">
              <b-icon icon="close" size="is-small" />
            </div>
          </div>
        </div>
        <img :src="pic.url" />
      </figure>
    </div>
    <b-field class="file is-centered" position="is-centered">
      <b-upload v-model="file" @input="upload">
        <a class="button is-primary">
          <b-icon icon="upload"></b-icon>
          <span>Click to upload</span>
        </a>
      </b-upload>
    </b-field>
  </div>
</template>

<script>
export default {
  props: ["pictures"],
  data() {
    return {
      file: null
    };
  },
  methods: {
    upload(file) {
      const formData = new FormData();
      formData.append("picture", file[0]);
      this.$http
        .post("/pictures/me", formData, {
          errorHandle: false,
          "Content-Type": "multipart/form-data"
        })
        .then(res => {
          console.log(res);
        });
    },
    remove(id) {
      console.log(id);
    }
  }
};
</script>

<style scoped>
figure {
  display: inline-block;
}
.edit-overlay {
  position: absolute;
  display: flex;
  padding: 5px;
  justify-content: flex-end;
}
.close-button {
  background: transparent;
  border: none;
}
.close-button:hover {
  background: rgba(0, 0, 0, 0.1);
}
</style>
