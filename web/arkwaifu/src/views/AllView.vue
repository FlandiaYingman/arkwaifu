<template>
  <v-container fluid>
    <v-row>
      <v-col v-for="image in images" :key="image" cols="3">
        <v-card>
          <v-img :src="image" class="transparent-background"></v-img>
        </v-card>
      </v-col>
      <v-spacer />
    </v-row>
  </v-container>
</template>

<style scoped>
.transparent-background {
  background: linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    ),
    linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    );
  background-position: 0px 0, 10px 10px;
  background-origin: padding-box, padding-box;
  background-clip: border-box, border-box;
  background-size: 20px 20px, 20px 20px;
  box-shadow: none;
}
</style>

<script>
export default {
  name: "HomeView",
  components: {},
  data() {
    return {
      images: [],
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.fetchImages();
      },
      { immediate: true }
    );
  },
  methods: {
    async fetchImages() {
      const response = await fetch(`${this.$API_URL}/api/v0/resources/images`);
      this.images = (await response.json()).map((it) => `${this.$API_URL}/api/v0/resources/images/${it}`);
    },
  },
};
</script>
