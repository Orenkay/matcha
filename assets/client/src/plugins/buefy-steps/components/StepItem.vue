<template>
  <div class="step-content" :class="isActive && 'is-active'">
    <slot />
  </div>
</template>

<script>
export default {
  props: ["title", "desc", "next"],
  data() {
    return {
      isActive: false
    };
  },
  methods: {
    onNext(next) {
      if (typeof this.next === "function") {
        return this.next(() => next());
      }
      next();
    }
  },
  created() {
    if (!this.$parent.$data._isSteps) {
      this.$destroy();
      throw new Error("You should wrap bStepItem on a bSteps");
    }
    this.$parent.stepItems.push(this);
  },
  beforeDestroy() {
    const index = this.$parent.stepItems.indexOf(this);
    if (index >= 0) {
      this.$parent.stepItems.splice(index, 1);
    }
  }
};
</script>

<style>
</style>
