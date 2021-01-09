<template>
  <NavBar />
  <div class="container">
    <h1>API Ping test</h1>
    <button v-on:click="testPing" class="btn btn-success mb-2">Ping</button>
    <p>{{ state.pingValue }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import NavBar from "@/components/NavBar.vue";
import { ping } from "@/services/api/ping";
import { logoutUser } from "@/services/api/auth";

export default defineComponent({
  name: "AuthenticatedPing",
  components: { NavBar },
  setup() {
    const state = reactive({
      pingValue: "",
    });

    const testPing = () => {
      state.pingValue = "";
      ping()
        .then((data) => {
          state.pingValue = data.id;
        })
        .catch((error) => {
          alert(error);
        });
    };

    return { state, testPing };
  },
});
</script>
