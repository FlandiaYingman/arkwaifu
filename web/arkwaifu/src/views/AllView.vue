<template>
  <v-container fluid>
    <v-row>
      <v-col v-for="image in images" :key="image" cols="3">
        <resource-card :resName="image" resCategory="images"></resource-card>
      </v-col>
      <v-spacer />
    </v-row>
    <v-row>
      <v-col v-for="background in backgrounds" :key="background" cols="3">
        <resource-card :resName="background" resCategory="backgrounds"></resource-card>
      </v-col>
      <v-spacer />
    </v-row>
  </v-container>
</template>

<script>
import ResourceCard from "@/components/ResourceCard.vue";
export default {
  name: "HomeView",
  components: { ResourceCard },
  data() {
    return {
      images: [],
      backgrounds: [],
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.fetchImages();
        this.fetchBackgrounds();
      },
      { immediate: true }
    );
  },
  methods: {
    async fetchImages() {
      const response = await fetch(`${this.$API_URL}/api/v0/resources/images`);
      this.images = await response.json();
    },
    async fetchBackgrounds() {
      const response = await fetch(`${this.$API_URL}/api/v0/resources/backgrounds`);
      this.backgrounds = await response.json();
    },
  },
};
</script>
