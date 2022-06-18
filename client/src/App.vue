<template>
  <div class="main">
    <MainHeader class="main__header" />
    <InputGroup class="main__input-group" :update-toast="updateToast" />
    <ViewResult class="main__view-result" :update-toast="updateToast" />
    <InfoToast />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useToast } from "primevue/usetoast";
import InputGroup from "./components/InputGroup.vue";
import MainHeader from "./components/MainHeader.vue";
import ViewResult from "./components/ViewResult.vue";
import HandelSpiner from "./components/HandlerSpiner.vue";

export default defineComponent({
  name: "App",
  components: {
    InputGroup,
    MainHeader,
    ViewResult,
  },
  setup() {
    const toast = useToast();
    const updateToast = (info: string, text?: string) => {
      toast.add({
        severity: info,
        summary: text !== undefined ? text : "Let's do it",
        life: 1000,
      });
    };
    return {
      updateToast,
    };
  },
  data() {
    return {};
  },
  methods: {},
});
</script>

<style>
.main {
  width: 100%;
  height: 100vh;
  display: grid;
  grid-template-rows: 4rem min-content 1fr;
  grid-template-columns: 1fr;
}
.main__input-group {
  margin-top: 5rem;
}
.main__view-result {
  margin-top: 5rem;
}
</style>
