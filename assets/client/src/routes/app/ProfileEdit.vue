<template>
  <div>
    <b-tabs>
      <b-tab-item label="Profile">
        <personal-form ref="form" />
        <br />
        <button class="button is-black" @click="submit">Edit</button>
      </b-tab-item>

      <b-tab-item label="Pictures">
        <user-pictures :editable="true" :pictures="userData.pictures" />
      </b-tab-item>

      <b-tab-item label="Interests">
        <user-interests :editable="true" :tags="userData.interests" />
      </b-tab-item>

      <b-tab-item label="Localisation">
        <user-location :editable="true" :location="userData.loc" />
      </b-tab-item>

    </b-tabs>
  </div>
</template>

<script>
import PersonalForm from "../../components/forms/PersonalForm";
import UserPictures from "../../components/UserPictures";
import UserLocation from "../../components/UserLocation";
import UserInterests from "../../components/UserInterests";

export default {
  components: {
    PersonalForm,
    UserPictures,
    UserLocation,
    UserInterests
  },
  computed: {
    userData() {
      return this.$store.getters.userData;
    }
  },
  methods: {
    submit() {
      this.$refs.form.submit(() => this.$router.push("/app/profile/me"));
    }
  }
};
</script>

<style>
.b-tabs .autocomplete .dropdown-menu {
  position: initial;
}
</style>
