import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home from "../views/Home.vue";
import Login from "../views/Login.vue";
import AuthenticatedPing from "../views/AuthennticatedPing.vue";

import { isUserLoggedIn } from "@/services/api/auth";

const PUBLIC_PATHS = ["/login"];

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/ping",
    name: "AuthenticatedPing",
    component: AuthenticatedPing,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const unAuthenticatedAndPrivatePage = (path: string) => {
  return !PUBLIC_PATHS.includes(path) && !isUserLoggedIn();
};

router.beforeEach((to, from, next) => {
  if (unAuthenticatedAndPrivatePage(to.path)) {
    next(`/login?next=${to.path}`);
  } else {
    next();
  }
});

export default router;
