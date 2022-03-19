<template>
  <assets-show
    :assets="nonAvgAssets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'
import fp from 'lodash/fp'
import _ from 'lodash'

export default {
  components: { AssetsShow },
  computed: {
    avgAssets () {
      const stories = this.$store.state.avg.stories
      const assets = fp.flow(
        fp.flatMap(it => it.assets),
        fp.uniqWith(_.isEqual)
      )(stories)
      return assets
    },
    nonAvgAssets () {
      const allAssets = this.$store.state.assets.assets
      const avgAssets = this.avgAssets
      const nonAvgAssets = allAssets.filter(allEl => !avgAssets.find(avgEl => _.isEqual(allEl, avgEl)))
      return nonAvgAssets
    }
  }
}
</script>
