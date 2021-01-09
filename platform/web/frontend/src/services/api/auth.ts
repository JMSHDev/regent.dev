import axios, { AxiosError } from "axios";
import router from "@/router";

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";
const USERNAME = "username";

declare module "axios" {
  export interface AxiosRequestConfig {
    skipIntercept?: boolean;
  }
}

const tokenRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    "Content-Type": "application/json",
    accept: "application/json",
  },
});

const authRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  skipIntercept: false,
  headers: {
    Authorization: `Bearer ${localStorage.getItem(ACCESS_TOKEN)}`,
    "Content-Type": "application/json",
  },
});

const loginUser = async (username: string, password: string) => {
  const response = await tokenRequest.post("/api/token/both/", { username, password });
  localStorage.setItem(ACCESS_TOKEN, response.data.access);
  localStorage.setItem(REFRESH_TOKEN, response.data.refresh);
  localStorage.setItem(USERNAME, username);

  authRequest.defaults.headers.Authorization = `Bearer ${response.data.access}`;
};

const refreshToken = async () => {
  const refreshBody = { refresh: localStorage.getItem(REFRESH_TOKEN) };
  const response = await tokenRequest.post("/api/token/access/", refreshBody);
  localStorage.setItem(ACCESS_TOKEN, response.data.access);

  authRequest.defaults.headers.Authorization = `Bearer ${response.data.access}`;
};

const logoutUser = () => {
  localStorage.removeItem(ACCESS_TOKEN);
  localStorage.removeItem(REFRESH_TOKEN);
  localStorage.removeItem(USERNAME);
  authRequest.defaults.headers.Authorization = "";
  router.push({ name: "Login" }).then(() => {});
};

const errorInterceptor = async (error: AxiosError) => {
  const originalRequest = error.config;
  const status = error.response?.status;
  const accessToken = localStorage.getItem(ACCESS_TOKEN);
  if (status === 401 && accessToken && !originalRequest.skipIntercept) {
    try {
      await refreshToken();
      const newAccessToken = localStorage.getItem(ACCESS_TOKEN);
      originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
      originalRequest.skipIntercept = true;
      return authRequest(originalRequest);
    } catch (error) {
      logoutUser();
      throw error;
    }
  }
  throw error;
};

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

export { loginUser, logoutUser, userInfo, isUserLoggedIn, authRequest };
