<template>
  <v-container v-if="story">
    <story-card :story="story"></story-card>
    <br />
    <v-row>
      <v-col cols="12" class="text-h5">Images</v-col>
      <v-col v-for="(image, i) in story.ImageResList" :key="`${image.ID}-${i}`" cols="3">
        <v-card>
          <v-img
            :src="`http://localhost:7080/api/v0/resources/images/${image}`"
            class="transparent-background"
          ></v-img>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" class="text-h5">Backgrounds</v-col>
      <v-col v-for="(image, i) in story.BackgroundResList" :key="`${image.ID}-${i}`" cols="3">
        <v-card>
          <v-img
            :src="`http://localhost:7080/api/v0/resources/backgrounds/${image}`"
            class="transparent-background"
          ></v-img>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import StoryCard from "@/components/StoryCard.vue";
export default {
  name: "AvgGroupsView",
  components: { StoryCard },
  props: ["storyName"],
  data() {
    return {
      story: null,
    };
  },
  created() {
    this.fetchStory();
  },
  methods: {
    fetchStory() {
      fetch(`http://localhost:7080/api/v0/stories/${this.storyName}`)
        .then((resp) => resp.json())
        .then((json) => (this.story = json));
    },
  },
};
</script>
