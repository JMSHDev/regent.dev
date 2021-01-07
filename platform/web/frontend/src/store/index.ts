import { InjectionKey } from "vue";
import { createStore, useStore as baseUseStore, Store } from "vuex";

import { RootState, rootState } from "@/store/root";
import { AuthState, auth } from "@/store/auth";

export interface State extends RootState {
  auth: AuthState;
}

export const store = createStore({
  state: rootState,
  modules: { auth },
});

export const key: InjectionKey<Store<State>> = Symbol();

export function useStore() {
  return baseUseStore(key);
}
