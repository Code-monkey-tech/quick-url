import { createApp } from "vue";
import App from "./App.vue";
import Toast, { PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";
import store from "./store";

createApp(App).use(store).use(Toast).mount("#app");
