import { State } from "./state";

export default {
  setTokens(state: State, { accessToken, refreshToken }: { accessToken: string; refreshToken: string }) {
    state.accessToken = accessToken;
    state.refreshToken = refreshToken;
  },
  deleteTokens(state: State) {
    state.accessToken = "";
    state.refreshToken = "";
  },
};