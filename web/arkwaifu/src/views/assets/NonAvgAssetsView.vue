<template>
  <assets-show
    :images="nonAvgImages"
    :backgrounds="nonAvgBackgrounds"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'
import _ from 'lodash'
export default {
  components: { AssetsShow },
  computed: {
    avgImages () {
      let images = []
      images = _.flatMap(this.$store.state.avg.stories, (el) => el.images)
      images = _.uniq(images)
      return images
    },
    avgBackgrounds () {
      let backgrounds = []
      backgrounds = _.flatMap(this.$store.state.avg.stories, (el) => el.backgrounds)
      backgrounds = _.uniq(backgrounds)
      return backgrounds
    },
    nonAvgImages () {
      const { images } = this.$store.state.assets
      const { avgImages } = this
      return images.filter((el) => !avgImages.includes(el))
    },
    nonAvgBackgrounds () {
      const { backgrounds } = this.$store.state.assets
      const { avgBackgrounds } = this
      return backgrounds.filter((el) => !avgBackgrounds.includes(el))
    }
  }
}
</script>
