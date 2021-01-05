import { InjectionKey } from "vue";
import { createStore, useStore as baseUseStore, Store } from "vuex";
import auth from "./modules/auth";

export interface State {
  auth: typeof auth.state;
}

export const store = createStore<State>({
  modules: {
    auth,
  },
});

export const key: InjectionKey<Store<State>> = Symbol();

export function useStore() {
  return baseUseStore(key);
}
