import Vue from 'vue'
import vuetify from '@/plugins/vuetify'
import router from '@/router'
import lodash from 'lodash'
import App from '@/App.vue'
import store from './store'
import i18n from '@/i18n'
import API_URL from './api'
import AsyncComputed from 'vue-async-computed'

Object.defineProperty(Vue.prototype, '$API_URL', { value: API_URL })
Object.defineProperty(Vue.prototype, '_', { value: lodash })

Vue.use(AsyncComputed)

new Vue({
  vuetify,
  router,
  store,
  i18n,
  created () {
    store.dispatch('updateAll')
  },
  render: (h) => h(App)
}).$mount('#app')
