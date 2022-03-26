import Vue from 'vue'
import VueRouter from 'vue-router'
import AboutView from '@/views/AboutView.vue'
import HomeView from '@/views/HomeView.vue'
import StoryView from '@/views/assets/StoryView.vue'
import NonAvgView from '@/views/assets/NonAvgAssetsView.vue'
import AllAssetsView from '@/views/assets/AllAssetsView.vue'
import AvgMainView from '@/views/assets/AvgMainView.vue'
import AvgMajorView from '@/views/assets/AvgMajorView.vue'
import AvgMiniView from '@/views/assets/AvgMiniView.vue'
import AvgOthersView from '@/views/assets/AvgOthersView.vue'
import { RouteConfigSingleView } from 'vue-router/types/router'
import AssetView from '@/views/assets/AssetView.vue'

Vue.use(VueRouter)

const routes: RouteConfigSingleView[] = [
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
  },
  {
    path: '/assets/:kind/:id',
    component: AssetView,
    props: true
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
