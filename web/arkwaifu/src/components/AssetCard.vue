<template>
  <div>
    <div class="text-caption">
      {{ kindAbbr }}/{{ assetId }}
    </div>
    <v-card
      :to="`/assets/${assetKind}/${encodeURIComponent(assetId)}`"
      :href="assetURL"
    >
      <v-img
        :src="assetThumbnailURL"
        class="transparent-background"
      />
    </v-card>
  </div>
</template>

<script>
import API_URL from '@/api'

export default {
  name: 'AssetCard',
  props: {
    assetId: String(),
    assetKind: String()
  },
  data () {
    return {
    }
  },
  computed: {
    assetThumbnailURL () {
      return `${API_URL}/api/v0/asset/variants/${encodeURIComponent(this.assetKind)}/${encodeURIComponent(this.assetId)}/timg/file`
    },
    assetURL () {
      return `${API_URL}/api/v0/asset/variants/${encodeURIComponent(this.assetKind)}/${encodeURIComponent(this.assetId)}/img/file`
    },
    kindAbbr () {
      switch (this.assetKind) {
        case 'images':
          return 'img'
        case 'backgrounds':
          return 'bg'
        case 'characters':
          return 'char'
        default:
          return this.assetKind
      }
    }
  }
}
</script>

<style scoped>
.transparent-background {
  background: linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    ), linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    );
  background-position: 0 0, 10px 10px;
  background-origin: padding-box, padding-box;
  background-clip: border-box, border-box;
  background-size: 20px 20px, 20px 20px;
  box-shadow: none;
}
</style>
