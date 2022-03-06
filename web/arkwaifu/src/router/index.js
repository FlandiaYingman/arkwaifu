import Vue from "vue";
import VueRouter from "vue-router";
import HomeView from "../views/HomeView.vue";
import AllView from "@/views/AllView.vue";
import GroupsMainThemeView from "@/views/GroupsMainThemesView.vue";
import GroupsMajorEventsView from "@/views/GroupsMajorEventsView.vue";
import GroupsVegnettesView from "@/views/GroupsVignettesView.vue";
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
    path: "/avgs/main_themes",
    component: GroupsMainThemeView,
    meta: {
      title: "Main Themes",
    },
  },
  {
    path: "/avgs/major_events",
    component: GroupsMajorEventsView,
    meta: {
      title: "Major Events",
    },
  },
  {
    path: "/avgs/vignettes",
    component: GroupsVegnettesView,
    meta: {
      title: "Vignettes",
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
