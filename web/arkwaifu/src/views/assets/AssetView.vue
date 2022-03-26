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
      Variants
    </p>
    <v-expansion-panels>
      <v-expansion-panel
        v-for="variant in variants"
        :key="variant"
      >
        <v-expansion-panel-header>{{ variant }}</v-expansion-panel-header>
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
