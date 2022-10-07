<template>
  <qrcode-vue
    :value="$store.getters.getCorrectShortUrl"
    :size="size"
    level="H"
  />
</template>
<script lang="ts">
import { defineComponent } from "vue";
import QrcodeVue from "qrcode.vue";

export default defineComponent({
  name: "QrcodeGenerate",
  components: {
    QrcodeVue,
  },
  data() {
    return {
      size: 140,
    };
  },
});
</script>
<style scoped></style>
