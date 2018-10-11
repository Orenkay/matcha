<template>
  <div>
    <b-tabs>
      <b-tab-item label="Profile">
        <m-form :fields="form.fields" ref="form" />
        <br />
        <button class="button is-black" @click="submit">Edit</button>
      </b-tab-item>

      <b-tab-item label="Pictures">
        TODO
      </b-tab-item>

      <b-tab-item label="Localisation">
        TODO
      </b-tab-item>

    </b-tabs>
  </div>
</template>

<script>
import Form from "../../components/Form";
export default {
  components: {
    "m-form": Form
  },
  data() {
    const profile = this.$store.getters.profile;

    return {
      form: {
        fields: {
          firstName: {
            label: "Firstname",
            value: profile.firstName,
            required: true
          },
          lastName: {
            label: "Lastname",
            value: profile.lastName,
            required: true
          },
          gender: {
            label: "Gender",
            value: profile.gender,
            type: "select",
            options: [
              { value: "male", label: "Male" },
              { value: "female", label: "Female" }
            ],
            required: true
          },
          attraction: {
            label: "Sexual Orientation",
            value: profile.attraction,
            type: "select",
            options: [
              { value: "bi", label: "Bisexual" },
              { value: "hetero", label: "Heterosexual" },
              { value: "homo", label: "Homosexual" }
            ],
            required: true
          },
          bio: {
            label: "Bio",
            value: profile.bio,
            props: () => ({ type: "textarea", maxlength: 200 }),
            required: true
          }
        }
      }
    };
  },
  methods: {
    submit() {
      this.$refs.form.submit(data => {
        this.$http.put("/profiles/edit", data).then(res => {
          this.$store.commit("setUserData", ["profile", res.data.data]);
          this.$router.push("/app/profile");
        });
      });
    }
  }
};
</script>
