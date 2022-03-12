<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      width="288"
      permanant
      app
    >
      <v-img
        src="/logo_phantom.png"
        width="128"
        height="128"
        class="mx-4"
      />

      <v-list nav>
        <v-subheader>NAVIGATION</v-subheader>

        <v-list-item
          link
          to="/"
        >
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Home</v-list-item-title>
            <router-link to="/home" />
          </v-list-item-content>
        </v-list-item>

        <v-list-group
          link
          prepend-icon="mdi-television"
        >
          <template #activator>
            <v-list-item-content>
              <v-list-item-title>AVG</v-list-item-title>
              <v-list-item-subtitle>Sort by Stories</v-list-item-subtitle>
            </v-list-item-content>
          </template>
          <v-list-item to="/avgs/main_themes">
            <v-list-item-icon><v-icon>mdi-cube</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Main Themes</v-list-item-title>
              <v-list-item-subtitle>Main Stories</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/major_events">
            <v-list-item-icon><v-icon>mdi-music-accidental-sharp</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Major Events</v-list-item-title>
              <v-list-item-subtitle>Intermezzi & Side Stories</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/vignettes">
            <v-list-item-icon><v-icon>mdi-music-accidental-flat</v-icon></v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Vignettes</v-list-item-title>
              <v-list-item-subtitle>Story Collections</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
          <v-list-item to="/avgs/others">
            <v-list-item-icon>
              <v-icon>mdi-music-rest-quarter</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Others</v-list-item-title>
              <v-list-item-subtitle>Operator Records</v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>

        <v-list-item
          link
          to="/all"
        >
          <v-list-item-icon>
            <v-icon>mdi-all-inclusive</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>All</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          link
          to="/about"
        >
          <v-list-item-icon>
            <v-icon>mdi-information</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>About</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

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
export default {
  name: 'App',

  components: {},

  data: () => ({
    title: '',
    drawer: null,
  }),
  watch: {
    $route: {
      handler(to) {
        this.title = to.meta.title;
      },
      immediate: true,
    },
  },
};
</script>
