<template>
  <NavBar />
  <div class="container">
    {{ state.deviceList }}
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, reactive } from "vue";
import NavBar from "@/components/NavBar.vue";

import { getDeviceList, Device } from "@/services/api/devices";

export default defineComponent({
  name: "Devices",
  components: { NavBar },
  setup() {
    const state = reactive({
      deviceList: new Array<Device>(),
    });

    onMounted(async () => state.deviceList = await getDeviceList());

    return { state };
  },
});
</script>
