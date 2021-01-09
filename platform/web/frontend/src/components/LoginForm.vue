<template>
  <div class="my-3 mx-2">
    <div>
      <p v-if="state.invalidCredentials" class="errornote">
        Please enter the correct username and password for a staff account. Note that both fields may be case-sensitive.
      </p>
    </div>
    <div>
      <div class="form">
        <form :submit="handleLogin">
          <div class="mb-3">
            <label for="username" class="form-label">Username</label>
            <input type="text" class="form-control" id="username" v-model="state.username" />
          </div>
          <div class="mb-3">
            <label for="password" class="form-label">Password</label>
            <input type="password" class="form-control" id="password" v-model="state.password" />
          </div>
          <div class="text-center">
            <button type="submit" class="btn btn-primary text-center" @click="handleLogin">
              <span v-if="state.submitting" class="spinner-border spinner-border-sm"></span>
              Log in
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import { loginUser } from "@/services/api/auth";

export default defineComponent({
  name: "LoginForm",
  setup() {
    const state = reactive({
      submitting: false,
      invalidCredentials: false,
      nextPath: "/",
      username: "",
      password: "",
    });

    const handleLogin = (event: Event) => {
      state.submitting = true;
      event.preventDefault();

      loginUser(state.username, state.password, state.nextPath).catch(() => {
        state.invalidCredentials = true;
        state.password = "";
        state.submitting = false;
      });
    };

    return { state, handleLogin };
  },
});
</script>

<style scoped>
.errornote {
  font-weight: bold;
  display: block;
  padding: 10px 12px;
  color: #ba2121;
  border: 1px solid #ba2121;
  border-radius: 4px;
  background-color: #fff;
  background-position: 5px 12px;
}
</style>
