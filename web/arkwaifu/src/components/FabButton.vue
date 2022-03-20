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
        <v-icon>{{ dataValue ? iconOn : iconOff }}</v-icon>
      </v-btn>
    </template>
    <span class="text-caption">{{ dataValue ? captionOn : captionOff }}</span>
  </v-tooltip>
</template>
<script>
export default {
  name: 'FabButton',
  props: {
    iconOn: {
      type: String,
      default: ''
    },
    iconOff: {
      type: String,
      default: ''
    },
    captionOn: {
      type: String,
      default: ''
    },
    captionOff: {
      type: String,
      default: ''
    },
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
