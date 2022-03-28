<template>
  <assets-show
    :assets="allAssets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'

export default {
  name: 'AllAssetsView',
  components: { AssetsShow },
  computed: {
    allAssets () {
      const allAssets = this.$store.state.assetsMap
      return Object.values(allAssets)
    }
  }
}
</script>
