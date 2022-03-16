<template>
  <v-tooltip left>
    <template #activator="{ on, attrs }">
      <v-btn
        fab
        large
        fixed
        bottom
        right
        color="primary"
        v-bind="attrs"
        @click="dataValue = !dataValue"
        v-on="on"
      >
        <v-icon v-if="!dataValue">
          mdi-fingerprint-off
        </v-icon>
        <v-icon v-else>
          mdi-fingerprint
        </v-icon>
      </v-btn>
    </template>
    <span class="text-caption">Distinct: {{ dataValue ? "ON" : "OFF" }}</span>
  </v-tooltip>
</template>
<script>
export default {
  name: 'DistinctButton',
  props: {
    value: Boolean
  },
  data () {
    return {
      dataValue: this.value
    }
  },
  watch: {
    dataValue (nv, ov) {
      if (nv !== ov) this.$emit('input', nv)
    }
  }
}
</script>
