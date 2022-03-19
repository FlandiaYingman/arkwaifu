<template>
  <assets-show
    :assets="nonAvgAssets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'
import _ from 'lodash'
export default {
  components: { AssetsShow },
  computed: {
    avgAssets () {
      let assets = []
      assets = _.flatMap(this.$store.state.avg.stories, (el) => el.assets)
      assets = _.uniq(assets)
      return assets
    },
    nonAvgAssets () {
      const { assets } = this.$store.state.assets
      const { avgAssets } = this
      return assets.filter((el) => !avgAssets.includes(el))
    }
  }
}
</script>
