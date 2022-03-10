<template>
  <v-container class="px-8">
    <v-row v-for="group in groups" :key="group.id">
      <v-col cols="12">{{ group.name }}</v-col>
      <v-col v-for="(story, i) in group.stories" :key="`${story.id}-${i}`" cols="6" sm="3" lg="2">
        <story-card :story="story" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import StoryCard from "./StoryCard.vue";
export default {
  name: "GroupsPage",
  components: { StoryCard },
  props: ["type"],
  data() {
    return {
      groups: [],
    };
  },
  created() {
    fetch(`${this.$API_URL}/api/v0/groups`)
      .then((resp) => resp.json())
      .then((groups) => groups.filter((it) => it.actType == this.type))
      .then((groups) => (this.groups = groups));
  },
};
</script>
