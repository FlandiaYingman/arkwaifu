<template>
  <v-card :to="`/avgs/stories/${story.id}`" :disabled="current">
    <v-card-title>{{ story.name }}</v-card-title>
    <v-card-subtitle>
      {{ story.code ? `${story.code} ` : "" }}{{ story.tag }}
      <br />
      {{ story.id }}
    </v-card-subtitle>
  </v-card>
</template>

<script>
export default {
  name: "StoryCard",
  props: ["story", "current"],
  data() {
    return {};
  },
  methods: {
    frontPic() {
      if (this.story.images.length > 0) {
        return `${this.$API_URL}/api/v0/resources/images/${this.story.images[0]}`;
      } else {
        return `${this.$API_URL}/api/v0/resources/backgrounds/${this.story.backgrounds[0]}`;
      }
    },
  },
};
</script>
