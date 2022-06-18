<template>
  <div class="view-result">
    <span class="view-result__title"> Here you can copy a short URL or scan a QR code </span>
    <div class="view-result__response response">
      <div class="response__short-url short-url">
        <span class="short-url__link">{{ shortURL }}</span>
        <OkButton
          icon="pi pi-copy"
          class="p-button-rounded p-button-text p-button-lg short-url__btn"
          @click="handleCopyAndToast(shortURL)"
        />
      </div>
      <div v-if="test > 0" class="response__qr qr">
        <img src="../assets/default-qr.png" alt="QR-code" class="qr__img" />
      </div>
      <HandlerSpiner v-else class="view-result__spiner spinner" :size-spinner="sizeSpinner" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import useClipboard from "vue-clipboard3";
import HandlerSpiner from "./HandlerSpiner.vue";

export default defineComponent({
  components: { HandlerSpiner },
  // @click="copy(shortURL)"
  // props: ["updateToast"],
  props: {
    updateToast: {
      type: Function,
      required: true,
    },
  },
  setup() {
    const { toClipboard } = useClipboard();
    const toastText = "Copied to clipboard";
    const copy = async (text: string): Promise<any> => {
      // this.updateToast();
      try {
        await toClipboard(text);
      } catch (e) {
        console.log(e);
      }
    };
    return { copy };
  },
  data() {
    return {
      shortURL: `www.google.com`,
      toastInfo: "success",
      test: 1,
      sizeSpinner: {
        widthSpinner: 10,
        heightSpinner: 10,
      },
    };
  },
  methods: {
    handleCopyAndToast(text: string) {
      this.copy(text);
      this.updateToast(this.toastInfo, "Copied");
    },
  },
});
</script>

<style scoped>
.view-result {
  width: 100%;
  display: grid;
  grid-template-rows: min-content 1fr;
  justify-content: center;
}
.view-result__title {
  font-size: 2rem;
}
.view-result__response {
  margin-top: 2rem;
}
.response {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.response__short-url {
  display: flex;
  align-items: center;
}
.short-url__link {
  font-size: 1.8rem;
}
.short-url__btn {
  margin-left: 1rem;
}
.response__qr {
  margin-top: 2rem;
}
.qr {
  max-width: 30rem;
  max-height: 30rem;
}
.view-result__spiner {
  margin-top: 10rem;
}
</style>
