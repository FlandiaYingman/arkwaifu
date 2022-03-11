import Vue from "vue";
import { API_URL } from "./api";
import App from "@/App.vue";
import vuetify from "@/plugins/vuetify";
import router from "@/router";
import store from "./store";
import lodash from "lodash";

Object.defineProperty(Vue.prototype, "$API_URL", { value: API_URL });
Object.defineProperty(Vue.prototype, "_", { value: lodash });

new Vue({
  vuetify,
  router,
  store,
  render: (h) => h(App),
  created() {
    store.dispatch("updateAll");
  },
}).$mount("#app");
