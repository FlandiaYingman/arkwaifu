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
  </v-container>
</template>
<script>
import AssetCard from '@/components/AssetCard'
import _ from 'lodash'
export default {
  components: { AssetCard },
  props: {
    images: Array,
    backgrounds: Array,
    distinct: Boolean
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

<i18n>
{
  "en": {
    "images": "Images",
    "backgrounds": "Backgrounds"
  },
  "zh": {
    "images": "图片",
    "backgrounds": "背景"
  }

}
</i18n>
