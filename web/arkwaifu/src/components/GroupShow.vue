<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        {{ group.name }}
      </v-col>
      <v-col
        v-for="(story, i) in group.stories"
        :key="i"
        cols="6"
        sm="3"
        lg="2"
      >
        <story-card
          :story="story"
          :current="currentStoryId === story.id"
        />
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
import StoryCard from './StoryCard.vue'

export default {
  name: 'GroupShow',
  components: { StoryCard },
  props: {
    groupId: String(),
    currentStoryId: String()
  },
  computed: {
    group () {
      const group = this.$store.getters.groupByID(this.groupId)
      return group
    }
  }
}
</script>
