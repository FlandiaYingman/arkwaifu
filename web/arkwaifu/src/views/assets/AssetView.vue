<template>
  <v-container max-width="800px">
    <p class="text-h4">
      {{ kind }}/{{ id }}
    </p>
    <AssetImg
      :kind="kind"
      :name="id"
    />
    <p class="text-h4">
      {{ $t('variants') }}
    </p>
    <v-expansion-panels>
      <v-expansion-panel
        v-for="variant in variants"
        :key="variant"
      >
        <v-expansion-panel-header>
          {{ $t(`variant.${variant}`) }}
        </v-expansion-panel-header>
        <v-expansion-panel-content>
          <v-lazy>
            <AssetImg
              :kind="kind"
              :name="id"
              :variant="variant"
            />
          </v-lazy>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-container>
</template>

<script>
import AssetImg from '@/components/AssetImg'
import Api from '@/api'

export default {
  name: 'AssetView',
  components: { AssetImg },
  props: {
    kind: {
      required: true,
      type: String
    },
    id: {
      required: true,
      type: String
    }
  },
  computed: {},
  asyncComputed: {
    variants () {
      return fetch(`${Api}/api/v0/assets/kinds/${this.kind}/names/${this.id}/variants`)
        .then(res => res.json())
    }
  }
}
</script>

<i18n>{
  "en": {
    "variants": "Variants",
    "variant": {
      "img": "Raw Image",
      "timg": "Thumbnail Image",
      "real-esrgan": "Super-Resolution: Real-ESRGAN",
      "real-cugan": "Super-Resolution: Real-CUGAN"
    }
  },
  "zh":    {
  "variants": "变体",
  "variant": {
    "img": "原始图",
    "timg": "缩略图",
    "real-esrgan": "超分辨率: Real-ESRGAN",
    "real-cugan": "超分辨率: Real-CUGAN"
  }
}
}</i18n>
