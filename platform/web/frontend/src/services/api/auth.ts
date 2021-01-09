import { AxiosError } from "axios";
import axios from "axios";

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";
const USERNAME = "username";

const tokenRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    "Content-Type": "application/json",
    accept: "application/json",
  },
});

const loginUser = async (username: string, password: string) => {
  const response = await tokenRequest.post("/api/token/both/", { username, password });
  localStorage.setItem(ACCESS_TOKEN, response.data.refresh);
  localStorage.setItem(REFRESH_TOKEN, response.data.refresh);
  localStorage.setItem(USERNAME, username);
};

const refreshToken = async () => {
  const refreshBody = { refresh: localStorage.getItem(REFRESH_TOKEN) };
  const response = await tokenRequest.post("/api/token/access/", refreshBody);
  localStorage.setItem(ACCESS_TOKEN, response.data.access);
};

const authRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    Authorization: `Bearer ${localStorage.getItem(ACCESS_TOKEN)}`,
    "Content-Type": "application/json",
  },
});

const logoutUser = () => {
  localStorage.removeItem(ACCESS_TOKEN);
  localStorage.removeItem(REFRESH_TOKEN);
  localStorage.removeItem(USERNAME);
  authRequest.defaults.headers.Authorization = "";
};

const errorInterceptor = async (error: AxiosError) => {
  const originalRequest = error.config;
  const status = error.response?.status;
  const accessToken = localStorage.getItem(ACCESS_TOKEN);
  if (status === 401 && accessToken) {
    try {
      await refreshToken();
      const headerAuthorization = `Bearer ${accessToken}`;
      authRequest.defaults.headers.Authorization = headerAuthorization;
      originalRequest.headers.Authorization = headerAuthorization;
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
