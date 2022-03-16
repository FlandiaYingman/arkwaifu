<template>
  <v-container>
    <v-row
      v-for="group in groups"
      :key="group.id"
    >
      <group-show :group-id="group.id" />
    </v-row>
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow'

export default {
  name: 'GroupsPage',
  components: { GroupShow },
  props: {
    type: String()
  },
  data () {
    return {}
  },
  computed: {
    groups () {
      return this.$store.state.avg.groupsTypeMap[this.type]
    }
  }
}
</script>
