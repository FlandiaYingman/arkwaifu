import Vue from 'vue'
import vuetify from '@/plugins/vuetify'
import router from '@/router'
import lodash from 'lodash'
import App from '@/App.vue'
import store from './store'
import API_URL from './api'

Object.defineProperty(Vue.prototype, '$API_URL', { value: API_URL })
Object.defineProperty(Vue.prototype, '_', { value: lodash })

new Vue({
  vuetify,
  router,
  store,
  created () {
    store.dispatch('updateAll')
  },
  render: (h) => h(App)
}).$mount('#app')
