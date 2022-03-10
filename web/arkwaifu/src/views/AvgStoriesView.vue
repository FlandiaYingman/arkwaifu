<template>
  <v-container class="px-8" v-if="story">
    <!-- <group-show v-if="prevGroup" :groupID="prevGroup.id" limited="true"></group-show> -->
    <group-show v-if="group" :groupID="group.id" :currentStoryID="story.id"></group-show>
    <!-- <group-show v-if="nextGroup" :groupID="nextGroup.id" limited="true"></group-show> -->
    <br />
    <v-row>
      <v-col cols="12" class="text-h5">Images</v-col>
      <v-col v-for="(image, i) in story.images" :key="`${image.ID}-${i}`" cols="6" sm="3" lg="2">
        <resource-card :resName="image" resCategory="images"></resource-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" class="text-h5">Backgrounds</v-col>
      <v-col v-for="(image, i) in story.backgrounds" :key="`${image.ID}-${i}`" cols="6" sm="3" lg="2">
        <resource-card :resName="image" resCategory="backgrounds"></resource-card>
      </v-col>
    </v-row>
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
      story: null,
      groups: [],
      group: null,
      // prevGroup: null,
      // nextGroup: null,
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.fetchStory(this.storyID).then(() => this.fetchGroups());
      }
    );
    this.fetchStory(this.storyID).then(() => this.fetchGroups());
  },
  methods: {
    async fetchStory(storyID) {
      return fetch(`${this.$API_URL}/api/v0/stories/${storyID}`)
        .then((resp) => resp.json())
        .then((story) => (this.story = story));
    },
    async fetchGroups() {
      return fetch(`${this.$API_URL}/api/v0/groups`)
        .then((resp) => resp.json())
        .then((group) => (this.groups = group))
        .then(() => {
          let i = this.groups.findIndex((it) => it.id == this.story.groupID);
          this.group = this.groups[i];
          // if (i - 1 >= 0) {
          //   this.prevGroup = this.groups[i - 1];
          // }
          // if (i + 1 < this.groups.length) {
          //   this.nextGroup = this.groups[i + 1];
          // }
        });
    },
  },
};
</script>
