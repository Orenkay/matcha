<template>
  <div class="steps">
    <div 
      v-for="(stepItem, index) in stepItems" 
      :key="index" 
      :class="[activeStep === index && 'is-active', index < activeStep && 'is-completed']" 
      class="step-item">
      <div class="step-marker">{{ index + 1 }}</div>
      <div class="step-details">
        <p class="step-title">{{ stepItem.title }}</p>
        <p>{{ stepItem.desc }}</p>
      </div>
    </div>
    <div class="steps-content">
      <slot />
    </div>
    <div class="steps-actions">
      <div class="steps-action">
        <a 
          :disabled="activeStep === 0" 
          href="#" 
          data-nav="previous" 
          class="button is-light" 
          @click="prev">Previous</a>
      </div>
      <div class="steps-action">
        <a 
          :disabled="activeStep === stepItems.length -1" 
          href="#" 
          data-nav="next" 
          class="button is-light" 
          @click="next">Next</a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: ['index'],
  data() {
    return {
      activeStep: this.index || 0,
      stepItems: [],
      _isSteps: true
    }
  },
  computed: {
    itemClass() {
      return []
    }
  },
  mounted() {
    if (this.stepItems.length) {
      this.stepItems[this.activeStep].isActive = true
    }
  },
  methods: {
    next() {
      const index = this.activeStep + 1
      if (index < this.stepItems.length) {
        this.stepItems[this.activeStep].onNext(() => this.changeStep(index))
      }
    },
    prev() {
      const index = this.activeStep - 1
      if (index >= 0) {
        this.changeStep(index)
      }
    },
    changeStep(index) {
      if (index === this.activeStep) return

      this.stepItems[this.activeStep].isActive = false
      this.stepItems[index].isActive = true
      this.activeStep = index
    }
  }
}
</script>

<style>
</style>
