<template>
  <div>
    <img
      ref="img"
      :src="fileURL"
      alt=""
    >
    <p class="text-caption">
      {{ kind }}/{{ name }} ({{ variant }}), {{ imgWidth }}×{{ imgHeight }} {{ filename }}
    </p>
    <p>
      <v-btn
        color="primary"
        :href="fileURL"
        :download="filename"
      >
        <v-icon left>
          {{ mdiDownload }}
        </v-icon>
        {{ $t('download') }}
      </v-btn>
    </p>
  </div>
</template>

<script>
import { mdiDownload } from '@mdi/js'
import Api from '@/api'

export default {
  name: 'AssetImg',
  props: {
    kind: {
      type: String,
      required: true
    },
    name: {
      type: String,
      required: true
    },
    variant: {
      type: String,
      required: false,
      default: 'img'
    }
  },
  data: function () {
    return {
      imgWidth: 0,
      imgHeight: 0,
      filename: '',
      mdiDownload
    }
  },
  computed: {
    fileURL: function () {
      return `${Api}/api/v0/asset/variants/${encodeURIComponent(this.kind)}/${encodeURIComponent(this.name)}/${encodeURIComponent(this.variant)}/file`
    },
    url: function () {
      return `${Api}/api/v0/asset/variants/${encodeURIComponent(this.kind)}/${encodeURIComponent(this.name)}/${encodeURIComponent(this.variant)}`
    }
  },
  created () {
    fetch(this.url)
      .then(res => res.json())
      .then(asset => (this.filename = asset.filename))
  },
  mounted () {
    this.$refs.img.onload = () => {
      const img = this.$refs.img
      this.imgWidth = img.naturalWidth
      this.imgHeight = img.naturalHeight
    }
  }
}
</script>

<style scoped>
  img {
    max-width: 100%;
    max-height: 100%;
  }
</style>
