<template>
  <div class="col">
    <div class="card" style="min-width: 20rem">
      <div class="card-body">
        <h5 class="card-title">{{ name }}</h5>
        <h6 class="card-subtitle mb-2 text-success" :class="agentStatus === 'online' ? 'text-success' : 'text-danger'">
          {{ agentStatus }}
        </h6>
        <p class="card-text" v-if="agentStatus === 'online'">
          Program status:
          <span :class="programStatus === 'up' ? 'text-success' : 'text-danger'">{{ programStatus }}</span>
        </p>
        <p class="card-text">Last updated: {{ lastUpdated }}</p>
        <p class="card-text">Activated: {{ activated }}</p>
        <router-link class="btn btn-primary" :to="{ name: 'Device', params: { pk: pk } }">Go to device</router-link>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
export default defineComponent({
  name: "DeviceCard",
  props: {
    name: { type: String, required: true },
    agentStatus: { type: String, required: true },
    programStatus: { type: String, required: true },
    lastUpdated: { type: String, required: true },
    activated: { type: Boolean, required: true },
    pk: { type: String, required: true },
  },
});
</script>

<style scoped></style>
