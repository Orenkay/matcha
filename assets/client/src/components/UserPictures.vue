<template>
  <div>
    <div v-if="onlyPP">
      <figure>
        <div 
          :style="`background-image: url('${pp.path}');'`" 
          class="image pp is-rounded"/>
      </figure>
    </div>
    <div v-else>
      <div class="has-text-centered">
        <figure 
          v-for="(pic, index) in pictures" 
          :key="index" 
          :class="pic.isPP && 'pp'" 
          class="image">
          <div class="edit-overlay">
            <div class="edit-overlay-content">
              <div 
                class="button close-button" 
                @click="updatePP(pic.id)">
                <b-icon 
                  icon="star" 
                  size="is-small" />
              </div>
              <div 
                class="button close-button" 
                @click="remove(pic.id)">
                <b-icon 
                  icon="close" 
                  size="is-small" />
              </div>
            </div>
          </div>
          <div 
            :class="pic.isPP && 'pp'" 
            :style="`background-image: url('${pic.path}');`" 
            class="image"/>
        </figure>
      </div>
      <b-field 
        v-if="editable" 
        class="file is-centered" 
        position="is-centered">
        <b-upload 
          v-model="file" 
          @input="upload">
          <a class="button is-primary">
            <b-icon icon="upload"/>
            <span>Click to upload</span>
          </a>
        </b-upload>
      </b-field>
    </div>
  </div>
</template>

<script>
export default {
  props: ['onlyPP', 'pictures', 'editable'],
  data() {
    return {
      file: null
    }
  },
  computed: {
    pp() {
      return this.pictures.find(p => p.isPP)
    }
  },
  methods: {
    upload(file) {
      const formData = new FormData()
      formData.append('picture', file[0])
      this.$http
        .post('/pictures/me', formData, {
          'Content-Type': 'multipart/form-data'
        })
        .then(res => {
          this.$store.commit('addPicture', res.data.data)
        })
        .catch(err => {
          if (err.response && err.response.status === 400) {
            this.$toast.error(err.response.data.error)
          }
        })
    },
    remove(id) {
      this.$http.delete(`/pictures/me/${id}`).then(res => {
        this.$store.commit('removePicture', id)
      })
    },
    updatePP(id) {
      this.$http.patch(`/pictures/me/${id}/pp`).then(res => {
        this.$store.commit('setPP', id)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
figure {
  display: inline-block;
}
.edit-overlay {
  position: absolute;

  display: flex;
  justify-content: flex-end;

  width: 100%;
  height: 100%;

  padding: 5px;
  background: rgba(0, 0, 0, 0.6);
  opacity: 0;
  transition: opacity 0.2s ease;
  z-index: 1;

  &:hover {
    opacity: 1;
  }
}

.image {
  opacity: 0.6;
  height: 128px;
  width: 128px;
  background-size: 100%;
  background-position: 50%;
}

.pp {
  opacity: 1;
}

.is-rounded {
  border-radius: 100%;
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
