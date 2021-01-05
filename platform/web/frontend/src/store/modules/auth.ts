// helper functions
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

interface TokenPair {
  accessToken: string;
  refreshToken: string;
}

// initial state
const state: TokenPair = getTokensFromLocal();

// mutations
const mutations = {
  setTokens(state: TokenPair, newTokens: TokenPair) {
    state.accessToken = newTokens.accessToken;
    state.refreshToken = newTokens.refreshToken;
  },
  deleteTokens(state: TokenPair) {
    state.accessToken = "";
    state.refreshToken = "";
  },
};

//actions
const actions = {};

// getters
const getters = {};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
};
