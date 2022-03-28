<template>
  <div>
    <v-lazy
      v-for="(group, i) in groups"
      :key="i"
    >
      <group-show :group-id="group.id" />
    </v-lazy>
    <fab-button
      v-model="descending"
      :icon-on="mdiSortCalendarDescending"
      :icon-off="mdiSortCalendarAscending"
      caption-on="Sort: Descending"
      caption-off="Sort: Ascending"
    />
  </div>
</template>

<script>
import GroupShow from '@/components/GroupShow'
import FabButton from '@/components/FabButton'
import _ from 'lodash'
import { mdiSortCalendarAscending, mdiSortCalendarDescending } from '@mdi/js'

export default {
  name: 'GroupsShow',
  components: { GroupShow, FabButton },
  props: {
    type: String()
  },
  data () {
    return {
      descending: true,
      mdiSortCalendarDescending: mdiSortCalendarDescending,
      mdiSortCalendarAscending: mdiSortCalendarAscending
    }
  },
  computed: {
    groups () {
      let groups = [...this.$store.state.groupsTypeMap[this.type]]
      if (groups) {
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
