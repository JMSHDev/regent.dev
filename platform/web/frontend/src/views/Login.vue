<template>
  <div class="container">
    <div class="login-page">
      <div class="text-center bg-primary text-white">
        <h4 class="py-3">Command Center</h4>
      </div>
      <LoginForm />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import LoginForm from "../components/LoginForm.vue";

export default defineComponent({
  name: "LoginPage",
  components: {
    LoginForm,
  },
});
</script>

<style scoped>
.login-page {
  background: #fff;
  border: 1px solid #eaeaea;
  border-radius: 4px;
  overflow: hidden;
  max-width: 28rem;
  margin: 8rem auto;
}
</style>
