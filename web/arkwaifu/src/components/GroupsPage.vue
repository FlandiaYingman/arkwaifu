<template>
  <v-container class="px-8">
    <v-row
      v-for="group in groups"
      :key="group.id"
    >
      <v-col cols="12">
        {{ group.name }}
      </v-col>
      <v-col
        v-for="(story, i) in group.stories"
        :key="`${story.id}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <story-card :story="story" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import StoryCard from './StoryCard.vue';

export default {
  name: 'GroupsPage',
  components: { StoryCard },
  props: {
    type: String(),
  },
  data() {
    return {};
  },
  computed: {
    groups() {
      const { groups } = this.$store.state.avg;
      return groups.filter((el) => el.actType === this.type);
    },
  },
};
</script>
