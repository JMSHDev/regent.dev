import axios, { AxiosError, AxiosRequestConfig } from "axios";

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";
const USERNAME = "username";

declare module "axios" {
  export interface AxiosRequestConfig {
    skipIntercept?: boolean;
  }
}

const anonRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    "Content-Type": "application/json",
  },
});

const authRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  skipIntercept: false,
  headers: {
    "Content-Type": "application/json",
  },
});

const loginUser = async (username: string, password: string) => {
  const response = await anonRequest.post("/api/token/both/", { username, password });
  localStorage.setItem(ACCESS_TOKEN, response.data.access);
  localStorage.setItem(REFRESH_TOKEN, response.data.refresh);
  localStorage.setItem(USERNAME, username);
};

const refreshToken = async () => {
  const refreshBody = { refresh: localStorage.getItem(REFRESH_TOKEN) };
  const response = await anonRequest.post("/api/token/access/", refreshBody);
  localStorage.setItem(ACCESS_TOKEN, response.data.access);
};

const logoutUser = () => {
  localStorage.removeItem(ACCESS_TOKEN);
  localStorage.removeItem(REFRESH_TOKEN);
  localStorage.removeItem(USERNAME);
};

const errorInterceptor = async (error: AxiosError) => {
  const originalConfig = error.config;
  const status = error.response?.status;
  if (status === 401 && !originalConfig.skipIntercept) {
    try {
      await refreshToken();
      originalConfig.skipIntercept = true;
      return authRequest(originalConfig);
    } catch (refreshTokenError) {
      await logoutUser();
      throw refreshTokenError;
    }
  }
  throw error;
};

const authHeaderInterceptor = (requestConfig: AxiosRequestConfig) => {
  requestConfig.headers.Authorization = `Bearer ${localStorage.getItem(ACCESS_TOKEN)}`;
  return requestConfig;
};

authRequest.interceptors.request.use(authHeaderInterceptor);

authRequest.interceptors.response.use(
  (response) => response,
  (error) => errorInterceptor(error)
);

const userInfo = () => {
  return {
    username: localStorage.getItem(USERNAME),
  };
};

const isUserLoggedIn = () => {
  return userInfo().username != null;
};

export { loginUser, logoutUser, userInfo, isUserLoggedIn, authRequest, anonRequest };
