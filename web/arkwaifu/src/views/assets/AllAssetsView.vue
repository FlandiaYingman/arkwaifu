<template>
  <assets-show
    :images="images"
    :backgrounds="backgrounds"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'

export default {
  components: { AssetsShow },
  computed: {
    images () {
      return this.$store.state.assets.images
    },
    backgrounds () {
      return this.$store.state.assets.backgrounds
    }
  }
}
</script>
