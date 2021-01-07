import axios from "axios";

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";

const tokenRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    "Content-Type": "application/json",
    accept: "application/json",
  },
});

export const loginUser = async (username: string, password: string): Promise<void> => {
  try {
    const response = await tokenRequest.post("/api/token/both/", { username, password });
    localStorage.setItem(ACCESS_TOKEN, response.data.refresh);
    localStorage.setItem(REFRESH_TOKEN, response.data.refresh);
  } catch (error) {
    console.log(error);
  }
};

const authRequest = axios.create({
  baseURL: process.env.VUE_APP_ROOT_API,
  timeout: process.env.VUE_APP_API_TIMEOUT,
  headers: {
    Authorization: `Bearer ${window.localStorage.getItem(ACCESS_TOKEN)}`,
    "Content-Type": "application/json",
  },
});

export const logoutUser = async (): Promise<void> => {
  localStorage.removeItem(ACCESS_TOKEN);
  localStorage.removeItem(REFRESH_TOKEN);
  authRequest.defaults.headers.Authorization = "";
};
