import { Commit } from "vuex";
import { RouteLocationRaw } from "vue-router";

import router from "@/router";
import { putTokensToLocal, delTokensFromLocal } from "./util";

export default {
  login: async function login(
    { commit }: { commit: Commit },
    {
      accessToken,
      refreshToken,
      routerRedirect,
    }: { accessToken: string; refreshToken: string; routerRedirect: RouteLocationRaw }
  ) {
    putTokensToLocal(accessToken, refreshToken);
    commit("setTokens", { accessToken, refreshToken });
    await router.push(routerRedirect);
  },
  logout: async function logout(
    { commit }: { commit: Commit },
    { routerRedirect }: { routerRedirect: RouteLocationRaw }
  ) {
    delTokensFromLocal();
    commit("deleteTokens");
    await router.push(routerRedirect);
  },
};
