import Vue from 'vue'
import VueRouter from 'vue-router'
import AboutView from '@/views/AboutView.vue'
import HomeView from '@/views/HomeView.vue'
import StoryView from '@/views/assets/StoryView'
import NonAvgView from '@/views/assets/NonAvgAssetsView'
import AllAssetsView from '@/views/assets/AllAssetsView'
import AvgMainView from '@/views/assets/AvgMainView'
import AvgMajorView from '@/views/assets/AvgMajorView'
import AvgMiniView from '@/views/assets/AvgMiniView'
import AvgOthersView from '@/views/assets/AvgOthersView'

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
    component: AvgMainView,
    meta: {
      title: 'Main Themes'
    }
  },
  {
    path: '/avgs/major_events',
    component: AvgMajorView,
    meta: {
      title: 'Major Events'
    }
  },
  {
    path: '/avgs/vignettes',
    component: AvgMiniView,
    meta: {
      title: 'Vignettes'
    }
  },
  {
    path: '/avgs/others',
    component: AvgOthersView,
    meta: {
      title: 'Others'
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
