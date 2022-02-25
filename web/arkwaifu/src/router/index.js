import Vue from "vue";
import VueRouter from "vue-router";
import HomeView from "../views/HomeView.vue";
import AllView from "@/views/AllView.vue";
import GroupsMainlineView from "@/views/GroupsMainlineView.vue";
import GroupsActivityView from "@/views/GroupsActivityView.vue";
import GroupsMiniView from "@/views/GroupsMiniView.vue";
import GroupsOthersView from "@/views/GroupsOthersView.vue";
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
    component: GroupsMainlineView,
  },
  {
    path: "/avgs/activity",
    name: "activity",
    component: GroupsActivityView,
  },
  {
    path: "/avgs/mini",
    name: "mini",
    component: GroupsMiniView,
  },
  {
    path: "/avgs/others",
    name: "others",
    component: GroupsOthersView,
  },
  {
    path: "/avgs/stories/:storyID",
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
