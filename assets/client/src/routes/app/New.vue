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
        <location-form ref="locForm" />
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
import LocationForm from "../../components/forms/LocationForm";
import UserPictures from "../../components/UserPictures";
import UserInterests from "../../components/UserInterests";

export default {
  components: {
    PersonalForm,
    LocationForm,
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
      if (this.userData.pictures.length === 0) return 2;
      if (this.userData.loc.address === undefined) return 3;
    }
  },
  methods: {
    step1next(next) {
      this.$refs.personalForm.submit(data => {
        const req =
          this.$store.getters.profile.lastName === undefined
            ? this.$http.post("/profiles", data)
            : this.$http.patch("/profiles/edit", data);
        req.then(res => {
          this.$store.commit("setUserData", ["profile", res.data.data]);
          next();
        });
      });
    },
    step2next(next) {
      if (this.userData.interests.length === 0) {
        return this.$toast.open({message:"You must add atleast one interest", queue: false, type: "is-danger"});
      }
      next();
    },
    step3next(next) {
      if (this.userData.pictures.findIndex(p => p.isPP) < 0) {
        return this.$toast.open({message:"You must have atleast one PP", queue: false, type: "is-danger"});
      }
      next();
    },
    step4next(next) {
      this.$refs.locForm.submit(data => {
        const req=
          this.$store.getters.loc.address === undefined
              ? this.$http.post("/loc/me", {placeId: data.place_id})
              : this.$http.patch("/loc/me/edit", {placeId: data.place_id});
        req.then(res => {
          this.$store.commit("setUserData", ["loc", res.data.data]);
          next();
        });
      })
    }
  }
};
</script>
