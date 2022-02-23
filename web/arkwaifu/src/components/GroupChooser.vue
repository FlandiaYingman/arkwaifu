<template>
  <v-container>
    <v-row v-for="group in groups" :key="group.ID">
      <v-col cols="12">{{ group.Name }}</v-col>
      <v-col v-for="(story, i) in group.StoryList" :key="`${story.Name}-${i}`" cols="auto">
        <story-card :story="story"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import StoryCard from './StoryCard.vue';
export default {
  components: { StoryCard },
  name: "GroupChooser",
  // props: ["groups"],
  data: () => ({
    groups: [],
  }),
  created() {
    fetch("http://localhost:7080/api/v0/groups")
      .then((resp) => resp.json())
      .then((json) => (this.groups = json));
  },
};
</script>
