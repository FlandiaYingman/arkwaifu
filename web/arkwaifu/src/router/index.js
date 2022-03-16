import Vue from 'vue'
import VueRouter from 'vue-router'
import AboutView from '@/views/AboutView.vue'
import HomeView from '@/views/HomeView.vue'
import GroupsView from '@/views/assets/GroupsView'
import StoryView from '@/views/assets/StoryView'
import NonAvgView from '@/views/assets/NonAvgAssetsView'
import AllAssetsView from '@/views/assets/AllAssetsView'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: HomeView,
    meta: {
      title: 'Home'
    }
  },
  {
    path: '/avgs/main_themes',
    component: GroupsView,
    meta: {
      title: 'Main Themes'
    },
    props: {
      type: 'MAIN_STORY'
    }
  },
  {
    path: '/avgs/major_events',
    component: GroupsView,
    meta: {
      title: 'Major Events'
    },
    props: {
      type: 'ACTIVITY_STORY'
    }
  },
  {
    path: '/avgs/vignettes',
    component: GroupsView,
    meta: {
      title: 'Vignettes'
    },
    props: {
      type: 'MINI_STORY'
    }
  },
  {
    path: '/avgs/others',
    component: GroupsView,
    meta: {
      title: 'Others'
    },
    props: {
      type: 'NONE'
    }
  },
  {
    path: '/avgs/stories/:storyID',
    component: StoryView,
    props: true,
    meta: {
      title: 'Story'
    }
  },
  {
    path: '/non_avgs',
    component: NonAvgView,
    meta: {
      title: 'Non-AVG'
    }
  },
  {
    path: '/all',
    component: AllAssetsView,
    meta: {
      title: 'All'
    }
  },
  {
    path: '/about',
    component: AboutView,
    meta: {
      title: 'About'
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
