const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,

  pluginOptions: {
    i18n: {
      locale: 'zh',
      fallbackLocale: 'en',
      localeDir: 'locales',
      enableInSFC: true,
      enableBridge: false
    }
  }
})
