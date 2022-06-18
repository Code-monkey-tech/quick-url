/* eslint-disable import/extensions */
import { createApp } from "vue";
import PrimeVue from "primevue/config";
import InputText from "primevue/inputtext";
import OkButton from "primevue/button";
import InfoToast from "primevue/toast";
import ToastService from "primevue/toastservice";
import ProgressSpinner from "primevue/progressspinner";
import App from "./App.vue";
import store from "./store";

import "primevue/resources/themes/saga-blue/theme.css";
import "primevue/resources/primevue.min.css";
import "primeicons/primeicons.css";
import "./style/globalStyle.css";

const app = createApp(App);

app.use(PrimeVue);
app.use(ToastService);
app.component("OkButton", OkButton);
app.component("InputText", InputText);
app.component("InfoToast", InfoToast);
app.component("ProgressSpinner", ProgressSpinner);

app.use(store).mount("#app");
