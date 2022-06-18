<template>
  <div>
    <ProgressSpinner class="p-progress-spinner-color size" :style="cssProps" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: {
    sizeSpinner: {
      type: Object,
      required: true,
    },
  },
  computed: {
    cssProps() {
      return {
        "--width-spinner": `${this.sizeSpinner.widthSpinner}rem`,
        "--height-spinner": `${this.sizeSpinner.heightSpinner}rem`,
      };
    },
  },
});
</script>

<style scoped>
.size {
  height: var(--width-spinner);
  width: var(--height-spinner);
}
@keyframes p-progress-spinner-color {
  100%,
  0% {
    stroke: #d62d20;
  }
  40% {
    stroke: #0057e7;
  }
  66% {
    stroke: #008744;
  }
  80%,
  90% {
    stroke: #ffa700;
  }
}
</style>
