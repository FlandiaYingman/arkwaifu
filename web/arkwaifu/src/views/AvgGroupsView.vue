<template>
  <group-chooser-vue></group-chooser-vue>
</template>

<script>
import GroupChooserVue from '@/components/GroupChooser.vue';
export default {
  name: "AvgGroupsView",
  components: {
      GroupChooserVue
  },
  data() {
    return {
      images: [],
    };
  },
  created() {
    this.$watch(
      () => this.$route.params,
      () => {
        this.fetchImages();
      },
      { immediate: true }
    );
  },
  methods: {
    async fetchImages() {
      const response = await fetch("http://localhost:7080/api/v0/resources/images");
      this.images = (await response.json()).map((it) => `http://localhost:7080/api/v0/resources/images/${it}`);
    },
  },
};
</script>
