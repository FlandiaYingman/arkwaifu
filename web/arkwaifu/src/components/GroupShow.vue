<template>
  <v-container v-if="group">
    <v-row>
      <v-col cols="12">
        {{ group.name }}
      </v-col>
      <v-col
        v-for="(story, i) in !limited
          ? group.stories
          : group.stories.slice(0, 3)"
        :key="`${story.id}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <story-card
          :story="story"
          :current="currentStoryID == story.id"
        />
      </v-col>
      <v-col
        v-if="limited"
        cols="auto"
      >
        <v-card>
          <v-card-title>...</v-card-title>
          <v-card-subtitle>
            ...
            <br>
          </v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
import StoryCard from './StoryCard.vue';

export default {
  name: 'GroupShow',
  components: { StoryCard },
  props: {
    groupID: String(),
    currentStoryID: String(),
    limited: Boolean,
  },
  data: () => ({}),
  computed: {
    group() {
      const group = this.$store.getters.groupByID(this.groupID);
      return group;
    },
  },
};
</script>
