<template>
  <v-card
    v-ripple="{ class: `primary--text` }"
    :to="`/avgs/stories/${story.id}`"
    :disabled="current"
  >
    <v-card-title>{{ story.name }}</v-card-title>
    <v-card-subtitle>
      {{ story.code ? `${story.code} ${story.tag}` : `${story.tag}` }}
      <span v-show="!$vuetify.breakpoint.mobile">
        <br>
        {{ story.id }}
      </span>
    </v-card-subtitle>
  </v-card>
</template>

<script>
export default {
  name: 'StoryCard',
  props: {
    story: Object(),
    current: Boolean
  },
  computed: {
    codeTag () {
      if (this.story.code !== '') {
        return `${this.story.code} ${this.story.tag}`
      } else {
        return this.story.tag
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.v-card--disabled::before {
  content: "";
  opacity: 0.12 !important;
  background-color: var(--v-primary-base) !important;
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  top: 0;
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}

.v-card--disabled {
  opacity: 1 !important;
  color: var(--v-anchor-base) !important;
  caret-color: var(--v-anchor-base) !important;
}

.v-card--disabled > *:not(.v-card__progress) {
  opacity: 1 !important;
}
</style>
