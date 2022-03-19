<template>
  <assets-show
    :assets="assets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'

export default {
  components: { AssetsShow },
  computed: {
    assets () {
      return this.$store.state.assets.assets
    }
  }
}
</script>
