<template>
  <assets-show
    :assets="nonAvgAssets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'

export default {
  name: 'NonAvgAssetsView',
  components: { AssetsShow },
  computed: {
    nonAvgAssets () {
      const allAssets = this.$store.state.assets
      const avgAssets = this.$store.state.storiesAssetsMap
      const nonAvgAssets = allAssets.filter(allEl => !(allEl.id in avgAssets))
      return nonAvgAssets
    }
  }
}
</script>
