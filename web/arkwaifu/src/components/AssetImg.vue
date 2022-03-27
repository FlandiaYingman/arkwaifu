<template>
  <div>
    <img
      ref="img"
      :src="url"
      alt=""
    >
    <p class="text-caption">
      {{ kind }}/{{ name }} ({{ variant }}), {{ imgWidth }}×{{ imgHeight }} {{ imgFormat }}
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
      imgFormat: '',
      imgWidth: 0,
      imgHeight: 0
    }
  },
  computed: {
    url: function () {
      return `${Api}/api/v0/assets/kinds/${this.kind}/names/${this.name}/variants/${this.variant}/file`
    }
  },
  mounted () {
    this.$refs.img.onload = () => {
      const img = this.$refs.img
      this.imgWidth = img.naturalWidth
      this.imgHeight = img.naturalHeight
    }
    fetch(this.url, { method: 'HEAD' })
      .then(resp => resp.headers.get('Content-Type'))
      .then(type => (this.imgFormat = type))
  }
}

</script>

<style scoped>
  img {
    max-width: 100%;
    max-height: 100%;
  }
</style>
