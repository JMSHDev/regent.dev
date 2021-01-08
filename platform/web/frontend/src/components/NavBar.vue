<template>
  <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container-fluid">
      <button class="navbar-toggler" type="button">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse">
        <router-link class="navbar-brand" to="/">Home</router-link>
      </div>
    </div>
  </nav>
</template>

<script>
import { defineComponent } from "vue";

export default defineComponent({
  name: "NavBar",
});
</script>

<style scoped></style>
