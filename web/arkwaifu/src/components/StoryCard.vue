<template>
  <v-card :to="`/avgs/stories/${story.ID}`">
    <v-card-title>{{ story.Name }}</v-card-title>
    <v-card-subtitle>
      {{ story.Code ? `${story.Code} ` : "" }}{{ story.Tag }}
      <br />
      {{ story.ID }}
    </v-card-subtitle>
  </v-card>
</template>

<script>
export default {
  name: "StoryCard",
  props: ["story"],
};
</script>
