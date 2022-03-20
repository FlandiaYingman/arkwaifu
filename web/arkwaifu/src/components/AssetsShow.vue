<template>
  <v-container>
    <v-row
      v-for="kind in assetKinds"
      :key="kind"
    >
      <v-col
        cols="12"
        class="text-h5"
      >
        {{ $t(kind) }}
      </v-col>
      <v-col
        v-for="(asset, i) in assetsMap[kind]"
        :key="i"
        cols="6"
        sm="3"
        lg="2"
      >
        <asset-card
          :asset-id="asset.id"
          :asset-kind="kind"
        />
      </v-col>
    </v-row>
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
import AssetCard from '@/components/AssetCard'
import _ from 'lodash'
import FabButton from '@/components/FabButton'

export default {
  name: 'AssetsShow',
  components: { AssetCard, FabButton },
  props: {
    assets: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      distinct: true
    }
  },
  computed: {
    assetKinds () {
      return Object.keys(this.$store.state.assetsKindMap)
    },
    assetsMap () {
      let assets = this.assets
      if (this.distinct) {
        assets = _.uniqBy(assets, el => el.id)
      }
      return _.groupBy(assets, el => el.kind)
    }
  }
}
</script>

<i18n>{
  "en": {
    "images": "Images",
    "backgrounds": "Backgrounds",
    "distinctOn": "Distinct: ON",
    "distinctOff": "Distinct: Off"
  },
  "zh": {
    "images": "图片",
    "backgrounds": "背景",
    "distinctOn": "显示重复：关",
    "distinctOff": "显示重复：开"
  }
}</i18n>
