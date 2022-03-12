<template>
  <v-container class="px-8">
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        Images
      </v-col>
      <v-col
        v-for="image in nonAvgImages"
        :key="image"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card
          :res-name="image"
          res-category="images"
        />
      </v-col>
      <v-spacer />
    </v-row>
    <v-row>
      <v-col
        cols="12"
        class="text-h5"
      >
        Backgrounds
      </v-col>
      <v-col
        v-for="background in nonAvgBackgrounds"
        :key="background"
        cols="6"
        sm="3"
        lg="2"
      >
        <resource-card
          :res-name="background"
          res-category="backgrounds"
        />
      </v-col>
      <v-spacer />
    </v-row>
  </v-container>
</template>

<script>
import _ from 'lodash';
import ResourceCard from '@/components/ResourceCard.vue';

export default {
  name: 'NonAvgView',
  components: { ResourceCard },
  data() {
    return {
    };
  },
  computed: {
    avgImages() {
      let images = [];
      images = _.flatMap(this.$store.state.avg.stories, (el) => el.images);
      images = _.uniq(images);
      return images;
    },
    avgBackgrounds() {
      let backgrounds = [];
      backgrounds = _.flatMap(this.$store.state.avg.stories, (el) => el.backgrounds);
      backgrounds = _.uniq(backgrounds);
      return backgrounds;
    },
    nonAvgImages() {
      const { images } = this.$store.state.assets;
      const { avgImages } = this;
      return images.filter((el) => !avgImages.includes(el));
    },
    nonAvgBackgrounds() {
      const { backgrounds } = this.$store.state.assets;
      const { avgBackgrounds } = this;
      return backgrounds.filter((el) => !avgBackgrounds.includes(el));
    },
  },
};
</script>
