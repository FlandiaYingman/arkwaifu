<template>
  <v-container>
    <group-show
      v-if="story"
      :group-id="story.groupID"
      :current-story-id="story.id"
    />
    <assets-show
      v-if="story"
      :images="story.images"
      :backgrounds="story.backgrounds"
      :distinct="distinct"
    />
    <DistinctButton v-model="distinct" />
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow.vue'
import DistinctButton from '@/components/DistinctButton'
import AssetsShow from '@/components/AssetsShow'

export default {
  components: { AssetsShow, DistinctButton, GroupShow },
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
