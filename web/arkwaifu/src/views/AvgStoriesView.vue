<template>
  <v-container class="px-8" v-if="story">
    <!-- <group-show v-if="prevGroup" :groupID="prevGroup.id" limited="true"></group-show> -->
    <group-show v-if="group" :groupID="group.id" :currentStoryID="story.id"></group-show>
    <!-- <group-show v-if="nextGroup" :groupID="nextGroup.id" limited="true"></group-show> -->
    <br />
    <v-row>
      <v-col cols="12" class="text-h5">Images</v-col>
      <v-col
        v-for="(image, i) in distinct ? _.uniq(story.images) : story.images"
        :key="`${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card :resName="image" resCategory="images"></resource-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" class="text-h5">Backgrounds</v-col>
      <v-col
        v-for="(image, i) in distinct ? _.uniq(story.backgrounds) : story.backgrounds"
        :key="`${image.ID}-${i}`"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card :resName="image" resCategory="backgrounds"></resource-card>
      </v-col>
    </v-row>
    <v-tooltip left>
      <template v-slot:activator="{ on, attrs }">
        <v-btn fab large fixed bottom right color="primary" @click="distinct = !distinct" v-bind="attrs" v-on="on">
          <v-icon v-if="!distinct">mdi-fingerprint-off</v-icon>
          <v-icon v-else>mdi-fingerprint</v-icon>
        </v-btn>
      </template>
      <span class="text-caption">Distinct: {{ distinct ? "ON" : "OFF" }}</span>
    </v-tooltip>
  </v-container>
</template>

<script>
import GroupShow from "@/components/GroupShow.vue";
import ResourceCard from "@/components/ResourceCard.vue";
export default {
  name: "AvgGroupsView",
  components: { GroupShow, ResourceCard },
  props: ["storyID"],
  data() {
    return {
      distinct: false,
    };
  },
  computed: {
    groups() {
      const groups = this.$store.state.avg.groups;
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
