<template>
  <div>
    <div class="text-caption">
      {{ resCategory }}/{{ resName }}
    </div>
    <v-card
      :href="assetURL"
      target="_blank"
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
  name: 'ResourceCard',
  props: {
    resName: String(),
    resCategory: String()
  },
  data () {
    return {
    }
  },
  computed: {
    assetThumbnailURL () {
      return `${API_URL}/api/v0/resources/${this.resCategory}/thumbnail/${this.resName}`
    },
    assetURL () {
      return `${API_URL}/api/v0/resources/${this.resCategory}/raw/${this.resName}`
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
