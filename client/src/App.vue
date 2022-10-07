<template>
  <div class="main">
    <div class="header">
      <ShortyLogoSVG class="header__logo" />
      <div class="header__scheme-switcher">
        <input
          id="checkbox"
          type="checkbox"
          class="custom-checkbox"
          @change="toggleScheme"
        />
        <label for="checkbox" class="custom-label">
          <span>🌙</span>
          <span>☀️</span>
          <div
            class="swich-toggle"
            :class="{ 'switch-toggle-checked': colorScheme === 'dark' }"
          ></div>
        </label>
      </div>
    </div>
    <div class="leftContent">
      <span class="leftContent__one">Shorten</span>
      <span class="leftContent__two">your</span>
      <span class="leftContent__three">path.</span>
    </div>
    <div class="content">
      <input
        v-model.trim="longUrl"
        class="input content__long"
        type="text"
        placeholder="Your link"
        @keyup.enter="longUrl.length && sendLongUrl()"
      />

      <input
        v-model="$store.getters.getCorrectShortUrl"
        class="input content__short"
        type="text"
        placeholder="https://shorty.com/sGn2"
        disabled
      />
      <button
        :class="{
          'btn-disabled': !$store.getters.getCorrectShortUrl.length,
          'btn-active': $store.getters.getCorrectShortUrl.length,
        }"
        class="btn content__copy"
        @click="
          $store.getters.getCorrectShortUrl &&
            copyURL($store.getters.getCorrectShortUrl)
        "
      >
        Copy URL
      </button>
      <div class="content__qr qr">
        <qrcode-generate
          :class="{
            'qr__img-disabled': !$store.getters.getCorrectShortUrl.length,
          }"
          class="qr__img"
        />

        <button
          :class="{
            'btn-disabled': !$store.getters.getCorrectShortUrl.length,
            'btn-active': $store.getters.getCorrectShortUrl.length,
          }"
          class="btn qr__btn"
        >
          Dowload QR
        </button>
      </div>
      <img class="content__cut" src="./assets/cut.png" alt="cut" />
      <button
        :class="{
          'btn-disabled': !longUrl.length,
          'btn-active': longUrl.length,
        }"
        class="btn content__btn"
        @click="longUrl.length && sendLongUrl()"
      >
        Shorten URL
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { requestExpandUrl } from "./api";
import ShortyLogoSVG from "./assets/ShortyLogo.vue";
import { useToast } from "vue-toastification";
import QrcodeGenerate from "./components/QrcodeGenerate.vue";
import { IUrlData } from "./components/types";

export default defineComponent({
  name: "App",
  components: {
    ShortyLogoSVG,
    QrcodeGenerate,
  },
  setup() {
    const toast = useToast();
    return { toast };
  },
  data() {
    return {
      longUrl: "",
      colorScheme: "light",
    };
  },
  computed: {
    hashParam() {
      const path = window.location.pathname.slice(1);
      return path;
    },
  },
  async created() {
    if (this.hashParam) {
      const { url } = await requestExpandUrl(this.hashParam);
      url && (window.location.href = url);
    }
  },
  async mounted() {
    const initialScheme = this.getSettings() || this.getMediaPrefers();
    this.setSettings(initialScheme);
  },
  methods: {
    async copyURL(mytext: string) {
      try {
        await navigator.clipboard.writeText(mytext);
        this.toast.success("You copied link! Use it");
      } catch ($e) {
        this.toast.error("So sorry, cannot copy!");
      }
    },
    async sendLongUrl() {
      const requestData: IUrlData = {
        url: this.longUrl,
      };
      await this.$store.dispatch("getShortUrl", requestData, {
        root: true,
      });
    },
    toggleScheme() {
      const localSettings = localStorage.getItem("user-scheme");
      if (localSettings === "light") {
        this.setSettings("dark");
      } else {
        this.setSettings("light");
      }
    },
    getSettings() {
      return localStorage.getItem("user-scheme");
    },
    setSettings(scheme: string) {
      localStorage.setItem("user-scheme", scheme);
      this.colorScheme = scheme;
      document.documentElement.className = scheme;
    },
    getMediaPrefers() {
      const setMediaPrefers = window.matchMedia(
        "(prefers-color-sheme: dark)"
      ).matches;
      if (setMediaPrefers) {
        return "dark";
      } else {
        return "light";
      }
    },
  },
});
</script>

<style scoped>
@import "./styles/style.css";
</style>
