import { Commit } from "vuex";
import { loginUser, logoutUser } from "@/services/api/auth"

// types
export interface AuthState {
  user: string | null;
  isLoggedIn: boolean;
}

// state
const state: AuthState = { user: null, isLoggedIn: false };

// mutations
const mutations = {
  loginSuccess(state: AuthState, userId: string) {
    state.user = userId;
    state.isLoggedIn = true;
  },
  logout(state: AuthState) {
    state.user = null;
    state.isLoggedIn = false;
  },
};

// actions
const actions = {
  login({ commit }: { commit: Commit }, { username, password }: { username: string; password: string }) {
    return loginUser(username, password)
      .then(() => {
        commit({ type: "loginSuccess", username });
        return Promise.resolve();
      })
      .catch((error: any) => {
        commit({ type: "logout" });
        return Promise.reject(error);
      });
  },
  logout({ commit }: { commit: Commit }) {
    logoutUser();
    commit("logout");
  },
};

export const auth = {
  state: state,
  mutations: mutations,
  actions: actions,
  string: true,
};
