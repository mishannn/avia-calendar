import * as VueRouter from "vue-router";

import MainPage from "@/pages/MainPage.vue";
import ResultsPage from "@/pages/ResultsPage.vue";

const routes = [
  { path: "/", component: MainPage },
  { path: "/results", component: ResultsPage },
];

const router = VueRouter.createRouter({
  history: VueRouter.createWebHashHistory(),
  routes, // short for `routes: routes`
});

export default router;
