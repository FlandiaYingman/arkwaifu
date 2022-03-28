import Vue from 'vue'
import Vuetify, { colors } from 'vuetify/lib'
import { ThemeOptions } from 'vuetify/types/services/theme'

const mainTheme: ThemeOptions = {
  dark: false,
  default: 'light',
  themes: {
    light: {
      primary: colors.deepPurple.base,
      secondary: colors.pink.lighten2,
      accent: colors.indigo.base
    }
  },
  options: { customProperties: true }
}

Vue.use(Vuetify, {
  options: {
    customProperties: true
  }
})

export default new Vuetify({
  lang: {
    current: 'zhHans'
  },
  icons: {
    iconfont: 'mdiSvg'
  },
  theme: mainTheme
})

declare module '@/plugins/vuetify' {
  interface Vue {
    $vuetify: any
  }
}
