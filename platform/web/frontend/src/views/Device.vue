<template>
  <NavBar />
  <div class="container">
    <h1>{{ state.device.name }}</h1>
    <h4 :class="state.device.status === 'online' ? 'text-success' : 'text-danger'">{{ state.device.status }}</h4>
    <p>Last updated: {{ state.device.lastUpdated }}</p>
    <p>Activated: {{ state.device.activated }}</p>
    <router-link class="btn btn-primary" :to="{ name: 'Devices' }">Back to list</router-link>
    <h4 class="mt-5">Telemetry</h4>
    <div class="telemetry-desc">
      <p>To get device's telemetry use the <a :href="telemetryApi">telemetry API</a>.</p>
      <p>
        For example to get telemetry between 2021/02/10 and 2021/02/12 use GET
        <a :href="telemetryExample + state.device.name">{{ telemetryExample + state.device.name}}</a>
      </p>
      <p>
        To browse the API you have to login as admin. To use the API from an external application you need to add the
        following authorization header: <b>Authorization: Token token_value</b> where token_value can be found at
        <a :href="adminToken">admin tokens</a>.
      </p>
    </div>
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

    const telemetryApi = process.env.VUE_APP_ROOT_API + "/api/telemetry/";
    const telemetryExample = process.env.VUE_APP_ROOT_API + "/api/telemetry/?start=2021-02-10&end=2021-02-12&device=";
    const adminToken = process.env.VUE_APP_ROOT_API + "/admin/authtoken/tokenproxy/";

    return { state, telemetryApi, adminToken, telemetryExample };
  },
});
</script>

<style scoped>
.telemetry-desc {
  background-color: white;
  padding: 1rem;
  margin-right: 0;
  margin-left: 0;
  border: 1px black solid;
  border-radius: 0.25rem;
}
</style>
