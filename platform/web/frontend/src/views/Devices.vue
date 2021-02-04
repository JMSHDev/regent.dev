<template>
  <NavBar />
  <div class="container">
    <div class="row row-cols-1 row-cols-md-2 row-cols-xl-3 g-4 mt-5">
      <DeviceCard
        :name="device.name"
        :status="device.status"
        :lastUpdated="device.lastUpdated"
        :activated="device.activated"
        v-for="device in state.deviceList"
        :key="device.name"
      />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive } from "vue";
import NavBar from "@/components/NavBar.vue";
import DeviceCard from "@/components/DeviceCard.vue";
import { getDeviceList, Device } from "@/services/api/devices";

export default defineComponent({
  name: "Devices",
  components: { NavBar, DeviceCard },
  setup() {
    const state = reactive({
      deviceList: new Array<Device>(),
    });

    onMounted(async () => (state.deviceList = await getDeviceList()));

    return { state };
  },
});
</script>
