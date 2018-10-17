<template>
  <m-form ref="form">
    <m-form-field name="firstName" label="Firstname" :value="profile.firstName" type="input" :required="true" :props="{placeholder: 'Your Firstname'}" />
    <m-form-field name="lastName" label="Lastname" type="input" :value="profile.lastName" :required="true" :props="{placeholder: 'Your Lastname'}" />
    <m-form-field name="birthdate" label="Birthdate" type="datepicker" :value="birthdate" :required="true" :props="{'min-date': minDate, 'max-date': maxDate, placeholder: 'Your birthdate', }" />
    <m-form-field name="gender" type="select" label="Gender" :value="profile.gender" :required="true" :props="{placeholder: 'Your Gender'}">
      <option value="male">Male</option>
      <option value="female">Female</option>
    </m-form-field>
    <m-form-field name="attraction" type="select" label="Sexual Orientation" :value="profile.attraction" :required="true" :props="{placeholder: 'Your Sexual Orientation'}">
      <option value="hetero">Heterosexual</option>
      <option value="bi">Bisexual</option>
      <option value="homo">Homosexual</option>
    </m-form-field>
    <m-form-field name="bio" label="Bio" type="input" :value="profile.bio" :required="true" :props="{type: 'textarea', placeholder: 'Describe you in few lines', maxlength: 200}" />
  </m-form>
</template>

<script>
export default {
  methods: {
    submit(cb) {
      this.$refs.form.submit(data => {
        data.birthdate = data.birthdate.getTime() / 1000;
        const req =
          this.$store.getters.profile.lastName === undefined
            ? this.$http.post("/profiles/me", data)
            : this.$http.patch("/profiles/me/edit", data);
        req.then(res => {
          this.$store.commit("setUserData", ["profile", res.data.data]);
          if (typeof cb === "function") {
            cb();
          }
        });
      });
    }
  },
  data() {
    const maxDate = new Date();
    maxDate.setFullYear(maxDate.getFullYear() - 18);
    return {
      minDate: new Date((1 << 31) * 1000),
      maxDate: maxDate
    };
  },
  computed: {
    birthdate() {
      return this.profile.birthdate && new Date(this.profile.birthdate * 1000);
    },
    profile() {
      return this.$store.getters.userData.profile;
    }
  }
};
</script>

<style>
</style>
