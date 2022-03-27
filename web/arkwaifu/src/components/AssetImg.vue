<template>
  <div>
    <img
      ref="img"
      :src="url"
      alt=""
    >
    <p class="text-caption">
      {{ kind }}/{{ name }} ({{ variant }}), {{ imgWidth }}×{{ imgHeight }} {{ fileName }}
    </p>
  </div>
</template>

<script>
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
      fileName: ''
    }
  },
  computed: {
    url: function () {
      return `${Api}/api/v0/assets/kinds/${this.kind}/names/${this.name}/variants/${this.variant}/file`
    }
  },
  created () {
    fetch(`${Api}/api/v0/assets/kinds/${this.kind}/names/${this.name}/variants/${this.variant}`)
      .then(res => res.json())
      .then(asset => (this.fileName = asset.fileName))
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
