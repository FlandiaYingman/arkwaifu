<template>
  <v-app>
    <NavigationDrawer v-model="drawer" />

    <v-app-bar app>
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-toolbar-title>{{ title }}</v-toolbar-title>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
import NavigationDrawer from '@/NavigationDrawer'

export default {
  name: 'App',

  components: { NavigationDrawer },

  data: () => ({
    title: '',
    drawer: null
  }),
  watch: {
    $route: {
      handler (to) {
        this.title = to.meta.title
      },
      immediate: true
    }
  }
}
</script>
