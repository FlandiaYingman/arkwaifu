<template>
  <v-app>
    <nav-drawer v-model="drawer" />

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
import NavDrawer from '@/components/NavDrawer'

export default {
  name: 'App',

  components: { NavDrawer },

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
