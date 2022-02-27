<template>
  <div>
    <div class="text-caption">{{ this.resCategory }}/{{ this.resName }}</div>
    <v-card :href="this.resURL" target="_blank">
      <v-img :src="`${this.resURL}?resType=thumbnail`" class="transparent-background"></v-img>
    </v-card>
  </div>
</template>

<style scoped>
.transparent-background {
  background: linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    ),
    linear-gradient(
      45deg,
      rgba(0, 0, 0, 0.25) 25%,
      transparent 25%,
      transparent 75%,
      rgba(0, 0, 0, 0.25) 75%,
      rgba(0, 0, 0, 0.25) 0
    );
  background-position: 0px 0, 10px 10px;
  background-origin: padding-box, padding-box;
  background-clip: border-box, border-box;
  background-size: 20px 20px, 20px 20px;
  box-shadow: none;
}
.overlay-text {
  /* position: absolute; */
  /* top: 0; */
  /* left: 0; */
}
</style>

<script>
export default {
  name: "ResourceCard",
  props: ["resName", "resCategory"],
  data() {
    return {
      resURL: `${this.$API_URL}/api/v0/resources/${this.resCategory}/${this.resName}`,
    };
  },
};
</script>
