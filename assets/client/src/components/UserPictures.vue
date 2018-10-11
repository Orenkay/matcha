<template>
  <div>
    <div class="has-text-centered">
      <figure v-for="(pic, index) in pictures" :key="index" class="image" :class="pic.isPP && 'pp'">
        <div class="edit-overlay">
          <div class="edit-overlay-content">
            <div class="button close-button" @click="pp(pic.id)">
              <b-icon icon="star" size="is-small" />
            </div>
            <div class="button close-button" @click="remove(pic.id)">
              <b-icon icon="close" size="is-small" />
            </div>
          </div>
        </div>
        <img :src="`http://192.168.99.100:3000/pictures/${pic.path}`" />
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
          "Content-Type": "multipart/form-data"
        })
        .then(res => {
          this.$store.commit('addPicture', res.data.data)
        });
    },
    remove(id) {
      this.$http
        .delete(`/pictures/me/${id}`)
        .then(res => {
          this.$store.commit('removePicture', id)
        });
    },
    pp(id) {
      console.log(id)
      this.$http
        .patch(`/pictures/me/${id}/pp`)
        .then(res => {
          this.$store.commit('setPP', id)
        });
    }
  }
};
</script>

<style lang="scss" scoped>
figure {
  display: inline-block;
  width: 128px;
}
.edit-overlay {
  position: absolute;

  display: flex;
  justify-content: flex-end;

  width: 100%;
  height: 100%;

  padding: 5px;
  background: rgba(0, 0, 0, .6);
  opacity: 0;
  transition: opacity .2s ease;

  &:hover {
    opacity: 1;
  }
}

.image {
  opacity: .6;
}

.pp {
  opacity: 1;
}

.close-button {
  background: transparent;
  border: none;
  color: white;
}
.close-button:hover {
  background: rgba(0, 0, 0, 0.1);
}
</style>
