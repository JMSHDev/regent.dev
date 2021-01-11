<template>
  <router-view :key="$route.fullPath" />
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "App",
});
</script>

<style lang="scss">
html {
  width: 100%;
  height: 100%;
  background-size: contain;
  background-color: var(--bs-light);
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: var(--bs-light);
}
</style>
