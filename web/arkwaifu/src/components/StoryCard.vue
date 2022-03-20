<template>
  <v-card
    :to="`/avgs/stories/${story.id}`"
    :disabled="current"
  >
    <v-card-title>{{ story.name }}</v-card-title>
    <v-card-subtitle>
      {{ story.code ? `${story.code} ` : '' }}{{ story.tag }}
      <span>
        <br v-show="!$vuetify.breakpoint.mobile">
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
  }
}
</script>
