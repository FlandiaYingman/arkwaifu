<template>
  <v-container>
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        {{ $t("images") }}
      </v-col>
      <v-col
        v-for="(image, i) in imagesData"
        :key="`${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <asset-card
          :res-name="image"
          res-category="images"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        {{ $t("backgrounds") }}
      </v-col>
      <v-col
        v-for="(image, i) in backgroundsData"
        :key=" `${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <asset-card
          :res-name="image"
          res-category="backgrounds"
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
  components: { AssetCard, FabButton },
  props: {
    images: Array,
    backgrounds: Array
  },
  data () {
    return {
      distinct: true
    }
  },
  computed: {
    imagesData () {
      if (this.distinct) {
        return _.uniq(this.images)
      } else {
        return this.images
      }
    },
    backgroundsData () {
      if (this.distinct) {
        return _.uniq(this.backgrounds)
      } else {
        return this.backgrounds
      }
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
