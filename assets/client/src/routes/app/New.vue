<template>
  <section class="section">
    <b-steps :index="currentStep">
      <b-step-item title="Personal Informations" :next="step1next">
        <personal-form ref="personalForm" />
      </b-step-item>
      <b-step-item title="Your Interests" :next="step2next">
        <user-interests :editable="true" :tags="userData.interests" />
      </b-step-item>
      <b-step-item title="Your Pictures" :next="step3next">
        <user-pictures :editable="true" :pictures="userData.pictures" />
      </b-step-item>
      <b-step-item title="Your Location" :next="step4next">
        <user-location :editable="true" :location="userData.loc" />
      </b-step-item>
      <b-step-item title="Done!">
        <div class="has-text-centered">
          <h1 class="title">Welcome in Matcha !</h1>
          <router-link to="/app">Click here to enjoy matcha</router-link>
        </div>
      </b-step-item>
    </b-steps>
  </section>
</template>

<script>
import PersonalForm from "../../components/forms/PersonalForm";
import UserLocation from "../../components/UserLocation";
import UserPictures from "../../components/UserPictures";
import UserInterests from "../../components/UserInterests";

export default {
  components: {
    PersonalForm,
    UserLocation,
    UserPictures,
    UserInterests
  },
  computed: {
    userData() {
      return this.$store.getters.userData;
    },
    currentStep() {
      if (this.userData.profile.lastName === undefined) return 0;
      if (this.userData.interests.length === 0) return 1;
      if (this.userData.pictures.findIndex(p => p.isPP) < 0) return 2;
      if (this.userData.loc.address === undefined) return 3;
    }
  },
  methods: {
    step1next(next) {
      this.$refs.personalForm.submit(next);
    },
    step2next(next) {
      if (this.userData.interests.length === 0) {
        return this.$toast.error("You must add atleast one interest");
      }
      next();
    },
    step3next(next) {
      if (this.userData.pictures.findIndex(p => p.isPP) < 0) {
        return this.$toast.error("You must have atleast one PP");
      }
      next();
    },
    step4next(next) {
      if (this.userData.loc.address === undefined) {
        return this.$toast.error("You must set your location");
      }
      next();
    }
  }
};
</script>
