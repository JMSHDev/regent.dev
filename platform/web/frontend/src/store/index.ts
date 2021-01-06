import { InjectionKey } from "vue";
import { createStore, useStore as baseUseStore, Store } from "vuex";

import {state, State} from "@/store/state";
import mutations from "@/store/mutations";
import actions from "@/store/actions";


export const store = createStore<State>({
  state: state,
  mutations: mutations,
  actions: actions
});

export const key: InjectionKey<Store<State>> = Symbol();

export function useStore() {
  return baseUseStore(key);
}
