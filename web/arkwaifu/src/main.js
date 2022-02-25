import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './router'

Vue.prototype.$API_URL = process.env.VUE_APP_API_URL

new Vue({
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app')
