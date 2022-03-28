<template>
  <v-app>
    <nav-drawer
      v-model="drawer"
      class="elevation-4"
    />

    <v-app-bar
      app
      dark
      class="primary"
    >
      <v-app-bar-nav-icon
        @click="drawer = !drawer"
      />
      <v-toolbar-title>{{ title }}</v-toolbar-title>
    </v-app-bar>

    <v-main class="ma-8">
      <!--this container is for setting the max-width-->
      <v-container>
        <slot />
      </v-container>
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
