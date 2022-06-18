<template>
  <div class="inner">
    <span class="inner__title">Enter your URL to make it short</span>
    <div class="inner__input-group input-group">
      <InputText v-model="inputValue" class="p-inputtext-lg input-group__input" />
      <OkButton
        icon="pi pi-check"
        class="p-button-rounded p-button-text p-button-lg input-group__button"
        @click="updateToast(toastInfo)"
      />
    </div>
    <!--    <InfoToast />-->
  </div>
</template>

<script>
import { defineComponent } from "vue";
// import { useToast } from "primevue/usetoast";

export default defineComponent({
  props: {
    updateToast: {
      type: Function,
      required: true,
    },
  },
  // setup() {
  //   const toast = useToast();
  //   const updateToast = () => {
  //     toast.add({ severity: "info", summary: "Hello", life: 1000 });
  //   };
  //   return {
  //     updateToast,
  //   };
  // },
  data() {
    return {
      inputValue: "",
      toastInfo: "info",
    };
  },
  // methods: {
  //   addToast() {
  //     this.$toast.add({ severity: "info", summary: "Let`s doing it", life: 1000 });
  //   },
  // },
});
</script>

<style scoped>
.inner {
  width: 100%;
  display: grid;
  grid-template-rows: repeat(2, min-content);
  justify-content: center;
}
.inner__title {
  font-size: 2rem;
}
.inner__input-group {
  margin-top: 2rem;
}
.input-group {
  width: 100%;
  display: flex;
  align-items: center;
}

.input-group__input {
  width: 100%;
  max-width: 30rem;
}
.input-group__button {
  margin-left: 1rem;
}
</style>
