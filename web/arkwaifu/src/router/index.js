import Vue from "vue";
import VueRouter from "vue-router";
import HomeView from "../views/HomeView.vue";
import AllView from "@/views/AllView.vue";
import AvgGroupsView from "@/views/AvgGroupsView.vue";
import AvgStoriesView from "@/views/AvgStoriesView.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "home",
    component: HomeView,
  },
  {
    path: "/avgs/mainline",
    name: "mainline",
    component: AvgGroupsView,
  },
  {
    path: "/avgs/activity",
    name: "activity",
  },
  {
    path: "/avgs/operator-record",
    name: "operator-record",
  },
  {
    path: "/avgs/stories/:storyName",
    name: "story",
    component: AvgStoriesView,
    props: true,
  },
  {
    path: "/all",
    name: "all",
    component: AllView,
  },
  {
    path: "/about",
    name: "about",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ "../views/AboutView.vue"),
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
