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
      {{ $t('variants.title') }}
    </p>
    <p class="text-body-1">
      {{ $t('variants.desc') }}
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
    "variants": {
      "title": "Variants",
      "desc": "These are the variants of the asset. Note that some files are huge (~10MB) and may take a while to load."
    },
    "variant": {
      "img": "Raw Image",
      "timg": "Thumbnail Image",
      "real-esrgan": "Super-Resolution: Real-ESRGAN",
      "real-cugan": "Super-Resolution: Real-CUGAN"
    }
  },
  "zh": {
    "variants": {
      "title": "变体",
      "desc": "这些是资源的变体。注意，一些文件是很大的（~10MB），可能需要一段时间来加载。"
    },
    "variant": {
    "img": "原始图",
    "timg": "缩略图",
    "real-esrgan": "超分辨率: Real-ESRGAN",
    "real-cugan": "超分辨率: Real-CUGAN"
  }
}
}</i18n>
