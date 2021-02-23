<template>
  <NavBar />
  <div class="container">
    <h1>{{ state.device.name }}</h1>
    <h4 class="text-muted">{{ state.device.status }}</h4>
    <p>Last updated: {{ state.device.lastUpdated }}</p>
    <p>Activated: {{ state.device.activated }}</p>
    <router-link class="btn btn-primary" :to="{name: 'Devices'}">Device List</router-link>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive } from "vue";
import NavBar from "@/components/NavBar.vue";
import { nullDevice, getDevice } from "@/services/api/devices";

export default defineComponent({
  name: "Device",
  components: { NavBar },
  props: { pk: { type: String, required: true } },
  setup(props) {
    const state = reactive({
      device: nullDevice,
    });

    onMounted(async () => (state.device = await getDevice(props.pk)));

    return { state };
  },
});
</script>

<style scoped></style>
