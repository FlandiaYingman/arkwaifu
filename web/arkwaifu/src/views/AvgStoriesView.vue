<template>
  <v-container
    v-if="story"
    class="px-8"
  >
    <!-- <group-show v-if="prevGroup" :groupID="prevGroup.id" limited="true"></group-show> -->
    <group-show
      v-if="group"
      :group-i-d="group.id"
      :current-story-i-d="story.id"
    />
    <!-- <group-show v-if="nextGroup" :groupID="nextGroup.id" limited="true"></group-show> -->
    <br>
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        Images
      </v-col>
      <v-col
        v-for="(image, i) in distinct ? _.uniq(story.images) : story.images"
        :key="`${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card
          :res-name="image"
          res-category="images"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        Backgrounds
      </v-col>
      <v-col
        v-for="(image, i) in distinct ? _.uniq(story.backgrounds) : story.backgrounds"
        :key="`${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card
          :res-name="image"
          res-category="backgrounds"
        />
      </v-col>
    </v-row>
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
          @click="distinct = !distinct"
          v-on="on"
        >
          <v-icon v-if="!distinct">
            mdi-fingerprint-off
          </v-icon>
          <v-icon v-else>
            mdi-fingerprint
          </v-icon>
        </v-btn>
      </template>
      <span class="text-caption">Distinct: {{ distinct ? "ON" : "OFF" }}</span>
    </v-tooltip>
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow.vue';
import ResourceCard from '@/components/ResourceCard.vue';

export default {
  name: 'AvgGroupsView',
  components: { GroupShow, ResourceCard },
  props: {
    storyID: String(),
  },
  data() {
    return {
      distinct: true,
    };
  },
  computed: {
    groups() {
      const { groups } = this.$store.state.avg;
      return groups;
    },
    group() {
      const group = this.$store.getters.groupByID(this.story.groupID);
      return group;
    },
    story() {
      const story = this.$store.getters.storyByID(this.storyID);
      return story;
    },
  },
};
</script>
