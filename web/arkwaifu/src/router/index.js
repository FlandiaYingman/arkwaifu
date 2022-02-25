import Vue from "vue";
import VueRouter from "vue-router";
import HomeView from "../views/HomeView.vue";
import AllView from "@/views/AllView.vue";
import GroupsMainlineView from "@/views/GroupsMainlineView.vue";
import GroupsActivityView from "@/views/GroupsActivityView.vue";
import GroupsMiniView from "@/views/GroupsMiniView.vue";
import GroupsOthersView from "@/views/GroupsOthersView.vue";
import AvgStoriesView from "@/views/AvgStoriesView.vue";
import AboutView from "@/views/AboutView.vue"

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: HomeView,
    meta: {
      title: "Home",
    },
  },
  {
    path: "/avgs/mainline",
    component: GroupsMainlineView,
    meta: {
      title: "Mainline",
    },
  },
  {
    path: "/avgs/activity",
    component: GroupsActivityView,
    meta: {
      title: "Activity",
    },
  },
  {
    path: "/avgs/mini",
    component: GroupsMiniView,
    meta: {
      title: "Mini Activity",
    },
  },
  {
    path: "/avgs/others",
    component: GroupsOthersView,
    meta: {
      title: "Others",
    },
  },
  {
    path: "/avgs/stories/:storyID",
    component: AvgStoriesView,
    props: true,
    meta: {
      title: "Story",
    },
  },
  {
    path: "/all",
    component: AllView,
    meta: {
      title: "All",
    },
  },
  {
    path: "/about",
    component: AboutView,
    meta: {
      title: "About",
    },
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

export default router;
