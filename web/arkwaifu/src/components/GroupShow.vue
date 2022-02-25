<template>
  <v-container v-if="group">
    <v-row>
      <v-col cols="12">{{ group.name }}</v-col>
      <v-col
        v-for="(story, i) in !limited ? group.stories : group.stories.slice(0, 3)"
        :key="`${story.id}-${i}`"
        cols="auto"
      >
        <story-card :story="story" />
      </v-col>
      <v-col cols="auto" v-if="limited">
        <v-card>
          <v-card-title>...</v-card-title>
          <v-card-subtitle>
            ...
            <br />
          </v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
import StoryCard from "./StoryCard.vue";
export default {
  name: "group-show",
  props: ["groupID", "limited"],
  components: { StoryCard },
  data: () => ({ group: null }),
  created() {
    fetch(`${this.$API_URL}/api/v0/groups/${this.groupID}`)
      .then((resp) => resp.json())
      .then((json) => (this.group = json));
  },
  watch: {
    groupID: function () {
      fetch(`${this.$API_URL}/api/v0/groups/${this.groupID}`)
        .then((resp) => resp.json())
        .then((json) => (this.group = json));
    },
  },
};
</script>
