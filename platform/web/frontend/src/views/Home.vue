<template>
  <NavBar />
  <div class="container">
    Hola!
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import NavBar from "@/components/NavBar.vue";

export default defineComponent({
  name: "HomePage",
  components: { NavBar },
});
</script>
