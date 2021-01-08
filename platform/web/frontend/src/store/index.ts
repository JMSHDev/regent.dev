import { InjectionKey } from "vue";
import { createStore, useStore as baseUseStore, Store } from "vuex";

import { RootState, rootState } from "@/store/root";

export interface State extends RootState {}

export const store = createStore({
  state: rootState,
});

export const key: InjectionKey<Store<State>> = Symbol();

export function useStore() {
  return baseUseStore(key);
}
