<template>
  <div class="header">
    <a
      aria-label="Link to Github"
      class="header__github-link"
      href="https://github.com/code-monkey-tech/quick-url"
      target="_blank"
    >
      <i class="pi pi-github pi-github" style="font-size: 3rem"></i>
    </a>
    <span class="header__title">Quick URL</span>
  </div>
</template>

<script></script>

<style scoped>
.header {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  background: lavender;
  padding: 0 50px 0 20px;
}
.header__github-link {
  text-decoration: none;
  cursor: pointer;
  color: black;
}
.header__title {
  font-size: 22px;
  margin: 0 auto 0 auto;
  color: black;
}
</style>
