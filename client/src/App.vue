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
      />

      <input
        class="input content__short"
        type="text"
        disabled
        value="www.short.com"
      />
      <button
        :class="{
          'btn-disabled': !shortUrl.length,
          'btn-active': shortUrl.length,
        }"
        class="btn content__copy"
      >
        Copy URL
      </button>
      <div class="content__qr qr">
        <img
          :class="{
            'qr__img-disabled': !shortUrl.length,
          }"
          class="qr__img"
          src="./assets/qr_mock.png"
          alt="qr"
        />
        <button
          :class="{
            'btn-disabled': !shortUrl.length,
            'btn-active': shortUrl.length,
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
      >
        Shorten URL
      </button>
    </div>
  </div>
</template>

<script>
import { defineComponent } from "vue";
import ShortyLogoSVG from "./assets/ShortyLogo.vue";

export default defineComponent({
  name: "App",
  components: {
    ShortyLogoSVG,
  },
  data() {
    return {
      longUrl: "",
      shortUrl: "",
      colorScheme: "light",
    };
  },
  mounted() {
    const initialScheme = this.getSettings() || this.getMediaPrefers();
    this.setSettings(initialScheme);
  },
  methods: {
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
    setSettings(scheme) {
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
