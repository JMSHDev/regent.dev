<template>
  <div class="card" style="width: 18rem">
    <div class="card-body">
      <h5 class="card-title">{{ name }}</h5>
      <h6 class="card-subtitle mb-2 text-muted">{{ status }}</h6>
      <p class="card-text">Last modified: {{ lastUpdated }}</p>
      <p class="card-text">Activated: {{ activated }}</p>
      <a href="#" class="btn btn-primary">Go to device</a>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
  name: "DeviceCard",
  props: {
    name: { type: String, required: true },
    status: { type: String, required: true },
    lastUpdated: { type: String, required: true },
    activated: { type: Boolean, required: true },
  },
});
</script>

<style scoped></style>
