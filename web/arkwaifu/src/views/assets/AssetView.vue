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
    <p class="text-body-1">
      {{ $t('variants.descVariants') }}
    </p>
    <p class="text-body-1">
      {{ $t('variants.descSuperResolution') }}
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
      "desc": "These are the variants of the asset. Note that some files are huge (~10MB) and may take a while to load.",
      "descVariants": "The Raw Image is the original image extracted from the game. The Super-Resolution Images are the high-resolution versions of image after super-resolution with different models (Real-CUGAN, Real-ESRGAN, etc.). Thumbnail is the low-resolution and low-quality version of the image. ",
      "descSuperResolution": "Real-CUGAN and Real-ESRGAN are two different models. There are slightly differences between the produced images. For me, Real-CUGAN looks more massive, and Real-ESRGAN looks more fine and smooth. "
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
      "desc": "这些是资源的变体。注意，一些文件是很大的（~10MB），可能需要一段时间来加载。",
      "descVariants": "原始图指从游戏中提取的原始图像；超分辨率指使用不同模型（Real-CUGAN、Real-ESRGAN 等）进行超分辨率后的高清版原始图像；缩略图指原始图像降低分辨率和质量后的缩略图。",
      "descSuperResolution": "Real-CUGAN 和 Real-ESRGAN 是两个不同的模型，它们的超分辨率效果有细微的不同。就作者本人而言，Real-CUGAN 更厚重，而 Real-ESRGAN 更细腻。"
    },
    "variant": {
      "img": "原始图",
      "timg": "缩略图",
      "real-esrgan": "超分辨率: Real-ESRGAN",
      "real-cugan": "超分辨率: Real-CUGAN"
    }
  }
}</i18n>
