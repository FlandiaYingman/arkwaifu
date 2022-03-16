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
    <fab-button
      v-model="distinct"
      icon-on="mdi-fingerprint"
      icon-off="mdi-fingerprint-off"
      :caption-on="$t('distinctOn')"
      :caption-off="$t('distinctOff')"
    />
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow.vue'
import FabButton from '@/components/FabButton'
import AssetsShow from '@/components/AssetsShow'

export default {
  components: { AssetsShow, FabButton, GroupShow },
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

<i18n>{
  "en": {
    "distinctOn": "Distinct: ON",
    "distinctOff": "Distinct: Off"
  },
  "zh": {
    "distinctOn": "显示重复：关",
    "distinctOff": "显示重复：开"
  }
}</i18n>
