import Vue from "vue";
import Vuetify from "vuetify/lib/framework";
// Translation provided by Vuetify (javascript)
import zhHans from "vuetify/lib/locale/zh-Hans";

Vue.use(Vuetify);

export default new Vuetify({
  lang: {
    locales: { zhHans },
    current: "zhHans",
  },
});
