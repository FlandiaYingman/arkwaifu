<template>
  <assets-show
    :assets="$store.state.assets"
  />
</template>
<script>
import AssetsShow from '@/components/AssetsShow'

export default {
  name: 'AllAssetsView',
  components: { AssetsShow }
}
</script>
