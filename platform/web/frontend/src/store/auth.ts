import { Commit } from "vuex";
import { RouteLocationRaw } from "vue-router";
import router from "@/router";

// types
export interface AuthState {
  accessToken: string;
  refreshToken: string;
}

// helpers
function getTokensFromLocal() {
  const accessToken = localStorage.getItem("accessToken") || "";
  const refreshToken = localStorage.getItem("refreshToken") || "";
  return { accessToken, refreshToken };
}

function putTokensToLocal(accessToken: string, refreshToken: string) {
  localStorage.setItem("accessToken", accessToken);
  localStorage.setItem("refreshToken", refreshToken);
}

function delTokensFromLocal() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
}

// state
const state = getTokensFromLocal();

// mutations
const mutations = {
  setTokens(state: AuthState, { accessToken, refreshToken }: { accessToken: string; refreshToken: string }) {
    state.accessToken = accessToken;
    state.refreshToken = refreshToken;
  },
  deleteTokens(state: AuthState) {
    state.accessToken = "";
    state.refreshToken = "";
  },
};

// actions
const actions = {
  async login(
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
  async logout({ commit }: { commit: Commit }, { routerRedirect }: { routerRedirect: RouteLocationRaw }) {
    delTokensFromLocal();
    commit("deleteTokens");
    await router.push(routerRedirect);
  },
};

export const auth = {
  state: state,
  mutations: mutations,
  actions: actions,
  string: true
}
