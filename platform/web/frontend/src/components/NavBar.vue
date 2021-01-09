<template>
  <nav class="navbar navbar-expand-md navbar-dark bg-dark">
    <div class="container">
      <router-link class="navbar-brand" to="/">Command Center</router-link>
      <button class="navbar-toggler" type="button" @click="switchCollapse">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="navbar-collapse" :class="[state.collapse]">
        <ul class="navbar-nav me-auto mb-2 mb-md-0">
          <li class="nav-item">
            <router-link class="nav-link" to="/">Home</router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";

export default defineComponent({
  name: "NavBar",
  setup() {
    const state = reactive({
      collapse: "collapse",
    });

    const switchCollapse = () => {
      state.collapse === "collapse" ? (state.collapse = "collapse show") : (state.collapse = "collapse");
    };

    return { state, switchCollapse };
  },
});
</script>

<style scoped></style>
