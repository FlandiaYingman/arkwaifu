<template>
  <v-container>
    <v-row
      v-for="group in groups"
      :key="group.id"
    >
      <group-show :group-id="group.id" />
    </v-row>
    <fab-button
      v-model="descending"
      icon-on="mdi-sort-calendar-descending"
      icon-off="mdi-sort-calendar-ascending"
      caption-on="Sort: Descending"
      caption-off="Sort: Ascending"
    />
  </v-container>
</template>

<script>
import GroupShow from '@/components/GroupShow'
import FabButton from '@/components/FabButton'
import _ from 'lodash'

export default {
  name: 'GroupsPage',
  components: { GroupShow, FabButton },
  props: {
    type: String()
  },
  data () {
    return {
      descending: true
    }
  },
  computed: {
    groups () {
      let groups = this.$store.state.avg.groupsTypeMap[this.type]
      if (groups) {
        groups = [...groups]
        if (this.descending) {
          groups = _.reverse(groups)
        }
        return groups
      }
      return null
    }
  }
}
</script>
