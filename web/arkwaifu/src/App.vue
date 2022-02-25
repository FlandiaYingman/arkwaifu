<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" permanant app width="288">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="text-h4">ArkWaifu</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />

      <v-list nav>
        <v-subheader>NAVIGATION</v-subheader>

        <v-list-item link to="/">
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Home</v-list-item-title>
            <router-link to="/home"></router-link>
          </v-list-item-content>
        </v-list-item>

        <v-list-group link prepend-icon="mdi-television">
          <template #activator>
            <v-list-item-content>
              <v-list-item-title>AVGs</v-list-item-title>
              <v-list-item-subtitle>Group Resources</v-list-item-subtitle>
            </v-list-item-content>
          </template>
          <v-list-item to="/avgs/mainline">
            <v-list-item-icon><v-icon></v-icon></v-list-item-icon>
            <v-list-item-title>Mainline</v-list-item-title>
          </v-list-item>
          <v-list-item to="/avgs/activity">
            <v-list-item-icon><v-icon></v-icon></v-list-item-icon>
            <v-list-item-title>Activity</v-list-item-title>
          </v-list-item>
          <v-list-item to="/avgs/mini">
            <v-list-item-icon><v-icon></v-icon></v-list-item-icon>
            <v-list-item-title>Mini Activity</v-list-item-title>
          </v-list-item>
          <v-list-item to="/avgs/others">
            <v-list-item-icon><v-icon></v-icon></v-list-item-icon>
            <v-list-item-title>Others</v-list-item-title>
          </v-list-item>
        </v-list-group>

        <v-list-item link to="/all">
          <v-list-item-icon>
            <v-icon>mdi-all-inclusive</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>All</v-list-item-title>
            <v-list-item-subtitle>All Resources</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>

        <v-list-item link to="/about">
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
      <router-view> </router-view>
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: "App",

  components: {},

  data: () => ({
    title: "",
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
