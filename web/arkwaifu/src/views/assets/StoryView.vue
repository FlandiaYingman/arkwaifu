<template>
  <v-container>
    <group-show
      v-if="story"
      :group-id="story.groupID"
      :current-story-id="story.id"
    />
    <assets-show
      v-if="story"
      :assets="story.assets"
      :distinct="distinct"
    />
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow.vue'
import AssetsShow from '@/components/AssetsShow'

export default {
  components: { AssetsShow, GroupShow },
  props: {
    storyID: String()
  },
  data () {
    return {
      distinct: true
    }
  },
  computed: {
    story () {
      const story = this.$store.getters.storyByID(this.storyID)
      return story
    }
  }
}
</script>
