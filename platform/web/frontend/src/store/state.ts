import { getTokensFromLocal } from "./util";

export interface State {
  accessToken: string;
  refreshToken: string;
}

export const state = getTokensFromLocal();
