<template>
  <v-app>
    <nav-drawer v-model="drawer" />

    <v-app-bar app>
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-toolbar-title>{{ title }}</v-toolbar-title>
    </v-app-bar>

    <v-main class="ma-8">
      <slot />
    </v-main>
  </v-app>
</template>

<script>
import NavDrawer from '@/components/NavDrawer'

export default {
  name: 'MainLayout',
  components: { NavDrawer },
  props: {
    title: String()
  },
  data () {
    return {
      drawer: null
    }
  }
}
</script>
