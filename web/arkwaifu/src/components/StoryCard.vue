<template>
  <v-card
    :to="`/avgs/stories/${story.id}`"
    :disabled="current"
  >
    <v-card-title>{{ story.name }}</v-card-title>
    <v-card-subtitle>
      {{ story.code ? `${story.code} ` : '' }}{{ story.tag }}
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
  }
}
</script>
