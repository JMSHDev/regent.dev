<template>
  <nav class="navbar navbar-expand-md navbar-dark bg-dark">
    <div class="container">
      <router-link class="navbar-brand" to="/">regent.dev</router-link>
      <button class="navbar-toggler" type="button" @click="switchCollapse">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="navbar-collapse" :class="[state.collapse]">
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <router-link class="nav-link" to="/">Devices</router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link me-3" to="/">Docs</router-link>
          </li>
        </ul>

        <button class="btn btn-danger" @click="handleLogout">Logout</button>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { logoutUser } from "@/services/api/auth";
import router from "@/router";

export default defineComponent({
  name: "NavBar",
  setup() {
    const state = reactive({
      collapse: "collapse",
    });

    const switchCollapse = () => {
      state.collapse === "collapse" ? (state.collapse = "collapse show") : (state.collapse = "collapse");
    };

    const handleLogout = () => {
      logoutUser();
      router.push({ name: "Login" });
    };

    return { state, switchCollapse, handleLogout };
  },
});
</script>

<style scoped></style>
